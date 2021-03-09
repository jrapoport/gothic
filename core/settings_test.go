package core

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/health"
	"github.com/jrapoport/gothic/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
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
	err := providers.LoadProviders(a.config)
	require.NoError(t, err)
	settings := a.Settings()
	assert.EqualValues(t, settings.Health, health.Check(a.config))
	assert.Equal(t, "", settings.Provider.Internal)
	has := settings.Provider.External[provider.Google]
	assert.True(t, has)
}
