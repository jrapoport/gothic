package tutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// PathExists returns true if the path exits for tests.
func PathExists(t *testing.T, n string) bool {
	_, err := os.Stat(n)
	if os.IsNotExist(err) {
		return false
	}
	require.NoError(t, err)
	return true
}

// ProjectRoot finds and returns and the project root for tests.
func ProjectRoot(t *testing.T) string {
	wd, _ := os.Getwd()
	for {
		root := filepath.Join(wd, ".git")
		if PathExists(t, root) {
			break
		}
		wd = filepath.Dir(wd)
	}
	require.NotEmpty(t, wd)
	return wd
}
