package tconf

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
)

type mockKey struct{}

func mockWithContext(ctx context.Context, mock sqlmock.Sqlmock) context.Context {
	return context.WithValue(ctx, mockKey{}, mock)
}

// MockFromContext returns the mock from the context for tests
// The key is only valid when conf.Driver is MockTemp
func MockFromContext(ctx context.Context) sqlmock.Sqlmock {
	m, _ := ctx.Value(mockKey{}).(sqlmock.Sqlmock)
	return m
}

// MockDB loads a mock sql db for tests.
func MockDB(t *testing.T) (context.Context, *config.Config) {
	c := DBConfig(t, MockTemp)
	return mockDB(t, c)
}

func mockDB(t *testing.T, c *config.Config) (context.Context, *config.Config) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	require.NotNil(t, db)
	cfg := mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}
	ctx := context.Background()
	ctx = drivers.ConfigWithContext(ctx, cfg)
	ctx = mockWithContext(ctx, mock)
	c.DB.Driver = drivers.MySQL
	c.DB.DSN = ""
	return ctx, c
}
