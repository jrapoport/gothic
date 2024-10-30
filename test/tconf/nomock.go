package tconf

import (
	"fmt"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mssql"
	"github.com/orlangure/gnomock/preset/mysql"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"  // mysql driver
	_ "github.com/lib/pq"               // postgres driver
	_ "github.com/mattn/go-sqlite3"     // sqlite driver
	_ "github.com/microsoft/go-mssqldb" // mssql driver
)

const (
	username = "gothic"
	password = "password"
)

func mysqlPreset(c *config.Config) (string, gnomock.Preset) {
	const mysqlDSN = "%s:%s@tcp(%%s:%%d)/%s?parseTime=true&multiStatements=true"
	database := c.DB.Namespace + c.Name
	dsnFormat := fmt.Sprintf(mysqlDSN, username, password, database)
	p := mysql.Preset(
		mysql.WithDatabase(database),
		mysql.WithUser(username, password),
	)
	return dsnFormat, p
}

func mysqlDB(t *testing.T, c *config.Config) *config.Config {
	dsnFormat, p := mysqlPreset(c)
	c.DB.DSN = gnomockDB(t, p, dsnFormat)
	c.DB.Driver = drivers.MySQL
	return c
}

func postgresPreset(c *config.Config) (string, gnomock.Preset) {
	const postgresDSN = "host=%%s port=%%d user=%s password=%s dbname=%s sslmode=disable"
	database := c.DB.Namespace + c.Name
	dsnFormat := fmt.Sprintf(postgresDSN, username, password, database)
	p := postgres.Preset(
		postgres.WithDatabase(database),
		postgres.WithUser(username, password),
	)
	return dsnFormat, p
}

func postgresDB(t *testing.T, c *config.Config) *config.Config {
	dsnFormat, p := postgresPreset(c)
	c.DB.DSN = gnomockDB(t, p, dsnFormat)
	c.DB.Driver = drivers.Postgres
	return c
}

func mssqlPreset(c *config.Config) (string, gnomock.Preset) {
	const mssqlDSN = "sqlserver://sa:%s@%%s%%d?database=%s"
	database := c.DB.Namespace + c.Name
	dsnFormat := fmt.Sprintf(mssqlDSN, password, database)
	p := mssql.Preset(
		mssql.WithDatabase(database),
		mssql.WithAdminPassword(password),
		mssql.WithLicense(true),
	)
	return dsnFormat, p
}

func mssqlDB(t *testing.T, c *config.Config) *config.Config {
	dsnFormat, p := mssqlPreset(c)
	c.DB.DSN = gnomockDB(t, p, dsnFormat)
	c.DB.Driver = drivers.SQLServer
	return c
}

func gnomockDB(t *testing.T, p gnomock.Preset, dsnFormat string) string {
	dock, err := gnomock.Start(p)
	require.NoError(t, err)
	require.NotNil(t, dock)
	t.Cleanup(func() {
		err = gnomock.Stop(dock)
		require.NoError(t, err)
	})
	dsn := fmt.Sprintf(dsnFormat, dock.Host, dock.DefaultPort())
	return dsn
}

func gnomockDBs(t *testing.T, cfgs []*config.Config) []*config.Config {
	para := gnomock.InParallel()
	type index struct {
		dock int
	}
	idxs := make([]index, len(cfgs))
	var x int
	for i, c := range cfgs {
		var p gnomock.Preset
		switch c.DB.Driver {
		case MySQLTemp:
			c.DB.DSN, p = mysqlPreset(c)
			para = para.Start(p)
		case PostgresTemp:
			c.DB.DSN, p = postgresPreset(c)
			para = para.Start(p)
		case SQLServerTemp:
			c.DB.DSN, p = mssqlPreset(c)
			para = para.Start(p)
		default:
			continue
		}
		para = para.Start(p)
		idxs[i].dock = x
		x++
	}
	containers, err := para.Go()
	require.NoError(t, err)
	for i, c := range cfgs {
		var d drivers.Driver
		switch c.DB.Driver {
		case MySQLTemp:
			d = drivers.MySQL
			break
		case PostgresTemp:
			d = drivers.Postgres
			break
		case SQLServerTemp:
			d = drivers.SQLServer
			break
		default:
			continue
		}
		cont := containers[idxs[i].dock]
		c.DB.Driver = d
		c.DB.DSN = fmt.Sprintf(c.DB.DSN, cont.Host, cont.DefaultPort())
	}
	t.Cleanup(func() {
		err = gnomock.Stop(containers...)
		require.NoError(t, err)
	})
	return cfgs
}
