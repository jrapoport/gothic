package log

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseFileOutput(t *testing.T) {
	const emptyName = ""
	randName := func() string { return uuid.New().String() }
	tests := []struct {
		newLogger func(level Level) Logger
		file      string
		name      string
	}{
		{NewLogrusLoggerWithLevel, emptyName, emptyName},
		{NewLogrusLoggerWithLevel, randName(), emptyName},
		{NewLogrusLoggerWithLevel, randName(), randName()},

		{NewStdLoggerWithLevel, emptyName, emptyName},
		{NewStdLoggerWithLevel, randName(), emptyName},
		{NewStdLoggerWithLevel, randName(), randName()},

		{NewZapLoggerWithLevel, emptyName, emptyName},
		{NewZapLoggerWithLevel, randName(), emptyName},
		{NewZapLoggerWithLevel, randName(), randName()},
	}

	for _, test := range tests {
		file := filepath.Join(t.TempDir(), test.file+".log")
		logger := test.newLogger(ErrorLevel).WithName(test.name)
		logger = logger.UseFileOutput(file)
		require.NotNil(t, logger)
		logger.Error("test")
		assert.FileExists(t, file)
		b, err := ioutil.ReadFile(file)
		require.NoError(t, err)
		assert.NotEmpty(t, b)
		logger = logger.UseFileOutput("")
		assert.NotNil(t, logger)
	}
}
