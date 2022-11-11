// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package auth

import (
	"context"
	"sync"

	github "github.com/sourcegraph/sourcegraph/internal/extsvc/github"
)

// MockGitHubClient is a mock implementation of the GitHubClient interface
// (from the package
// github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/uploads/transport/http/auth)
// used for unit testing.
type MockGitHubClient struct {
	// GetRepositoryFunc is an instance of a mock function object
	// controlling the behavior of the method GetRepository.
	GetRepositoryFunc *GitHubClientGetRepositoryFunc
	// ListInstallationRepositoriesFunc is an instance of a mock function
	// object controlling the behavior of the method
	// ListInstallationRepositories.
	ListInstallationRepositoriesFunc *GitHubClientListInstallationRepositoriesFunc
}

// NewMockGitHubClient creates a new mock of the GitHubClient interface. All
// methods return zero values for all results, unless overwritten.
func NewMockGitHubClient() *MockGitHubClient {
	return &MockGitHubClient{
		GetRepositoryFunc: &GitHubClientGetRepositoryFunc{
			defaultHook: func(context.Context, string, string) (r0 *github.Repository, r1 error) {
				return
			},
		},
		ListInstallationRepositoriesFunc: &GitHubClientListInstallationRepositoriesFunc{
			defaultHook: func(context.Context, int) (r0 []*github.Repository, r1 bool, r2 int, r3 error) {
				return
			},
		},
	}
}

// NewStrictMockGitHubClient creates a new mock of the GitHubClient
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockGitHubClient() *MockGitHubClient {
	return &MockGitHubClient{
		GetRepositoryFunc: &GitHubClientGetRepositoryFunc{
			defaultHook: func(context.Context, string, string) (*github.Repository, error) {
				panic("unexpected invocation of MockGitHubClient.GetRepository")
			},
		},
		ListInstallationRepositoriesFunc: &GitHubClientListInstallationRepositoriesFunc{
			defaultHook: func(context.Context, int) ([]*github.Repository, bool, int, error) {
				panic("unexpected invocation of MockGitHubClient.ListInstallationRepositories")
			},
		},
	}
}

// NewMockGitHubClientFrom creates a new mock of the MockGitHubClient
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockGitHubClientFrom(i GitHubClient) *MockGitHubClient {
	return &MockGitHubClient{
		GetRepositoryFunc: &GitHubClientGetRepositoryFunc{
			defaultHook: i.GetRepository,
		},
		ListInstallationRepositoriesFunc: &GitHubClientListInstallationRepositoriesFunc{
			defaultHook: i.ListInstallationRepositories,
		},
	}
}

// GitHubClientGetRepositoryFunc describes the behavior when the
// GetRepository method of the parent MockGitHubClient instance is invoked.
type GitHubClientGetRepositoryFunc struct {
	defaultHook func(context.Context, string, string) (*github.Repository, error)
	hooks       []func(context.Context, string, string) (*github.Repository, error)
	history     []GitHubClientGetRepositoryFuncCall
	mutex       sync.Mutex
}

