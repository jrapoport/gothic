package storage

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/jrapoport/gothic/conf"
	"github.com/onrik/gorm-logrus"
	"github.com/pkg/errors"
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

// Dial will connect to that storage engine
func Dial(config *conf.GlobalConfiguration) (*Connection, error) {
	if config.DB.Driver == "" && config.DB.URL != "" {
		u, err := url.Parse(config.DB.URL)
		if err != nil {
			return nil, errors.Wrap(err, "parsing db connection url")
		}
		config.DB.Driver = u.Scheme
	}

	if config.DB.Database == "" {
		config.DB.Driver = "gothic"
		if config.DB.Namespace != "" {
			config.DB.Driver += "_" + config.DB.Namespace
		}
	}

	dialect := func() gorm.Dialector {
		switch config.DB.Driver {
		case "mysql":
			return mysql.Open(config.DB.URL)
		case "sqlserver":
			return sqlserver.Open(config.DB.URL)
		case "postgres":
			return postgres.New(postgres.Config{
				DSN:                  config.DB.URL,
				PreferSimpleProtocol: true,
			})
		case "sqlite":
			fallthrough
		default:
			u, _ := url.Parse(config.DB.URL)
			name := fmt.Sprintf("%s.sqlite", config.DB.Database)
			file := filepath.Join(u.Path, name)
			return sqlite.Open(file)
		}
	}

	namespace := func() string {
		if config.DB.Namespace != "" {
			return config.DB.Namespace + "_"
		}
		return ""
	}

	orm, err := gorm.Open(dialect(), &gorm.Config{
		Logger: gorm_logrus.New(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: namespace(), // table name prefix
		}})
	if err != nil {
		return nil, errors.Wrap(err, "opening database connection")
	}

	if logrus.StandardLogger().Level == logrus.DebugLevel {
		orm.Logger.LogMode(logger.Info)
	}

	if config.DB.Driver != "sqlite" && config.DB.Driver != "" {
		orm.Exec(fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS %s",
			config.DB.Database))
	}
	orm.Exec(fmt.Sprintf(
		"USE %s",
		config.DB.Database))

	conn := &Connection{DB: orm}
	conn = conn.withContext(context.Background(), config, uuid.Nil)

	if !config.DB.AutoMigrate {
		return conn, nil
	}

	if err = MigrateDatabase(conn); err != nil {
		return nil, errors.Wrap(err, "migrating database")
	}

	return conn, nil
}

func (c *Connection) Transaction(fn func(*Connection) error) error {
	if c.transaction {
		return fn(c)
	}
	return c.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&Connection{DB: tx, transaction: true})
	})
}
