package test

import (
	"testing"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
	"github.com/stretchr/testify/assert"
)

func SetupDBConnection(t *testing.T, globalConfig *conf.Configuration) (*storage.Connection, error) {
	globalConfig.DB.Driver = "sqlite"
	globalConfig.DB.URL = t.TempDir()
	conn, err := storage.Dial(globalConfig)
	assert.NoError(t, err)
	if err != err {
		t.FailNow()
	}
	err = storage.TruncateAll(conn)
	assert.NoError(t, err)
	return conn, err
}
