package log

import (
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithName(t *testing.T) {
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
		logger := named(test.logger, "name", ErrorLevel)
		test.assertNil(t, logger)
		if logger != nil {
			_, ok := logger.(Logger)
			assert.True(t, ok)
			// error
			logger.Error(testName)
			logger.Errorf("%s", testName)
		}
	}
	lg := NewLogrusLoggerWithLevel(255).WithName("name")
	assert.NotNil(t, lg)
	lg = NewLogrusLoggerWithLevel(255).WithName("name")
	assert.NotNil(t, lg)
	lg = NewStdLoggerWithLevel(255).WithName("name")
	assert.NotNil(t, lg)
	lg = NewStdLoggerWithLevel(255).WithName("name")
	assert.NotNil(t, lg)
	lg = NewZapLoggerWithLevel(255).WithName("name")
	assert.NotNil(t, lg)
	lg = NewZapLoggerWithLevel(255).WithName("name")
	assert.NotNil(t, lg)
}
