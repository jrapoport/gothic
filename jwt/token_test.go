package jwt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	var sub = uuid.New().String()
	tests := []struct {
		config config.JWT
		Err    assert.ErrorAssertionFunc
	}{
		{
			config.JWT{
				Secret:    "a-secret",
				Algorithm: string(jwa.HS256),
				Scope:     "foo,bar",
			},
			assert.NoError,
		},
		{
			config.JWT{
				PEM: config.PEM{
					PrivateKey: "./testdata/rsa-key.pem",
				},
				Algorithm: string(jwa.RS512),
				Scope:     "foo,bar",
			},
			assert.NoError,
		},
		{
			config.JWT{
				PEM: config.PEM{
					PrivateKey: "./testdata/does-not-exist",
				},
				Algorithm: string(jwa.RS512),
			},
			assert.Error,
		},
		{
			config.JWT{
				PEM: config.PEM{
					PrivateKey: "./testdata/rsa-key.bad",
				},
				Algorithm: string(jwa.RS512),
			},
			assert.Error,
		},
		{
			config.JWT{
				PEM: config.PEM{
					PrivateKey: "./testdata/rsa-key.pem",
				},
				Algorithm: string(jwa.HS256),
			},
			assert.Error,
		},
	}
	for _, test := range tests {
		claims := NewStandardClaims(sub)
		tok := NewToken(test.config, claims)
		sig, err := tok.Bearer()
		test.Err(t, err)
		if err != nil {
			continue
		}
		var parsed StandardClaims
		err = ParseClaims(test.config, sig, &parsed)
		assert.NoError(t, err)
		assert.Equal(t, sub, parsed.Subject())
		assert.Equal(t, []string{"foo", "bar"}, parsed.Scope())
	}
}

func TestNewSignedData(t *testing.T) {
	var sub = uuid.New().String()
	tests := []struct {
		config config.JWT
	}{
		{
			config.JWT{
				Secret:    "a-secret",
				Algorithm: string(jwa.HS256),
			},
		},
		{
			config.JWT{
				PEM: config.PEM{
					PrivateKey: "./testdata/rsa-key.pem",
				},
				Algorithm: string(jwa.RS512),
			},
		},
	}
	for _, test := range tests {
		data := map[string]interface{}{
			"custom": sub,
		}
		sig, err := NewSignedData(test.config, data)
		assert.NoError(t, err)
		data, err = ParseData(test.config, sig)
		assert.NoError(t, err)
		assert.Equal(t, sub, data["custom"])
	}
}
