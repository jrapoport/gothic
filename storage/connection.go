package storage

import (
	"context"

	"github.com/jrapoport/gothic/conf"
	"gorm.io/gorm"
)

const configKey = "config"

// Connection is the interface a storage provider must implement.
type Connection struct {
	*gorm.DB
	transaction bool
}

func (conn *Connection) withConfiguredContext(c *conf.Configuration) *Connection {
	ctx := context.WithValue(context.Background(), interface{}(configKey), c)
	return &Connection{DB: conn.DB.WithContext(ctx)}
}

// Config gets configuration that was used to initialize the database context
func (conn *Connection) Config() *conf.Configuration {
	if c, ok := conn.DB.Statement.Context.Value(configKey).(*conf.Configuration); ok {
		return c
	}
	return nil
}

// Transaction opens a database transaction and prevents nested transactions.
func (conn *Connection) Transaction(fn func(*Connection) error) error {
	if conn.transaction {
		return fn(conn)
	}
	return conn.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&Connection{tx, true})
	})
}

func (conn *Connection) MigrateDatabase() error {
	c := conn.Config()
	name := dbName(c)
	if conn.Name() != "sqlite" {
		create := "CREATE DATABASE IF NOT EXISTS " + name
		if err := conn.Exec(create).Error; err != nil {
			return err
		}
		use := "USE " + name
		if err := conn.Exec(use).Error; err != nil {
			return err
		}
	}
	return conn.AutoMigrate(migrations...)
}

func (conn *Connection) DropDatabase() error {
	var err error
	if conn.Name() == "sqlite" {
		err = DropAll(conn)
	} else {
		c := conn.Config()
		drop := "DROP SCHEMA IF EXISTS " + dbName(c)
		err = conn.Exec(drop).Error
	}
	if err != nil {
		return err
	}
	return conn.MigrateDatabase()
}

func TruncateAll(conn *Connection) error {
	return conn.Transaction(func(tx *Connection) error {
		names, err := tx.tableNames()
		if err != nil {
			return err
		}
		if tx.Name() != "sqlite" {
			if err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
				return err
			}
			defer tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
		}
		for _, name := range names {
			raw := "TRUNCATE TABLE " + name
			if tx.Name() == "sqlite" {
				raw = "DELETE FROM " + name
			}
			if err = tx.Exec(raw).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func DropAll(conn *Connection) error {
	return conn.Transaction(func(tx *Connection) error {
		names, err := tx.tableNames()
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

func (conn *Connection) tableNames() ([]string, error) {
	var names []string
	for _, m := range migrations {
		stmt := conn.Statement
		if err := stmt.Parse(m); err != nil {
			return nil, err
		}
		names = append(names, stmt.Schema.Table)
	}
	return names, nil
}
