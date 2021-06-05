package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	jwtSecret = "i-am-a-secret"
	jwtAlgo   = "HS384"
	jwtIss    = "foo"
	jwtAud    = "bar"
	jwtExp    = 100 * time.Minute
	jwtPrvKey = "./testdata/test-key.pem"
	jwtPubKey = "./testdata/test-key.pem.pub"
	jwtBadKey = "./testdata/test-key.bad"
)

func TestJWT(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		j := c.JWT
		assert.Equal(t, jwtSecret+test.mark, j.Secret)
		assert.Equal(t, jwtPrvKey+test.mark, j.PEM.PrivateKey)
		assert.Equal(t, jwtPrvKey+test.mark, j.PEM.PublicKey)
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
			assert.Equal(t, jwtPrvKey, j.PEM.PrivateKey)
			assert.Equal(t, jwtPrvKey, j.PEM.PublicKey)
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
	j.PEM.PrivateKey = jwtPrvKey
	j.PEM.PublicKey = jwtPrvKey
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
		{"", jwtPrvKey, "", assert.NoError},
		{"", "", jwtPrvKey, assert.Error},
		{"", jwtPrvKey, "bad", assert.Error},
		{"", jwtPrvKey, jwtPrvKey, assert.NoError},
	}
	for _, test := range reqTests {
		j := JWT{}
		j.Secret = test.secret
		j.PEM.PrivateKey = test.private
		j.PEM.PublicKey = test.public
		err := j.CheckRequired()
		test.Err(t, err)
	}
}

func TestJWT_Keys(t *testing.T) {
	const noKey = "./does-not-exist"
	keyTests := []struct {
		config JWT
		Nil    assert.ValueAssertionFunc
	}{
		{
			JWT{Secret: jwtSecret},
			assert.NotNil,
		},
		{
			JWT{PEM: PEM{PrivateKey: jwtPrvKey}},
			assert.NotNil,
		},
		{
			JWT{PEM: PEM{
				PrivateKey: jwtPrvKey,
				PublicKey:  jwtPubKey,
			}},
			assert.NotNil,
		},
		{
			JWT{PEM: PEM{PrivateKey: noKey}},
			assert.Nil,
		},
		{
			JWT{PEM: PEM{PrivateKey: jwtBadKey}},
			assert.Nil,
		},
	}
	for _, test := range keyTests {
		assert.Nil(t, test.config.sk)
		assert.Nil(t, test.config.pk)
		test.Nil(t, test.config.PrivateKey())
		test.Nil(t, test.config.sk)
		test.Nil(t, test.config.PublicKey())
		test.Nil(t, test.config.pk)
	}
}
