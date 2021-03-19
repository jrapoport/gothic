package drivers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"

	_ "gorm.io/driver/clickhouse" // for tests
)

func TestDialect(t *testing.T) {
	const (
		myDSN    = "root@tcp(0.0.0.0:3306)/test?parseTime=true"
		pgDSN    = "postgres://root:password@0.0.0.0:5432/test"
		msDSN    = "sqlserver://sa:password@0.0.0.0:5432?database=test"
		clickDSN = "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
	)
	myCtx := ConfigWithContext(context.Background(), mysql.Config{})
	pgCtx := ConfigWithContext(context.Background(), postgres.Config{})
	msCtx := ConfigWithContext(context.Background(), sqlserver.Config{})
	tests := []struct {
		drv  Driver
		ctx  context.Context
		dsn  string
		name string
		Err  assert.ErrorAssertionFunc
	}{
		{MySQL, nil, "", "", assert.Error},
		{MySQL, nil, "\n", "", assert.Error},
		{MySQL, nil, myDSN, "test", assert.NoError},
		{MySQL, myCtx, "", "", assert.NoError},
		{Postgres, nil, pgDSN, "test", assert.NoError},
		{Postgres, pgCtx, "", "", assert.NoError},
		{SQLServer, nil, msDSN, "test", assert.NoError},
		{SQLServer, msCtx, "", "", assert.NoError},
		{SQLite, nil, "test.db", "test", assert.NoError},
		{"clickhouse", nil, clickDSN, "", assert.NoError},
	}
	for _, test := range tests {
		d, err := NewDialect(test.ctx, test.drv, test.dsn)
		test.Err(t, err)
		if err == nil {
			assert.Equal(t, test.name, d.DBName())
		}
	}
}
