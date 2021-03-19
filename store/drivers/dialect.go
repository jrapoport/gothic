package drivers

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Dialect is a wrapper for the database dialect.
type Dialect struct {
	gorm.Dialector
	name string
}

// NewDialect returns a new configured db Dialect.
func NewDialect(ctx context.Context, drv Driver, dsn string) (*Dialect, error) {
	var d gorm.Dialector
	var cfg interface{}
	var name string
	var err error
	if ctx != nil {
		cfg = ConfigFromContext(ctx)
	}
	if cfg == nil {
		if dsn == "" {
			return nil, errors.New("dsn required")
		}
		name, dsn, err = NormalizeDSN("", drv, dsn)
		if err != nil {
			return nil, err
		}
	}
	switch drv {
	case MySQL:
		if dbc, ok := cfg.(mysql.Config); ok {
			d = mysql.New(dbc)
		} else {
			d = mysql.Open(dsn)
		}
	case Postgres:
		if dbc, ok := cfg.(postgres.Config); ok {
			d = postgres.New(dbc)
		} else {
			d = postgres.Open(dsn)
		}
	case SQLServer:
		if dbc, ok := cfg.(sqlserver.Config); ok {
			d = sqlserver.New(dbc)
		} else {
			d = sqlserver.Open(dsn)
		}
	case SQLite, SQLite3:
		d = sqlite.Open(dsn)
	default:
		db, err := sql.Open(string(drv), dsn)
		if err != nil {
			return nil, err
		}
		d = mysql.New(mysql.Config{Conn: db})
	}
	return &Dialect{d, name}, nil
}

// DBName returns the name of the db.
func (d Dialect) DBName() string {
	return d.name
}

type configKey struct{}

// ConfigWithContext adds a db config to the context
func ConfigWithContext(ctx context.Context, config interface{}) context.Context {
	return context.WithValue(ctx, configKey{}, config)
}

// ConfigFromContext gets a db config from the context
func ConfigFromContext(ctx context.Context) interface{} {
	return ctx.Value(configKey{})
}
