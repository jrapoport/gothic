package tconf

import (
	"fmt"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mysql"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
)

func mysqlDB(t *testing.T, c *config.Config) *config.Config {
	const (
		username = "gothic"
		password = "password"
	)
	database := c.DB.Namespace + c.Name
	dsnFormat := fmt.Sprintf(
		"%s:%s@tcp(%%s:%%d)/%s?parseTime=true&multiStatements=true",
		username, password, database)
	usrOpt := mysql.WithUser(username, password)
	dbOpt := mysql.WithDatabase(database)
	p := mysql.Preset(usrOpt, dbOpt)
	c.DB.DSN = gnomockDB(t, p, dsnFormat)
	c.DB.Driver = drivers.MySQL
	return c
}

func postgresDB(t *testing.T, c *config.Config) *config.Config {
	const (
		username = "gothic"
		password = "password"
	)
	database := c.DB.Namespace + c.Name
	dsnFormat := fmt.Sprintf(
		"host=%%s port=%%d user=%s password=%s dbname=%s sslmode=disable",
		username, password, database)
	usrOpt := postgres.WithUser(username, password)
	dbOpt := postgres.WithDatabase(database)
	p := postgres.Preset(usrOpt, dbOpt)
	c.DB.DSN = gnomockDB(t, p, dsnFormat)
	c.DB.Driver = drivers.Postgres
	return c
}

func gnomockDB(t *testing.T, p gnomock.Preset, dsnFormat string) string {
	dock, err := gnomock.Start(p)
	require.NoError(t, err)
	require.NotNil(t, dock)
	t.Cleanup(func() {
		_ = gnomock.Stop(dock)
	})
	dsn := fmt.Sprintf(dsnFormat, dock.Host, dock.DefaultPort())
	return dsn
}
