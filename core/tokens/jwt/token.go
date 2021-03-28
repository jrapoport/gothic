package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/segmentio/encoding/json"
)

// Token is a struct to hold extended jwt token.
// TODO: support public/private keys etc.
type Token struct {
	jwt.Token
	method jwa.SignatureAlgorithm
	secret []byte
}

func newToken(c config.JWT, v interface{}) *Token {
	if v == nil {
		return nil
	}
	tok := jwt.New()
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, tok)
	if err != nil {
		return nil
	}
	iss := c.Issuer
	iat := time.Now().UTC().Truncate(time.Microsecond)
	_ = tok.Set(jwt.IssuerKey, iss)
	_ = tok.Set(jwt.IssuedAtKey, iat)
	algo := jwa.SignatureAlgorithm(c.Algorithm)
	sec := []byte(c.Secret)
	return &Token{tok, algo, sec}
}

// NewToken returns a new jwt token for the claims.
// ignoring errors here is ok because we can tightly
// control all the incoming types.
func NewToken(c config.JWT, claims Claims) *Token {
	tok := newToken(c, claims)
	if tok == nil {
		return nil
	}
	iat := tok.IssuedAt()
	sub := claims.Subject()
	_ = tok.Set(jwt.SubjectKey, sub)
	if c.Audience != "" {
		aud := strings.Split(c.Audience, ",")
		_ = tok.Set(jwt.AudienceKey, aud)
	}
	if c.Expiration > 0 {
		exp := iat.Add(c.Expiration).Truncate(time.Microsecond)
		_ = tok.Set(jwt.ExpirationKey, exp)
	}
	return tok
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

// Expires returns the expiration for the token.
func (t Token) Expiration() time.Duration {
	iat := t.Token.IssuedAt()
	exp := t.ExpiresAt()
	if exp.IsZero() {
		return 0
	}
	return exp.Sub(iat)
}

// NewSignedData returns a signed jwt token for the Map.
func NewSignedData(c config.JWT, d types.Map) (string, error) {
	tok := newToken(c, d)
	if tok == nil {
		return "", errors.New("invalid token")
	}
	return tok.Bearer()
}
