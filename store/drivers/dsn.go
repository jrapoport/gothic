package drivers

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
	"github.com/jrapoport/gothic/utils"
	"net/url"
	"path/filepath"
)

// NormalizeDSN returns a normalized dsn.
func NormalizeDSN(name string, drv Driver, dsn string) (string, string, error) {
	switch drv {
	case MySQL:
		return mysqlDSN(dsn)
	case Postgres:
		return postgresDSN(dsn)
	case SQLServer:
		return sqlserverDSN(dsn)
	case SQLite, SQLite3:
		return sqlitePath(name, dsn)
	}
	return name, dsn, nil
}

func mysqlDSN(dsn string) (string, string, error) {
	if dsn == "" {
		return "", "", errors.New("invalid mysql dsn")
	}
	dbc, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", "", err
	}
	return dbc.DBName, dbc.FormatDSN(), nil
}

func postgresDSN(dsn string) (string, string, error) {
	if dsn == "" {
		return "", "", errors.New("invalid postgres dsn")
	}
	c, err := pgx.ParseConfig(dsn)
	if err != nil {
		return "", "", err
	}
	return c.Database, c.ConnString(), nil
}

func sqlserverDSN(dsn string) (string, string, error) {
	if dsn == "" {
		return "", "", errors.New("invalid postgres dsn")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		return "", "", err
	}
	name := u.Query().Get("database")
	return name, dsn, nil
}

func sqlitePath(name, dsn string) (string, string, error) {
	const SQLiteExt = "." + SQLite
	if name == "" {
		const defaultName = "db"
		name = defaultName
	}
	if dsn == "" {
		return name, name + SQLiteExt, nil
	}
	file := name
	if !utils.HasExt(file) {
		file += SQLiteExt
	}
	u, err := url.Parse(dsn)
	if err != nil {
		return "", "", err
	}
	if u.Path == "" {
		u.Path = file
	} else if utils.IsDirectory(u.Path) || !utils.HasExt(u.Path) {
		u.Path = filepath.Join(u.Path, file)
	} else {
		_, name = filepath.Split(u.Path)
	}
	if ext := filepath.Ext(name); ext != "" {
		name = name[0 : len(name)-len(ext)]
	}
	return name, u.String(), nil
}
