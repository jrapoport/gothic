package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	maskEmails   = false
	rateLimit    = 100 * time.Minute
	rootPassword = "password"
	recapKey     = "RECAPTCHA-KEY"
	recapLogin   = false
	userRx       = "[A-Za-z]{3}[0-9][A-Z]{2}[!@#$%^&*]"
	passRx       = "FOO[A-Z]{10}[0-9]{2}"
	duration     = 100 * time.Minute
)

func TestSecurity(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		s := c.Security
		assert.Equal(t, rootPassword+test.mark, s.RootPassword)
		assert.Equal(t, maskEmails, s.MaskEmails)
		assert.Equal(t, rateLimit, s.RateLimit)
		assert.Equal(t, recapKey+test.mark, s.Recaptcha.Key)
		assert.Equal(t, recapLogin, s.Recaptcha.Login)
		assert.Equal(t, userRx+test.mark, s.Validation.UsernameRegex)
		assert.Equal(t, passRx+test.mark, s.Validation.PasswordRegex)
		assert.Equal(t, duration, s.Cookies.Duration)
	})
}

// tests the ENV vars are correctly taking precedence
func TestSecurity_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			s := c.Security
			assert.Equal(t, rootPassword, s.RootPassword)
			assert.Equal(t, maskEmails, s.MaskEmails)
			assert.Equal(t, rateLimit, s.RateLimit)
			assert.Equal(t, recapKey, s.Recaptcha.Key)
			assert.Equal(t, recapLogin, s.Recaptcha.Login)
			assert.Equal(t, userRx, s.Validation.UsernameRegex)
			assert.Equal(t, passRx, s.Validation.PasswordRegex)
			assert.Equal(t, duration, s.Cookies.Duration)
		})
	}
}

// test the *un-normalized* defaults with load
func TestSecurity_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := securityDefaults
	s := c.Security
	assert.Equal(t, def, s)
}

func TestSecurity_Normalization(t *testing.T) {
	s := Security{}
	s.RootPassword = rootPassword
	err := s.normalize(Service{
		Name:    service,
		SiteURL: siteURL,
	})
	assert.NoError(t, err)
	assert.Equal(t, cookieDuration, s.Cookies.Duration)
	s.Validation.PasswordRegex = "a(?=r)"
	err = s.normalize(serviceDefaults)
	assert.Error(t, err)
	s.Validation.UsernameRegex = "a(?=r)"
	err = s.normalize(serviceDefaults)
	assert.Error(t, err)
}
