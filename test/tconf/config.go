package tconf

import (
	"path/filepath"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/require"
)

// Config loads a configuration for tests.
func Config(t *testing.T) *config.Config {
	path := configPath(t)
	c, err := config.LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, c)
	return configDB(t, c, c.DB.Driver)
}

func configPath(t *testing.T) string {
	const configPath = "test/test.env"
	root := tutils.ProjectRoot(t)
	path := filepath.Join(root, configPath)
	require.FileExists(t, path)
	return path
}
