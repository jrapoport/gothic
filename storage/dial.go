package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/aklinkert/go-gorm-logrus-logger"
	"github.com/cenkalti/backoff/v4"
	"github.com/jrapoport/gothic/conf"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Connection is the interface a storage provider must implement.
type Connection struct {
	*gorm.DB
	transaction bool
}

// Dial will connect to the database.
func Dial(c *conf.Configuration, l *logrus.Entry) (*Connection, error) {
	d, err := dialect(c)
	if err != nil {
		err = fmt.Errorf("failed to parse database driver %w", err)
		return nil, err
	}
	ns := namespace(c)
	if l == nil {
		l = logrus.New().WithContext(nil)
	}
	log := gormlogruslogger.NewGormLogrusLogger(l, 100*time.Millisecond)
	log = log.LogMode(logLevel(l.Level))
	dbc := &gorm.Config{
		Logger:                                   log,
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
	if c.DB.Driver != "sqlite" && c.DB.Driver != "" {
		create := "CREATE DATABASE IF NOT EXISTS " + c.DB.Name
		use := "USE " + c.DB.Name
		db.Exec(create)
		db.Exec(use)
	}
	conn := &Connection{DB: db}
	conn = conn.withContext(context.Background(), c)
	if !c.DB.AutoMigrate {
		return conn, nil
	}
	if err = MigrateDatabase(conn); err != nil {
		err = fmt.Errorf("%w migrating database", err)
		return nil, err
	}
	return conn, nil
}

// Transaction opens a database transaction and prevents nested transactions.
func (c *Connection) Transaction(fn func(*Connection) error) error {
	if c.transaction {
		return fn(c)
	}
	return c.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&Connection{DB: tx, transaction: true})
	})
}

func name(c *conf.Configuration) string {
	n := c.DB.Name
	if n == "" {
		n = "gothic"
		if c.DB.Namespace != "" {
			n += "_" + c.DB.Namespace
		}
	}
	return n
}

func namespace(c *conf.Configuration) string {
	if c.DB.Namespace != "" {
		return c.DB.Namespace + "_"
	}
	return ""
}

func driver(c *conf.Configuration) (string, error) {
	dvr := c.DB.Driver
	if dvr == "" && c.DB.URL != "" {
		u, err := url.Parse(c.DB.URL)
		if err != nil {
			err = fmt.Errorf("%w parsing db connection url", err)
			return "", err
		}
		dvr = u.Scheme
	}
	return dvr, nil
}

func dialect(c *conf.Configuration) (gorm.Dialector, error) {
	dvr, err := driver(c)
	if err != nil {
		return nil, err
	}
	switch dvr {
	case "mysql":
		return mysql.Open(c.DB.URL), nil
	case "sqlserver":
		return sqlserver.Open(c.DB.URL), nil
	case "postgres":
		return postgres.New(postgres.Config{
			DSN:                  c.DB.URL,
			PreferSimpleProtocol: true,
		}), nil
	case "sqlite":
		fallthrough
	default:
		u, _ := url.Parse(c.DB.URL)
		fn := fmt.Sprintf("%s.sqlite", name(c))
		file := filepath.Join(u.Path, fn)
		return sqlite.Open(file), nil
	}
}

func disableFKey(c *conf.Configuration) bool {
	switch c.DB.Driver {
	case "mysql", "sqlserver", "postgres":
		return false
	default:
		return true
	}
}

func maxRetries(c *conf.Configuration) uint64 {
	return uint64(c.DB.MaxRetries)
}

func logLevel(l logrus.Level) logger.LogLevel {
	lvl := logger.Silent
	switch l {
	case logrus.DebugLevel, logrus.TraceLevel:
		lvl = logger.Info
	case logrus.WarnLevel:
		lvl = logger.Warn
	case logrus.ErrorLevel:
		lvl = logger.Error
	}
	return lvl
}
