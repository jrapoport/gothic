package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	jwtSecret = "i-am-a-secret"
	jwtKey    = "./testdata/test-key.pem"
	jwtAlgo   = "HS384"
	jwtIss    = "foo"
	jwtAud    = "bar"
	jwtExp    = 100 * time.Minute
)

func TestJWT(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		j := c.JWT
		assert.Equal(t, jwtSecret+test.mark, j.Secret)
		assert.Equal(t, jwtKey+test.mark, j.PrivateKey)
		assert.Equal(t, jwtKey+test.mark, j.PublicKey)
		assert.Equal(t, jwtAlgo+test.mark, j.Algorithm)
		assert.Equal(t, jwtIss+test.mark, j.Issuer)
		assert.Equal(t, jwtAud+test.mark, j.Audience)
		assert.Equal(t, jwtExp, j.Expiration)
	})
}

// tests the ENV vars are correctly taking precedence
func TestJWT_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			j := c.JWT
			assert.Equal(t, jwtSecret, j.Secret)
			assert.Equal(t, jwtKey, j.PrivateKey)
			assert.Equal(t, jwtKey, j.PublicKey)
			assert.Equal(t, jwtAlgo, j.Algorithm)
			assert.Equal(t, jwtIss, j.Issuer)
			assert.Equal(t, jwtAud, j.Audience)
			assert.Equal(t, jwtExp, j.Expiration)
		})
	}
}

// test the *un-normalized* defaults with load
func TestJWT_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := jwtDefaults
	j := c.JWT
	assert.Equal(t, def, j)
}

func TestJWT_Normalization(t *testing.T) {
	j := JWT{}
	j.Secret = jwtSecret
	j.PrivateKey = jwtKey
	j.PublicKey = jwtKey
	j.normalize(Service{
		Name:    service,
		SiteURL: siteURL,
	}, jwtDefaults)
	assert.Equal(t, strings.ToLower(service), j.Issuer)
}

func TestJWT_CheckRequired(t *testing.T) {
	reqTests := []struct {
		secret  string
		private string
		public  string
		Err     assert.ErrorAssertionFunc
	}{
		{"", "", "", assert.Error},
		{"secret", "", "", assert.NoError},
		{"", "bad", "", assert.Error},
		{"", jwtKey, "", assert.NoError},
		{"", "", jwtKey, assert.Error},
		{"", jwtKey, "bad", assert.Error},
		{"", jwtKey, jwtKey, assert.NoError},
	}
	for _, test := range reqTests {
		j := JWT{}
		j.Secret = test.secret
		j.PrivateKey = test.private
		j.PublicKey = test.public
		err := j.CheckRequired()
		test.Err(t, err)
	}
}
