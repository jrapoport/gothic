package jwt

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
)

func TestNewStandardClaims(t *testing.T) {
	c := tconf.Config(t)
	c.JWT.Audience = "test"
	c.JWT.Expiration = 100
	claims := NewStandardClaims(c.JWT)
	assert.Equal(t, c.JWT.Secret, string(claims.Secret()))
	assert.Equal(t, c.JWT.Algorithm, claims.Method().Alg())
	s := claims.Standard()
	assert.Equal(t, c.Service.Name, s.Issuer)
}
