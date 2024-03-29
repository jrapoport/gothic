package config

import (
	"testing"

	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	internal    = false
	redirectURL = "http://example.com/redirect"
	callback    = "http://example.com/callback"
	clientKey   = "foo"
	secret      = "i-am-a-secret"
)

func TestProviders(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		p := c.Authorization
		assert.Equal(t, internal, p.UseInternal)
		assert.Equal(t, redirectURL+test.mark, p.RedirectURL)
		names := []provider.Name{
			provider.Google,
			provider.GitLab,
		}
		assert.Len(t, p.Providers, len(names))
		for _, name := range names {
			v, ok := p.Providers[name]
			require.True(t, ok)
			assert.Equal(t, clientKey+test.mark, v.ClientKey)
			assert.Equal(t, secret+test.mark, v.Secret)
			assert.Equal(t, callback+test.mark, v.CallbackURL)
		}
	})
	a := Authorization{
		Providers: map[provider.Name]Provider{},
	}
	a.Providers[provider.Google] = Provider{
		ClientKey:   "foo",
		CallbackURL: "",
	}
	err := a.normalize(Service{}, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, a.Providers[provider.Google].CallbackURL)
	a.Providers[provider.Google] = Provider{
		ClientKey:   "foo",
		CallbackURL: "\n",
	}
	err = a.normalize(Service{}, "")
	assert.Error(t, err)
	a.RedirectURL = "\n"
	err = a.normalize(Service{}, "")
	assert.Error(t, err)
}

// tests the ENV vars are correctly taking precedence
func TestProviders_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			p := c.Authorization
			assert.Equal(t, internal, p.UseInternal)
			assert.Equal(t, redirectURL, p.RedirectURL)
			names := []provider.Name{
				provider.Google,
				provider.GitLab,
			}
			assert.Len(t, p.Providers, len(names))
			for _, name := range names {
				v, ok := p.Providers[name]
				require.True(t, ok)
				assert.Equal(t, clientKey, v.ClientKey)
				assert.Equal(t, secret, v.Secret)
				assert.Equal(t, callback, v.CallbackURL)
			}
		})
	}
}

// test the *un-normalized* defaults with load
func TestProviders_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := authorizationDefaults
	p := c.Authorization
	assert.Equal(t, def, p)
}

func TestAuthorization_InternalProvider(t *testing.T) {
	clearEnv()
	setEnv(t, ENVPrefix+"_SITE_URL", siteURL)
	setEnv(t, ENVPrefix+"_ROOT_PASSWORD", rootPassword)
	setEnv(t, ENVPrefix+"_JWT_SECRET", jwtSecret)
	setEnv(t, ENVPrefix+"_DB_DSN", dsn)
	c, err := loadNormalized("")
	assert.NoError(t, err)
	assert.EqualValues(t, c.Service.Name, c.Provider())
	c.Authorization.UseInternal = false
	assert.EqualValues(t, provider.Unknown, c.Provider())
}
