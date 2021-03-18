package log

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	logrus.FieldLogger
}

var _ logger.Interface = (*Logger)(nil)

// WithLogger takes a logrus.FieldLogger and returns a
// *Logger that conforms to the gorm logger.Interface.
func WithLogger(l logrus.FieldLogger) *Logger {
	if l == nil {
		l = logrus.New()
	}
	return &Logger{l}
}

// LogMode sets the log level for the logger.
func (g Logger) LogMode(level logger.LogLevel) logger.Interface {
	lvl := logLevel(level)
	switch l := g.FieldLogger.(type) {
	case *logrus.Entry:
		l.Level = lvl
	case *logrus.Logger:
		l.Level = lvl
	}
	return g
}

// Warn logs an formatted string with info level.
func (g Logger) Info(_ context.Context, s string, i ...interface{}) {
	g.Debugf(s, i...)
}

// Warn logs an formatted string with warn level.
func (g Logger) Warn(_ context.Context, s string, i ...interface{}) {
	g.Warnf(s, i...)
}

// Error logs an formatted string with error level.
func (g Logger) Error(_ context.Context, s string, i ...interface{}) {
	g.Errorf(s, i...)
}

// Trace logs an trace formatted string.
func (g Logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	duration := float64(elapsed.Nanoseconds()) / 1e6
	switch {
	case err != nil:
		// record not found is an expected error and thus not logged
		if err == gorm.ErrRecordNotFound {
			return
		}
		g.WithFields(logrus.Fields{
			"error":    err,
			"rows":     rows,
			"duration": duration,
		}).Warn(sql)

	case elapsed > 100*time.Millisecond:
		g.WithFields(logrus.Fields{
			"rows":     rows,
			"duration": duration,
		}).Warn(sql)

	default:
		g.WithFields(logrus.Fields{
			"rows":     rows,
			"duration": duration,
		}).Debug(sql)
	}
}

func logLevel(l logger.LogLevel) logrus.Level {
	switch l {
	case logger.Silent:
		return logrus.FatalLevel
	case logger.Error:
		return logrus.ErrorLevel
	case logger.Warn:
		return logrus.WarnLevel
	case logger.Info:
		return logrus.DebugLevel
	}
	return logrus.FatalLevel
}
