package tdb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/data-dog/go-sqlmock.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB returns a temp SQLite *gorm.DB for tests.
func DB(t *testing.T) *gorm.DB {
	f := filepath.Join(t.TempDir(), "test.db")
	d := sqlite.Open(f)
	db, err := gorm.Open(d, nil)
	require.NoError(t, err)
	return db
}

// MockDB returns a mocked db instance
func MockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	msql, mock, err := sqlmock.New()
	require.NoError(t, err)
	cfg := mysql.New(mysql.Config{
		Conn:                      msql,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(cfg, &gorm.Config{})
	require.NoError(t, err)
	t.Cleanup(func() {
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	return db, mock
}
