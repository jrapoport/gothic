package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	service  = "example"
	siteURL  = "http://example.com"
	siteLogo = "http://example.com/logo.png"
)

func TestService(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		assert.Equal(t, service+test.mark, c.Name)
		assert.Equal(t, siteURL+test.mark, c.SiteURL)
		assert.Equal(t, siteLogo+test.mark, c.SiteLogo)
	})
}

// tests the ENV vars are correctly taking precedence
func TestService_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			assert.Equal(t, service, c.Name)
			assert.Equal(t, siteURL, c.SiteURL)
			assert.Equal(t, siteLogo, c.SiteLogo)
		})
	}
}

// test the *un-normalized* defaults with load
func TestService_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	s := c.Service
	def := serviceDefaults
	assert.Equal(t, def, s)
}

func TestService_Normalization(t *testing.T) {
	s := Service{
		SiteURL: siteURL,
	}
	err := s.normalize()
	assert.NoError(t, err)
	assert.Equal(t, BuildVersion(), s.Version())
	s.SiteURL = "\n"
	err = s.normalize()
	assert.Error(t, err)
}
