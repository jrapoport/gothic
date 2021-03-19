package tconf

import (
	"path/filepath"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/jrapoport/gothic/utils"
)

const (
	// MockTemp is a mock sql db for tests.
	MockTemp = "mock"

	// MySQLTemp is a temp mysql container for tests.
	MySQLTemp = "mysql-temp"

	// MySQLDTemp is a temp mysql daemon for tests.
	MySQLDTemp = "mysqld-temp"

	// PostgresTemp is a temp postgres container for tests.
	PostgresTemp = "postgres-temp"

	// SQLServerTemp is a temp ms sql server container for tests.
	SQLServerTemp = "sqlserver-temp"

	// SQLiteTemp is a temp sqlite db for tests.
	SQLiteTemp = "sqlite-temp"

	// DBTemp is a temp db for tests (defaults to sqlite).
	DBTemp = SQLiteTemp
)

// DBConfig loads a configuration using driver for tests.
func DBConfig(t *testing.T, d drivers.Driver) *config.Config {
	c := Config(t)
	if d == "" {
		d = c.DB.Driver
	}
	return configDB(t, c, d)
}

func configDB(t *testing.T, c *config.Config, d drivers.Driver) *config.Config {
	c.DB.Namespace = "test"
	switch d {
	case MySQLTemp:
		c = mysqlDB(t, c)
	case MySQLDTemp:
		c = mysqldDB(t, c)
	case PostgresTemp:
		c = postgresDB(t, c)
	case SQLServerTemp:
		c = mssqlDB(t, c)
	case SQLiteTemp:
		_, file := filepath.Split(c.DB.DSN)
		if file == "" {
			file = "db"
		}
		if !utils.HasExt(file) {
			file += "." + drivers.SQLite
		}
		c.DB.DSN = filepath.Join(t.TempDir(), file)
		fallthrough
	case drivers.SQLite, drivers.SQLite3:
		c = sqliteDB(t, c)
	default:
		break
	}
	return c
}

// DBConfigs db configs for testing
func DBConfigs(t *testing.T, dvrs []drivers.Driver) []*config.Config {
	c := Config(t)
	cfgs := make([]*config.Config, len(dvrs))
	for i, d := range dvrs {
		cfgs[i] = &(*c)
		switch d {
		case SQLiteTemp,
			MySQLDTemp,
			drivers.SQLite,
			drivers.SQLite3:
			cfgs[i] = configDB(t, cfgs[i], d)
		default:
			cfgs[i].DB.Driver = d
		}
	}
	return gnomockDBs(t, cfgs)
}
