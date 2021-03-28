package jwt

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseClaims(t *testing.T) {
	c := tconf.Config(t).JWT
	// errors
	c.Audience = "test"
	c.Expiration = 100
	sign := jwt.New()
	alg := jwa.SignatureAlgorithm(c.Algorithm)
	bad, err := jwt.Sign(sign, alg, []byte(c.Secret))
	require.NoError(t, err)
	claims := &StandardClaims{}
	err = ParseClaims(c, string(bad), claims)
	assert.Error(t, err)
	// expired
	c.Audience = ""
	claims = NewStandardClaims(uuid.New().String())
	tok := NewToken(c, claims)
	expired, err := tok.Bearer()
	require.NoError(t, err)
	exp := tok.Expiration()
	assert.EqualValues(t, 0, exp)
	claims = &StandardClaims{}
	err = ParseClaims(c, expired, claims)
	assert.Error(t, err)
	c.Expiration = time.Minute
	// good
	claims = NewStandardClaims(uuid.New().String())
	tok = NewToken(c, claims)
	good, err := tok.Bearer()
	require.NoError(t, err)
	s := &StandardClaims{}
	err = ParseClaims(c, good, s)
	assert.NoError(t, err)
	d := tok.Expiration()
	assert.EqualValues(t, c.Expiration, d)
	assert.Equal(t, claims.Subject(), s.Subject())
	c.Expiration = 0
	claims = NewStandardClaims(uuid.New().String())
	tok = NewToken(c, claims)
	assert.EqualValues(t, 0, tok.Expiration())
	// bad subject
	claims = NewStandardClaims("")
	tok = NewToken(c, claims)
	nosub, err := tok.Bearer()
	require.NoError(t, err)
	err = ParseClaims(c, nosub, s)
	assert.Error(t, err)
	// bad claims
	tok = NewToken(c, Claims(nil))
	assert.Nil(t, tok)
}

func TestParseValue(t *testing.T) {
	c := tconf.Config(t)
	email := tutils.RandomEmail()
	in := types.Map{
		key.Token: utils.SecureToken(),
		key.Email: email,
	}
	tok, err := NewSignedData(c.JWT, in)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)
	out, err := ParseData(c.JWT, tok)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out[key.Token], in[key.Token])
	assert.Equal(t, out[key.Email], in[key.Email])
	// bad signature
	decoded, err := base64.RawStdEncoding.DecodeString(tok)
	dec := strings.Replace(string(decoded), email, tutils.RandomEmail(), 1)
	enc := base64.RawStdEncoding.EncodeToString([]byte(dec))
	_, err = ParseData(c.JWT, enc)
	assert.Error(t, err)
	// bad token
	_, err = ParseData(c.JWT, "bad")
	assert.Error(t, err)
	// bad algo
	c.JWT.Algorithm = ""
	_, err = ParseData(c.JWT, "")
	assert.Error(t, err)
	_, err = NewSignedData(c.JWT, in)
	assert.Error(t, err)
}
