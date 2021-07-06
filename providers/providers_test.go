package providers

import (
	"os"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T, key, value string) {
	err := os.Setenv(config.ENVPrefix+"_"+key, value)
	require.NoError(t, err)
}

func TestLoadProviders(t *testing.T) {
	c := tconf.ProvidersConfig(t)
	os.Clearenv()
	tests := []struct {
		key   string
		value string
		Err   assert.ErrorAssertionFunc
	}{
		{"", "", assert.Error},
		{config.Auth0DomainEnv, "example.com", assert.Error},
		{config.OpenIDConnectURLEnv, "http://example.com", assert.Error},
		{config.OpenIDConnectURLEnv, tconf.MockOpenIDConnect(t), assert.NoError},
	}
	for _, test := range tests {
		if test.key != "" {
			setEnv(t, test.key, test.value)
		}
		err := LoadProviders(c)
		test.Err(t, err)
	}
	err := LoadProviders(c)
	assert.NoError(t, err)
}

func TestUseProvider(t *testing.T) {
	c := tconf.ProvidersConfig(t)
	setEnv(t, config.TwitterAuthorizeEnv, "1")
	for name, v := range c.Providers {
		err := useProvider(name, v.ClientKey, v.Secret, v.CallbackURL, v.Scopes...)
		assert.NoError(t, err, name)
	}
	err := useProvider("", "", "", "")
	assert.Error(t, err)
	err = useProvider("unknown", "", "", "")
	assert.Error(t, err)
}

func TestIsEnabled(t *testing.T) {
	const badProvider = "bad"
	clearProviders()
	tests := make([]provider.Name, len(provider.External))
	var i = 0
	for name := range provider.External {
		tests[i] = name
		i++
	}
	tests = append(tests, provider.Unknown, badProvider)
	for _, name := range tests {
		err := IsEnabled(name)
		assert.Error(t, err, name)
	}
	c := tconf.ProvidersConfig(t)
	err := LoadProviders(c)
	assert.NoError(t, err)
	tests = append(tests, c.Provider())
	for _, name := range tests {
		err = IsEnabled(name)
		switch name {
		case provider.Unknown, badProvider:
			assert.Error(t, err, name)
		default:
			assert.NoError(t, err, name)
		}
	}
}
