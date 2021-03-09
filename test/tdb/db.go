package tdb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
