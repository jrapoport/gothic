package tconn

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/require"
)

// TempConn creates a new Conn for tests with the configured test db.
func TempConn(t *testing.T) (*store.Connection, *config.Config) {
	c := tconf.TempDB(t)
	c.DB.AutoMigrate = true
	conn, err := store.Dial(c, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	err = conn.AutoMigrate()
	require.NoError(t, err)
	return conn, c
}