// GetRepository delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockGitHubClient) GetRepository(v0 context.Context, v1 string, v2 string) (*github.Repository, error) {
	r0, r1 := m.GetRepositoryFunc.nextHook()(v0, v1, v2)
	m.GetRepositoryFunc.appendCall(GitHubClientGetRepositoryFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the GetRepository method
// of the parent MockGitHubClient instance is invoked and the hook queue is
// empty.
func (f *GitHubClientGetRepositoryFunc) SetDefaultHook(hook func(context.Context, string, string) (*github.Repository, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetRepository method of the parent MockGitHubClient instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *GitHubClientGetRepositoryFunc) PushHook(hook func(context.Context, string, string) (*github.Repository, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *GitHubClientGetRepositoryFunc) SetDefaultReturn(r0 *github.Repository, r1 error) {
	f.SetDefaultHook(func(context.Context, string, string) (*github.Repository, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *GitHubClientGetRepositoryFunc) PushReturn(r0 *github.Repository, r1 error) {
	f.PushHook(func(context.Context, string, string) (*github.Repository, error) {
		return r0, r1
	})
}

func (f *GitHubClientGetRepositoryFunc) nextHook() func(context.Context, string, string) (*github.Repository, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GitHubClientGetRepositoryFunc) appendCall(r0 GitHubClientGetRepositoryFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of GitHubClientGetRepositoryFuncCall objects
// describing the invocations of this function.
func (f *GitHubClientGetRepositoryFunc) History() []GitHubClientGetRepositoryFuncCall {
	f.mutex.Lock()
	history := make([]GitHubClientGetRepositoryFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GitHubClientGetRepositoryFuncCall is an object that describes an
// invocation of method GetRepository on an instance of MockGitHubClient.
type GitHubClientGetRepositoryFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *github.Repository
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c GitHubClientGetRepositoryFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GitHubClientGetRepositoryFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// GitHubClientListInstallationRepositoriesFunc describes the behavior when
// the ListInstallationRepositories method of the parent MockGitHubClient
// instance is invoked.
type GitHubClientListInstallationRepositoriesFunc struct {
	defaultHook func(context.Context, int) ([]*github.Repository, bool, int, error)
	hooks       []func(context.Context, int) ([]*github.Repository, bool, int, error)
	history     []GitHubClientListInstallationRepositoriesFuncCall
	mutex       sync.Mutex
}

// ListInstallationRepositories delegates to the next hook function in the
// queue and stores the parameter and result values of this invocation.
func (m *MockGitHubClient) ListInstallationRepositories(v0 context.Context, v1 int) ([]*github.Repository, bool, int, error) {
	r0, r1, r2, r3 := m.ListInstallationRepositoriesFunc.nextHook()(v0, v1)
	m.ListInstallationRepositoriesFunc.appendCall(GitHubClientListInstallationRepositoriesFuncCall{v0, v1, r0, r1, r2, r3})
	return r0, r1, r2, r3
}

// SetDefaultHook sets function that is called when the
// ListInstallationRepositories method of the parent MockGitHubClient
// instance is invoked and the hook queue is empty.
func (f *GitHubClientListInstallationRepositoriesFunc) SetDefaultHook(hook func(context.Context, int) ([]*github.Repository, bool, int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// ListInstallationRepositories method of the parent MockGitHubClient
// instance invokes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *GitHubClientListInstallationRepositoriesFunc) PushHook(hook func(context.Context, int) ([]*github.Repository, bool, int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *GitHubClientListInstallationRepositoriesFunc) SetDefaultReturn(r0 []*github.Repository, r1 bool, r2 int, r3 error) {
	f.SetDefaultHook(func(context.Context, int) ([]*github.Repository, bool, int, error) {
		return r0, r1, r2, r3
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *GitHubClientListInstallationRepositoriesFunc) PushReturn(r0 []*github.Repository, r1 bool, r2 int, r3 error) {
	f.PushHook(func(context.Context, int) ([]*github.Repository, bool, int, error) {
		return r0, r1, r2, r3
	})
}

func (f *GitHubClientListInstallationRepositoriesFunc) nextHook() func(context.Context, int) ([]*github.Repository, bool, int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GitHubClientListInstallationRepositoriesFunc) appendCall(r0 GitHubClientListInstallationRepositoriesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// GitHubClientListInstallationRepositoriesFuncCall objects describing the
// invocations of this function.
func (f *GitHubClientListInstallationRepositoriesFunc) History() []GitHubClientListInstallationRepositoriesFuncCall {
	f.mutex.Lock()
	history := make([]GitHubClientListInstallationRepositoriesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GitHubClientListInstallationRepositoriesFuncCall is an object that
// describes an invocation of method ListInstallationRepositories on an
// instance of MockGitHubClient.
type GitHubClientListInstallationRepositoriesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*github.Repository
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 bool
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 int
	// Result3 is the value of the 4th result returned from this method
	// invocation.
	Result3 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c GitHubClientListInstallationRepositoriesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GitHubClientListInstallationRepositoriesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2, c.Result3}
}
