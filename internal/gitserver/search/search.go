package search

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/errgroup"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/gitserver/protocol"
	"github.com/sourcegraph/sourcegraph/internal/search/result"
)

// Git formatting directives as described in man git-log (see PRETTY FORMATS)
const (
	hash           = "%H"
	refNames       = "%D"
	sourceRefs     = "%S"
	authorName     = "%aN"
	authorEmail    = "%aE"
	authorDate     = "%at"
	committerName  = "%cN"
	committerEmail = "%cE"
	committerDate  = "%ct"
	rawBody        = "%B"
	parentHashes   = "%P"
)

var (
	partsPerCommit = len(formatWithRefs)

	formatWithRefs = []string{
		hash,
		refNames,
		sourceRefs,
		authorName,
		authorEmail,
		authorDate,
		committerName,
		committerEmail,
		committerDate,
		rawBody,
		parentHashes,
	}

	formatWithoutRefs = []string{
		hash,
		"",
		sourceRefs,
		authorName,
		authorEmail,
		authorDate,
		committerName,
		committerEmail,
		committerDate,
		rawBody,
		parentHashes,
	}

	baseLogArgs = []string{
		"log",
		"--decorate=full",
		"-z",
		"--no-merges",
	}

	// TODO(@camdencheek) support adding refs (issue #25356)
	// logArgsWithRefs    = append(baseLogArgs, "--format=format:"+strings.Join(formatWithRefs, "%x00")+"%x00")
	logArgsWithoutRefs = append(baseLogArgs, "--format=format:"+strings.Join(formatWithoutRefs, "%x00")+"%x00")
	sep                = []byte{0x0}
)

type job struct {
	batch      []*RawCommit
	resultChan chan *protocol.CommitMatch
}

const (
	// The size of a batch of commits sent in each worker job
	batchSize  = 512
	numWorkers = 4
)

type CommitSearcher struct {
	RepoDir     string
	Query       MatchTree
	Revisions   []protocol.RevisionSpecifier
	IncludeDiff bool
}

// Search runs a search for commits matching the given predicate across the revisions passed in as revisionArgs.
//
// We have some slightly complex logic here in order to run searches in parallel (big benefit to diff searches),
// but also return results in order. We first iterate over all the commits using the hard-coded git log arguments.
// We batch the shallowly-parsed commits, then send them on the jobs channel along with a channel that results for
// that job should be sent down. We then read from the result channels in the same order that the jobs were sent.
// This allows our worker pool to run the jobs in parallel, but we still emit matches in the same order that
// git log outputs them.
func (cs *CommitSearcher) Search(ctx context.Context, onMatch func(*protocol.CommitMatch)) error {
	g, ctx := errgroup.WithContext(ctx)

	jobs := make(chan job, 128)
	resultChans := make(chan chan *protocol.CommitMatch, 128)

	// Start feeder
	g.Go(func() error {
		defer close(resultChans)
		defer close(jobs)
		return cs.feedBatches(ctx, jobs, resultChans)
	})

	// Start workers
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			return cs.runJobs(ctx, jobs)
		})
	}

	// Consumer goroutine that consumes results in the order jobs were
	// submitted to the job queue
	g.Go(func() error {
		for resultChan := range resultChans {
			for result := range resultChan {
				onMatch(result)
			}
		}

		return nil
	})

	return g.Wait()
}

func (cs *CommitSearcher) feedBatches(ctx context.Context, jobs chan job, resultChans chan chan *protocol.CommitMatch) error {
	revArgs := revsToGitArgs(cs.Revisions)
	cmd := exec.CommandContext(ctx, "git", append(logArgsWithoutRefs, revArgs...)...)
	cmd.Dir = cs.RepoDir
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	batch := make([]*RawCommit, 0, batchSize)
	sendBatch := func() {
		resultChan := make(chan *protocol.CommitMatch, 128)
		resultChans <- resultChan
		jobs <- job{
			batch:      batch,
			resultChan: resultChan,
		}
		batch = make([]*RawCommit, 0, batchSize)
	}

	scanner := NewCommitScanner(stdoutReader)
	for scanner.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		cv := scanner.NextRawCommit()
		batch = append(batch, cv)
		if len(batch) == batchSize {
			sendBatch()
		}
	}

	if len(batch) > 0 {
		sendBatch()
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return cmd.Wait()
}

func (cs *CommitSearcher) runJobs(ctx context.Context, jobs chan job) error {
	// Create a new diff fetcher subprocess for each worker
	diffFetcher, err := StartDiffFetcher(cs.RepoDir)
	if err != nil {
		return err
	}
	defer diffFetcher.Stop()

	startBuf := make([]byte, 1024)

	runJob := func(j job) error {
		defer close(j.resultChan)

		for _, cv := range j.batch {
			if ctx.Err() != nil {
				// ignore context error, and don't spend time running the job
				return nil
			}

			lc := &LazyCommit{
				RawCommit:   cv,
				diffFetcher: diffFetcher,
				LowerBuf:    startBuf,
			}
			commitMatches, highlights, err := cs.Query.Match(lc)
			if err != nil {
				return err
			}
			if commitMatches {
				cm, err := CreateCommitMatch(lc, highlights, cs.IncludeDiff)
				if err != nil {
					return err
				}
				j.resultChan <- cm
			}
		}
		return nil
	}

	var errors error
	for j := range jobs {
		multierror.Append(errors, runJob(j))
	}
	return errors
}

