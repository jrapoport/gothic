package log

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	levels := []Level{
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		PanicLevel,
		FatalLevel,
	}
	type NewLoggerFunc func(Level) Logger
	tests := []struct {
		name  string
		logFn NewLoggerFunc
	}{
		{"logrus", NewLogrusLoggerWithLevel},
		{"std", NewStdLoggerWithLevel},
		{"zap", NewZapLoggerWithLevel},
	}
	for _, level := range levels {
		for _, test := range tests {
			logger := test.logFn(level)
			// debug
			logger.Debug(test.name, " ", "debug")
			logger.Debugf("%s %s", test.name, "debug")
			// info
			logger.Info(test.name, " ", "info")
			logger.Infof("%s %s", test.name, "info")
			// warn
			logger.Warn(test.name, " ", "warn")
			logger.Warnf("%s %s", test.name, "warn")
			// error
			logger.Error(test.name, " ", "error")
			logger.Errorf("%s %s", test.name, "error")
			// panic
			assert.Panics(t, func() {
				logger.Panic(test.name, " ", "panic")
			})
			assert.Panics(t, func() {
				logger.Panicf("%s %s", test.name, "panic")
			})
		}
	}
}

func TestWithLevel(t *testing.T) {
	const (
		testName  = "test"
		emptyName = ""
	)
	tests := []struct {
		logger    interface{}
		name      string
		assertNil assert.ValueAssertionFunc
	}{
		{nil, emptyName, assert.Nil},
		{logrus.New(), emptyName, assert.NotNil},
		{logrus.New(), testName, assert.NotNil},
		{logrus.New().WithContext(nil), emptyName, assert.NotNil},
		{logrus.New().WithContext(nil), testName, assert.NotNil},
		{NewLogrusLoggerWithLevel(ErrorLevel), emptyName, assert.Nil},
		{NewLogrusLoggerWithLevel(ErrorLevel), testName, assert.Nil},
		{log.New(os.Stderr, "", 0), emptyName, assert.NotNil},
		{log.New(os.Stderr, "", 0), testName, assert.NotNil},
		{NewStdLoggerWithLevel(ErrorLevel), emptyName, assert.NotNil},
		{NewStdLoggerWithLevel(ErrorLevel), testName, assert.NotNil},
		{zap.NewExample(), emptyName, assert.NotNil},
		{zap.NewExample(), testName, assert.NotNil},
		{zap.NewExample().Sugar(), emptyName, assert.NotNil},
		{zap.NewExample().Sugar(), testName, assert.NotNil},
		{NewZapLoggerWithLevel(ErrorLevel), emptyName, assert.Nil},
		{NewZapLoggerWithLevel(ErrorLevel), testName, assert.Nil},
	}
	for _, test := range tests {
		logger := WithLevel(test.logger, ErrorLevel)
		require.NotNil(t, logger)
		logger.Error(testName)
	}
	lg := NewLogrusLoggerWithLevel(InfoLevel)
	require.NotNil(t, lg)
	lg = WithLevel(lg, 255)
	assert.NotNil(t, lg)
	lg = NewStdLoggerWithLevel(InfoLevel)
	require.NotNil(t, lg)
	lg = WithLevel(lg, 255)
	assert.NotNil(t, lg)
	lg = NewZapLoggerWithLevel(InfoLevel)
	require.NotNil(t, lg)
	lg = WithLevel(lg, 255)
	assert.NotNil(t, lg)
}

func TestLevelFromString(t *testing.T) {
	tests := []struct {
		in  string
		out Level
	}{
		{LevelDebug, DebugLevel},
		{LevelInfo, InfoLevel},
		{LevelWarn, WarnLevel},
		{LevelError, ErrorLevel},
		{LevelFatal, FatalLevel},
		{LevelPanic, PanicLevel},
		{"", PanicLevel},
	}
	for _, test := range tests {
		out := LevelFromString(test.in)
		assert.Equal(t, test.out, out)
	}
}
