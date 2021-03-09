package tconf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/jrapoport/gothic/test/tutils"
	mysqld "github.com/lestrrat-go/test-mysqld"
	"github.com/stretchr/testify/require"
)

// TempDB loads a temp sqlite db for tests.
func TempDB(t *testing.T) *config.Config {
	return DBConfig(t, DBTemp)
}

func mysqldDB(t *testing.T, c *config.Config) *config.Config {
	cfg := mysqld.NewConfig()
	db, err := mysqld.NewMysqld(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Cleanup(func() {
		db.Stop()
	})
	c.DB.Driver = drivers.MySQL
	c.DB.DSN = db.DSN(
		mysqld.WithDbname(""),
		mysqld.WithParseTime(true),
		mysqld.WithMultiStatements(true))
	return c
}

func sqliteDB(t *testing.T, c *config.Config) *config.Config {
	path := c.DB.DSN
	dir, file := filepath.Split(path)
	if !tutils.PathExists(t, dir) {
		root := tutils.ProjectRoot(t)
		path = filepath.Join(root, dir)
		// make sure the path exists
		err := os.MkdirAll(path, 0700)
		require.NoError(t, err)
		path = filepath.Join(path, file)
	}
	c.DB.Driver = drivers.SQLite
	c.DB.DSN = path
	return c
}
