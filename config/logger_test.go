package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testLevel    = "debug"
	logFile      = "./logs/debug.log"
	logColors    = false
	logTimestamp = time.RFC1123Z
)

var testFields = map[string]string{
	"source":   "peaches",
	"priority": "1",
}

func TestLogger(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		l := c.Logger
		assert.Equal(t, testLevel+test.mark, l.Level)
		assert.Equal(t, logFile+test.mark, l.File)
		assert.Equal(t, logColors, l.Colors)
		assert.Equal(t, logTimestamp+test.mark, l.Timestamp)
		assert.Len(t, l.Fields, 2)
		fields := newKeyValueMap(l.Fields)
		for k, v := range fields {
			assert.Equal(t, testFields[k]+test.mark, v)
		}
	})
}

// tests the ENV vars are correctly taking precedence
func TestLogger_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			l := c.Logger
			assert.Equal(t, testLevel, l.Level)
			assert.Equal(t, logFile, l.File)
			assert.Equal(t, logColors, l.Colors)
			assert.Equal(t, logTimestamp, l.Timestamp)
			assert.Len(t, l.Fields, 2)
			fields := newKeyValueMap(l.Fields)
			for k, v := range fields {
				assert.Equal(t, testFields[k], v)
			}
		})
	}
}

// test the *un-normalized* defaults with load
func TestLogger_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	l := c.Logger
	def := loggerDefaults
	assert.Equal(t, def, l)
}

func TestLogger_NewLogger(t *testing.T) {
	l := loggerDefaults
	l.File = t.TempDir() + "test.log"
	l.Fields = []string{
		"source=peaches",
		"priority=1",
	}
	log := l.NewLogger()
	assert.NotNil(t, log)
	log.Debug("hello world")
}
