package test

import (
	"testing"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
	"github.com/stretchr/testify/assert"
)

// SetupDBConnection sets up a test DB connection, if the driver is sqlite, it opens
// the db in tmp directory which will be automatically wiped when the test exists.
func SetupDBConnection(t *testing.T, config *conf.Configuration) (*storage.Connection, error) {
	if config.DB.Driver == "sqlite" {
		config.DB.URL = t.TempDir()
	}
	conn, err := storage.Dial(config, nil)
	assert.NoError(t, err)
	if err != err {
		t.FailNow()
	} else if conn == nil {
		t.FailNow()
	}
	err = conn.DropDatabase()
	assert.NoError(t, err)
	return conn, err
}
