package jwt

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseClaims(t *testing.T) {
	c := tconf.Config(t)
	c.JWT.Audience = "test"
	c.JWT.Expiration = 100
	// errors
	sign := jwt.New(jwt.GetSigningMethod(c.JWT.Algorithm))
	bad, err := sign.SignedString([]byte(c.JWT.Secret))
	require.NoError(t, err)
	s := &StandardClaims{}
	err = ParseClaims(c.JWT, bad, s)
	assert.Error(t, err)
	// expired
	c.JWT.Audience = ""
	claims := NewStandardClaims(c.JWT)
	claims.Subject = uuid.New().String()
	tok := NewToken(claims)
	good, err := tok.Bearer()
	require.NoError(t, err)
	d := tok.Expiration()
	assert.EqualValues(t, 0, d)
	s = &StandardClaims{}
	err = ParseClaims(c.JWT, good, s)
	assert.Error(t, err)
	c.JWT.Expiration = time.Minute
	// good
	claims = NewStandardClaims(c.JWT)
	claims.Subject = uuid.New().String()
	tok = NewToken(claims)
	good, err = tok.Bearer()
	require.NoError(t, err)
	s = &StandardClaims{}
	err = ParseClaims(c.JWT, good, s)
	assert.NoError(t, err)
	d = tok.Expiration()
	assert.EqualValues(t, c.JWT.Expiration, d)
	assert.Equal(t, claims.Subject, s.Subject)
	c.JWT.Expiration = 0
	x := NewStandardClaims(c.JWT)
	tok = NewToken(x)
	assert.EqualValues(t, 0, tok.Expiration())
}

func TestParseValue(t *testing.T) {
	c := tconf.Config(t)
	in := types.Map{
		key.Token: utils.SecureToken(),
		key.Email: tutils.RandomEmail(),
	}
	tok, err := NewSignedData(c.JWT, in)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)
	out, err := ParseData(c.JWT, tok)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out[key.Token], in[key.Token])
	assert.Equal(t, out[key.Email], in[key.Email])
	_, err = ParseData(c.JWT, "bad")
	assert.Error(t, err)
}
