package clickhouse_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/internal/testdb"
)

func TestClickHouse(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	migrationDir := filepath.Join("testdata", "migrations")
	db, cleanup, err := testdb.NewClickHouse()
	is.NoErr(err)
	t.Cleanup(cleanup)

	goose.SetDialect("clickhouse")

	currentVersion, err := goose.GetDBVersion(db)
	is.NoErr(err)
	is.Equal(currentVersion, int64(0))

	err = goose.Up(db, migrationDir)
	is.NoErr(err)

	currentVersion, err = goose.GetDBVersion(db)
	is.NoErr(err)
	is.Equal(currentVersion, int64(1))

	err = goose.DownTo(db, migrationDir, 0)
	is.NoErr(err)
	// TODO(mf): this will fail if SETTINGS mutations_sync = 0.. because the
	// delete operation above is async and that is the default. I updated the
	// dialect string, but not sure that's correct..
	time.Sleep(5 * time.Second)
	// Add retry here?
	currentVersion, err = goose.GetDBVersion(db)
	is.NoErr(err)
	is.Equal(currentVersion, int64(0))
}
