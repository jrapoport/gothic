package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jrapoport/gothic/conf"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Dial will connect to the database.
func Dial(c *conf.Configuration, l *logrus.Entry) (*Connection, error) {
	d, err := dialect(c)
	if err != nil {
		err = fmt.Errorf("failed to parse database driver %w", err)
		return nil, err
	}
	ns := Namespace(c)
	dbc := &gorm.Config{
		Logger:                                   withLogger(l),
		DisableForeignKeyConstraintWhenMigrating: disableFKey(c),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: ns, // table name prefix
		},
	}
	var db *gorm.DB
	try := func() error {
		db, err = gorm.Open(d, dbc)
		return err
	}
	retries := maxRetries(c)
	err = backoff.RetryNotify(try,
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retries),
		func(err error, duration time.Duration) {
			if l != nil && retries > 0 {
				l.WithError(err).Warn("database connection failed")
				l.WithField("duration", duration).Info("retrying in...")
			}
		},
	)
	if err != nil {
		err = fmt.Errorf("error opening database: %w", err)
		return nil, err
	}
	if db == nil {
		return nil, errors.New("failed to open database")
	}
	conn := &Connection{DB: db}
	conn = conn.withConfiguredContext(c)
	if !c.DB.AutoMigrate {
		return conn, nil
	}
	if err = conn.MigrateDatabase(); err != nil {
		err = fmt.Errorf("%w migrating database", err)
		return nil, err
	}
	return conn, nil
}
