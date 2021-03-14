package tconn

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/require"
)

// Conn creates a new Conn for tests with the configured test db.
func Conn(t *testing.T, c *config.Config) *store.Connection {
	c.DB.AutoMigrate = true
	conn, err := store.Dial(c, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	err = conn.AutoMigrate()
	require.NoError(t, err)
	return conn
}

// TempConn creates a new Conn for tests with the configured test db.
func TempConn(t *testing.T) (*store.Connection, *config.Config) {
	c := tconf.TempDB(t)
	return Conn(t, c), c
}
