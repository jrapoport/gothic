package log

import (
	"context"
	"fmt"
	"time"

	"github.com/jrapoport/gothic/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Logger wraps a logger.Interface
type Logger struct {
	log.Logger
}

var _ logger.Interface = (*Logger)(nil)

// WithLogger takes a log.Logger and returns a
// *Logger that conforms to the gorm logger.Interface.
func WithLogger(l log.Logger) *Logger {
	return &Logger{l}
}

// LogMode sets the log level for the logger.
func (g *Logger) LogMode(level logger.LogLevel) logger.Interface {
	lvl := logLevel(level)
	g.Logger = log.WithLevel(g.Logger, lvl)
	return g
}

// Info logs an formatted string with info level.
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
		err = fmt.Errorf("%d %f: %w", rows, duration, err)
		g.Error(nil, err.Error())
		g.Warn(nil, sql)
	case elapsed > 100*time.Millisecond:
		g.Warnf("%d %f: %s", rows, duration, sql)
	default:
		g.Debugf("%d %f: %s", rows, duration, sql)
	}
}

func logLevel(l logger.LogLevel) log.Level {
	switch l {
	case logger.Silent:
		return log.FatalLevel
	case logger.Error:
		return log.ErrorLevel
	case logger.Warn:
		return log.WarnLevel
	case logger.Info:
		return log.DebugLevel
	}
	return log.FatalLevel
}
