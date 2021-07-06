package log

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestWithLogger(t *testing.T) {
	fld := logrus.New().WithField("logger", "test")
	loggers := []logger.Interface{
		WithLogger(nil),
		WithLogger(fld),
	}
	for _, l := range loggers {
		levels := []logger.LogLevel{
			0,
			logger.Silent,
			logger.Error,
			logger.Warn,
			logger.Info,
		}
		for _, level := range levels {
			l.LogMode(level)
			l.Info(nil, "info %s", "test")
			l.Warn(nil, "warn %s", "test")
			l.Error(nil, "error %s", "test")
			l.Trace(nil, time.Now().UTC(), func() (string, int64) {
				return "test", 10
			}, errors.New("test"))
			l.Trace(nil, time.Now().UTC(), func() (string, int64) {
				return "test", 10
			}, gorm.ErrRecordNotFound)
			l.Trace(nil, time.Now().UTC(), func() (string, int64) {
				return "test", 10
			}, nil)
			tm := time.Now().UTC().Add(-1 * time.Minute)
			l.Trace(nil, tm, func() (string, int64) {
				return "test", 10
			}, nil)

		}
	}
}
