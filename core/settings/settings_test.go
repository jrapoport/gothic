package settings

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/health"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	c := tconf.Config(t)
	c.Mail.Host = "example.com"
	c.Mail.Port = 25
	c.Signup.Disabled = true
	p := config.Provider{
		ClientKey:   "key",
		Secret:      "secret",
		CallbackURL: "http://example.com",
	}
	c.UseInternal = true
	c.Providers[provider.Google] = p
	c.Providers[provider.GitLab] = p
	c.Providers[provider.Heroku] = p
	s := Current(c)
	assert.EqualValues(t, s.Health, health.Check(c))
	assert.True(t, s.Signup.Disabled)
	assert.False(t, s.Mail.Disabled)
	assert.Equal(t, c.Mail.Host, s.Host)
	assert.Equal(t, c.Mail.Port, s.Port)
	assert.EqualValues(t, c.Provider(), s.Provider.Internal)
	assert.True(t, s.Provider.External[provider.Google])
	assert.True(t, s.Provider.External[provider.GitLab])
	assert.True(t, s.Provider.External[provider.Heroku])
	assert.False(t, s.Provider.External[provider.BitBucket])
}