func revsToGitArgs(revs []protocol.RevisionSpecifier) []string {
	revArgs := make([]string, 0, len(revs))
	for _, rev := range revs {
		if rev.RevSpec != "" {
			revArgs = append(revArgs, rev.RevSpec)
		} else if rev.RefGlob != "" {
			revArgs = append(revArgs, "--glob="+rev.RefGlob)
		} else if rev.ExcludeRefGlob != "" {
			revArgs = append(revArgs, "--exclude="+rev.RefGlob)
		} else {
			revArgs = append(revArgs, "HEAD")
		}
	}
	return revArgs
}

// RawCommit is a shallow parse of the output of git log
type RawCommit struct {
	Hash           []byte
	RefNames       []byte
	SourceRefs     []byte
	AuthorName     []byte
	AuthorEmail    []byte
	AuthorDate     []byte
	CommitterName  []byte
	CommitterEmail []byte
	CommitterDate  []byte
	Message        []byte
	ParentHashes   []byte
}

type CommitScanner struct {
	scanner *bufio.Scanner
	next    *RawCommit
	err     error
}

// NewCommitScanner creates a scanner that does a shallow parse of the stdout of git log.
// Like the bufio.Scanner() API, call Scan() to ingest the next result, which will return
// false if it hits an error or EOF, then call NextRawCommit() to get the scanned commit.
func NewCommitScanner(r io.Reader) *CommitScanner {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1<<22)

	// Split by commit
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 { // should only happen when atEOF
			return 0, nil, nil
		}

		// See if we have enough null bytes to constitute a full commit
		// Look for one more than the number of parts because the each message ends with a null byte too
		sepCount := bytes.Count(data, sep)
		if sepCount < partsPerCommit+1 {
			if atEOF {
				if sepCount == partsPerCommit {
					return len(data), data, nil
				}
				return 0, nil, errors.Errorf("incomplete line")
			}
			return 0, nil, nil
		}

		// If we do, expand token to the end of that commit
		for i := 0; i < partsPerCommit; i++ {
			idx := bytes.IndexByte(data[len(token):], 0x0)
			if idx == -1 {
				panic("we already counted enough bytes in data")
			}
			token = data[:len(token)+idx+1]
		}
		return len(token) + 1, token, nil
	})

	return &CommitScanner{
		scanner: scanner,
	}
}

func (c *CommitScanner) Scan() bool {
	if !c.scanner.Scan() {
		return false
	}

	// Make a copy so the view can outlive the next scan
	buf := make([]byte, len(c.scanner.Bytes()))
	copy(buf, c.scanner.Bytes())

	parts := bytes.SplitN(buf, sep, partsPerCommit+1)
	if len(parts) < partsPerCommit+1 {
		c.err = errors.Errorf("invalid commit log entry: %q", parts)
		return false
	}

	c.next = &RawCommit{
		Hash:           parts[0],
		RefNames:       parts[1],
		SourceRefs:     parts[2],
		AuthorName:     parts[3],
		AuthorEmail:    parts[4],
		AuthorDate:     parts[5],
		CommitterName:  parts[6],
		CommitterEmail: parts[7],
		CommitterDate:  parts[8],
		Message:        bytes.TrimSpace(parts[9]),
		ParentHashes:   parts[10],
	}

	return true
}

func (c *CommitScanner) NextRawCommit() *RawCommit {
	return c.next
}

func (c *CommitScanner) Err() error {
	return c.err
}

func CreateCommitMatch(lc *LazyCommit, hc *MatchedCommit, includeDiff bool) (*protocol.CommitMatch, error) {
	if hc == nil {
		hc = &MatchedCommit{}
	}

	authorDate, err := lc.AuthorDate()
	if err != nil {
		return nil, err
	}

	committerDate, err := lc.CommitterDate()
	if err != nil {
		return nil, err
	}

	diff := result.MatchedString{}
	if includeDiff {
		rawDiff, err := lc.Diff()
		if err != nil {
			return nil, err
		}
		diff.Content, diff.MatchedRanges = FormatDiff(rawDiff, hc.Diff)
	}

	return &protocol.CommitMatch{
		Oid: api.CommitID(string(lc.Hash)),
		Author: protocol.Signature{
			Name:  string(lc.AuthorName),
			Email: string(lc.AuthorEmail),
			Date:  authorDate,
		},
		Committer: protocol.Signature{
			Name:  string(lc.CommitterName),
			Email: string(lc.CommitterEmail),
			Date:  committerDate,
		},
		Parents:    lc.ParentIDs(),
		SourceRefs: lc.SourceRefs(),
		Refs:       lc.RefNames(),
		Message: result.MatchedString{
			Content:       string(lc.Message),
			MatchedRanges: hc.Message,
		},
		Diff: diff,
	}, nil
}
