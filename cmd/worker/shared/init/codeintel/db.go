package codeintel

import (
	"database/sql"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/cmd/worker/memo"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores"
	"github.com/sourcegraph/sourcegraph/internal/conf"
	"github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
	connections "github.com/sourcegraph/sourcegraph/internal/database/connections/live"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

// InitDB initializes and returns a connection to the codeintel db.
func InitDB() (*sql.DB, error) {
	return initDBMemo.Init()
}

func InitDBWithLogger(logger log.Logger) (stores.CodeIntelDB, error) {
	rawDB, err := InitDB()
	if err != nil {
		return nil, err
	}

	return stores.NewCodeIntelDB(rawDB), nil
}

var initDBMemo = memo.NewMemoizedConstructor(func() (*sql.DB, error) {
	dsn := conf.GetServiceConnectionValueAndRestartOnChange(func(serviceConnections conftypes.ServiceConnections) string {
		return serviceConnections.CodeIntelPostgresDSN
	})
	db, err := connections.EnsureNewCodeIntelDB(dsn, "worker", &observation.TestContext)
	if err != nil {
		return nil, errors.Errorf("failed to connect to codeintel database: %s", err)
	}

	return db, nil
})