package tutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/require"
)

// ProjectRoot finds and returns and the project root for tests.
func ProjectRoot(t *testing.T) string {
	wd, _ := os.Getwd()
	for {
		root := filepath.Join(wd, ".git")
		if utils.PathExists(root) {
			break
		}
		wd = filepath.Dir(wd)
	}
	require.NotEmpty(t, wd)
	return wd
}
