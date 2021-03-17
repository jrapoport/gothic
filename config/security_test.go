package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	maskEmails = false
	rateLimit  = 100 * time.Minute
	requestID  = "foobar"
	jwtSecret  = "i-am-a-secret"
	jwtAlgo    = "HS384"
	jwtIss     = "foo"
	jwtAud     = "bar"
	jwtExp     = 100 * time.Minute
	recapKey   = "RECAPTCHA-KEY"
	recapLogin = false
	userRx     = "[A-Za-z]{3}[0-9][A-Z]{2}[!@#$%^&*]"
	passRx     = "FOO[A-Z]{10}[0-9]{2}"
	duration   = 100 * time.Minute
)

func TestSecurity(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		s := c.Security
		assert.Equal(t, maskEmails, s.MaskEmails)
		assert.Equal(t, rateLimit, s.RateLimit)
		assert.Equal(t, requestID+test.mark, s.RequestID)
		assert.Equal(t, jwtSecret+test.mark, s.JWT.Secret)
		assert.Equal(t, jwtAlgo+test.mark, s.JWT.Algorithm)
		assert.Equal(t, jwtIss+test.mark, s.JWT.Issuer)
		assert.Equal(t, jwtAud+test.mark, s.JWT.Audience)
		assert.Equal(t, jwtExp, s.JWT.Expiration)
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
			assert.Equal(t, maskEmails, s.MaskEmails)
			assert.Equal(t, rateLimit, s.RateLimit)
			assert.Equal(t, requestID, s.RequestID)
			assert.Equal(t, jwtSecret, s.JWT.Secret)
			assert.Equal(t, jwtAlgo, s.JWT.Algorithm)
			assert.Equal(t, jwtIss, s.JWT.Issuer)
			assert.Equal(t, jwtAud, s.JWT.Audience)
			assert.Equal(t, jwtExp, s.JWT.Expiration)
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
	s.JWT.Secret = jwtSecret
	err := s.normalize(Service{
		Name:    service,
		SiteURL: siteURL,
	})
	assert.NoError(t, err)
	assert.Equal(t, strings.ToLower(service), s.JWT.Issuer)
	assert.Equal(t, cookieDuration, s.Cookies.Duration)
	s.Validation.PasswordRegex = "a(?=r)"
	err = s.normalize(serviceDefaults)
	assert.Error(t, err)
	s.Validation.UsernameRegex = "a(?=r)"
	err = s.normalize(serviceDefaults)
	assert.Error(t, err)
}
