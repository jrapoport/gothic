package log

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseFileOutput(t *testing.T) {
	const (
		emptyName = ""
		testName  = "test.log"
	)
	tests := []struct {
		newLogger func(level Level) Logger

		name string
	}{
		{NewLogrusLoggerWithLevel, emptyName},
		{NewLogrusLoggerWithLevel, testName},

		{NewStdLoggerWithLevel, emptyName},
		{NewStdLoggerWithLevel, testName},

		{NewZapLoggerWithLevel, emptyName},
		{NewZapLoggerWithLevel, testName},
	}

	for i, test := range tests {
		file := fmt.Sprintf("%d-%s", i, testName)
		file = filepath.Join(t.TempDir(), file)
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
