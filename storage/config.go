package storage

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/aklinkert/go-gorm-logrus-logger"
	"github.com/jrapoport/gothic/conf"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Namespace(c *conf.Configuration) string {
	if c.DB.Namespace != "" {
		return c.DB.Namespace + "_"
	}
	return ""
}

func dbName(c *conf.Configuration) string {
	n := c.DB.Name
	if n == "" {
		n = "gothic"
	}
	return Namespace(c) + n
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
		fn := fmt.Sprintf("%s.sqlite", dbName(c))
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

func withLogger(l *logrus.Entry) logger.Interface {
	if l == nil {
		l = logrus.New().WithContext(nil)
	}
	gl := gormlogruslogger.NewGormLogrusLogger(l, 100*time.Millisecond)
	gl = gl.LogMode(logLevel(l.Level))
	return gl
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
