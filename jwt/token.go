package jwt

import (
	"strings"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// Token is a struct to hold extended jwt token.
type Token struct {
	jwt.Token
	method jwa.SignatureAlgorithm
	secret jwk.Key
}

// NewToken returns a new jwt token for the claims.
// ignoring errors here is ok because we can tightly
// control all the incoming types.
func NewToken(c config.JWT, claims Claims) *Token {
	tok := jwt.New()
	if claims != nil {
		tok, _ = claims.Clone()
	}
	iss := c.Issuer
	iat := time.Now().UTC().Truncate(time.Microsecond)
	_ = tok.Set(jwt.IssuerKey, iss)
	_ = tok.Set(jwt.IssuedAtKey, iat)
	if c.Audience != "" {
		aud := strings.Split(c.Audience, ",")
		_ = tok.Set(jwt.AudienceKey, aud)
	}
	if c.Expiration > 0 {
		exp := iat.Add(c.Expiration).Truncate(time.Microsecond)
		_ = tok.Set(jwt.ExpirationKey, exp)
	}
	if c.Scope != "" {
		scope := strings.Split(c.Scope, ",")
		_ = tok.Set(ScopeKey, strings.Join(scope, " "))
	}
	algo := jwa.SignatureAlgorithm(c.Algorithm)
	key := c.PrivateKey()
	return &Token{tok, algo, key}
}

// Bearer signs the claims and returns the result as a string.
func (t Token) Bearer() (string, error) {
	b, err := jwt.Sign(t.Token, t.method, t.secret)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ExpiresAt returns the expiration time for the token.
func (t Token) ExpiresAt() time.Time {
	return t.Token.Expiration()
}

// Expiration returns the expiration for the token.
func (t Token) Expiration() time.Duration {
	iat := t.Token.IssuedAt()
	exp := t.ExpiresAt()
	if exp.IsZero() {
		return 0
	}
	return exp.Sub(iat)
}

// NewSignedData returns a signed jwt token for the Map.
func NewSignedData(c config.JWT, data types.Map) (string, error) {
	tok := NewToken(c, nil)
	for k, v := range data {
		_ = tok.Set(k, v)
	}
	return tok.Bearer()
}
