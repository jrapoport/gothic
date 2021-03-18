package core

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/health"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	a.config.Mail.Host = "example.com"
	a.config.Mail.Port = 25
	a.config.Signup.Disabled = true
	a.config.UseInternal = false
	a.config.Providers = config.Providers{
		provider.Google: config.Provider{
			ClientKey:   "key",
			Secret:      "secret",
			CallbackURL: "http://exmaple.com",
		},
	}
	err := a.ext.LoadProviders(a.config)
	require.NoError(t, err)
	settings := a.Settings()
	assert.EqualValues(t, settings.Health, health.Check(a.config))
	assert.Equal(t, "", settings.Provider.Internal)
	has := settings.Provider.External[provider.Google]
	assert.True(t, has)
}
