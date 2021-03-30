package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/jrapoport/gothic/store/log"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Connection is the interface a storage provider must implement.
type Connection struct {
	*gorm.DB
}

// Dial will connect to the database.
func Dial(c *config.Config, l logrus.FieldLogger) (*Connection, error) {
	conn, err := NewConnection(context.Background(), c, l)
	if err != nil {
		return nil, err
	}
	if c.DB.AutoMigrate {
		err = conn.AutoMigrate()
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

// NewConnection returns a new db connection.
func NewConnection(ctx context.Context, c *config.Config, l logrus.FieldLogger) (*Connection, error) {
	if c == nil {
		return nil, errors.New("configuration required")
	}
	l = ensureLog(c, l)
	d, err := drivers.NewDialect(ctx, c.DB.Driver, c.DB.DSN)
	if err != nil {
		return nil, err
	}
	l = l.WithField("db", d.DBName())
	ns := utils.Namespaced(c.DB.Namespace, "")
	dbc := &gorm.Config{
		Logger:                   log.WithLogger(l),
		DisableNestedTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: ns, // table name prefix
		},
	}
	var db *gorm.DB
	max := uint64(c.DB.MaxRetries)
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), max)
	err = backoff.RetryNotify(
		func() error {
			db, err = gorm.Open(d, dbc)
			return err
		},
		retry,
		func(err error, duration time.Duration) {
			if l != nil && max > 0 {
				l.WithError(err).
					Warn("database connection failed")
				l.WithField("duration", duration).
					Info("retrying in...")
			}
		},
	)
	if err != nil {
		err = fmt.Errorf("error opening database: %w", err)
		return nil, err
	}
	return &Connection{db}, nil
}

func ensureLog(c *config.Config, l logrus.FieldLogger) logrus.FieldLogger {
	if l == nil {
		l = c.Log()
	}
	if l == nil {
		l = logrus.New()
	}
	return l
}

// Transaction opens a database transaction.
func (conn *Connection) Transaction(fn func(tx *Connection) error) error {
	return conn.DB.Transaction(func(gtx *gorm.DB) error {
		return fn(&Connection{gtx})
	})
}

// RunMigration runs the migration m.
func (conn *Connection) RunMigration(m *migration.Migration) error {
	return m.Run(conn.DB)
}

// Migrate runs the migration plan.
func (conn *Connection) Migrate(p *migration.Plan) error {
	return p.Run(conn.DB, true)
}

// AutoMigrate runs the global migration plan if AutoMigrate is true.
func (conn *Connection) AutoMigrate() error {
	return conn.Migrate(plan)
}

// DBName returns the current database name.
func (conn *Connection) DBName() string {
	return conn.Migrator().CurrentDatabase()
}

// TableNames returns the table names for the models.
func (conn *Connection) TableNames() ([]string, error) {
	database := conn.DBName()
	var tx *gorm.DB
	switch conn.Name() {
	case drivers.SQLite:
		tx = conn.Table("sqlite_master").
			Select("tbl_name").
			Where("type = ?", "table").
			Where("tbl_name NOT LIKE ?", "sqlite_%")
	case drivers.Postgres:
		database = "public"
		fallthrough
	default:
		tx = conn.Table("information_schema.tables").
			Select("table_name").
			Where("table_type = ?", "BASE TABLE").
			Where("table_schema = ?", database)
	}
	var names []string
	err := tx.Scan(&names).Error
	if err != nil {
		return nil, err
	}
	return names, nil
}

// TruncateAll truncates all tables in the database.
func (conn *Connection) TruncateAll() error {
	return conn.Transaction(func(tx *Connection) error {
		names, err := tx.TableNames()
		if err != nil {
			return err
		}
		if tx.Name() == drivers.MySQL {
			if err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
				return err
			}
			defer tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
		}
		for _, name := range names {
			raw := "TRUNCATE TABLE " + name
			if tx.Name() == drivers.SQLite {
				raw = "DELETE FROM " + name
			}
			if err = tx.Exec(raw).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DropAll drops all tables in the database.
func (conn *Connection) DropAll() error {
	return conn.Transaction(func(tx *Connection) error {
		names, err := tx.TableNames()
		if err != nil {
			return err
		}
		for _, name := range names {
			raw := "DROP TABLE " + name
			if err = tx.Exec(raw).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Has returns true if the store contains the object, otherwise false.
// If an error besides not found occurs, false and the error are returned.
func (conn *Connection) Has(v interface{}, c ...interface{}) (bool, error) {
	return Has(conn.DB, v, c...)
}

// HasLast returns true if the store contains the object, otherwise false.
// If an error besides not found occurs, false and the error are returned.
func (conn *Connection) HasLast(v interface{}, c ...interface{}) (bool, error) {
	return HasLast(conn.DB, v, c...)
}

/*
func (conn *Connection) withContext(ctx context.Context, l logrus.FieldLogger) *Connection {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.SetValue(ctx, interface{}(types.LoggerKey), l)
	conn.DB = conn.DB.WithContext(ctx)
	return conn
}

// Log gets the log was used to initialize the database context.
func (conn *Connection) Log() logrus.FieldLogger {
	return conn.DB.Statement.Context.Value(types.LoggerKey).(logrus.FieldLogger)
}
*/
