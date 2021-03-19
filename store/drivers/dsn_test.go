package drivers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDSN(t *testing.T) {
	const (
		myDSN = "root@tcp(0.0.0.0:3306)/test?parseTime=true"
		pgDSN = "postgres://root:password@0.0.0.0:5432/test"
		msDSN = "sqlserver://sa:password@0.0.0.0:5432?database=test"
	)
	var dir = t.TempDir()
	tests := []struct {
		drv     Driver
		name    string
		dsn     string
		nameOut string
		dsnOut  string
		Err     assert.ErrorAssertionFunc
	}{
		{
			MySQL, "", "",
			"", "", assert.Error,
		},
		{
			MySQL, "", "\n",
			"", "", assert.Error,
		},
		{
			MySQL, "", myDSN,
			"test", myDSN, assert.NoError,
		},
		{
			Postgres, "", "",
			"", "", assert.Error,
		},
		{
			Postgres, "", "\n",
			"", "", assert.Error,
		},
		{
			Postgres, "", pgDSN,
			"test", pgDSN, assert.NoError,
		},
		{
			SQLServer, "", "",
			"", "", assert.Error,
		},
		{
			SQLServer, "", "\n",
			"", "", assert.Error,
		},
		{
			SQLServer, "", msDSN,
			"test", msDSN, assert.NoError,
		},
		{
			SQLite, "", "", "db",
			"db.sqlite", assert.NoError,
		},
		{
			SQLite, "", "\n", "",
			"", assert.Error,
		},
		{
			SQLite, "", dir, "db",
			dir + "/db.sqlite", assert.NoError,
		},
		{
			SQLite, "", "db", "db",
			"db/db.sqlite", assert.NoError,
		},
		{
			SQLite, "", "db.sqlite",
			"db", "db.sqlite", assert.NoError,
		},
		{
			SQLite, "", dir + "/db.sqlite",
			"db", dir + "/db.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "", "test",
			"test.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "\n", "",
			"", assert.Error,
		},
		{
			SQLite, "test", dir, "test",
			dir + "/test.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "db", "test",
			"db/test.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "db.sqlite",
			"db", "db.sqlite", assert.NoError,
		},
		{
			SQLite, "test", dir + "/db.sqlite",
			"db", dir + "/db.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "foo/db/",
			"test", "foo/db/test.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "http://",
			"test", "http://test.sqlite", assert.NoError,
		},
		{
			SQLite, "test", "http://foo.",
			"test", "http://foo./test.sqlite", assert.NoError,
		},
		{
			"unkn", "test", "test",
			"test", "test", assert.NoError,
		},
	}
	for _, test := range tests {
		name, dsn, err := NormalizeDSN(test.name, test.drv, test.dsn)
		test.Err(t, err)
		assert.Equal(t, test.nameOut, name)
		assert.Equal(t, test.dsnOut, dsn)
	}
}
