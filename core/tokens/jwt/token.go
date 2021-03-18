package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types"
)

// Token is a struct to hold extended jwt token.
type Token struct {
	*jwt.Token
}

// NewToken returns a new jwt token for the claims.
func NewToken(claims Claims) *Token {
	t := jwt.NewWithClaims(claims.Method(), claims)
	return &Token{t}
}

// Claims returns the claims for the token.
func (t Token) Claims() Claims {
	return t.Token.Claims.(Claims)
}

// Bearer signs the claims and returns the result as a string.
func (t Token) Bearer() (string, error) {
	return t.SignedString(t.Claims().Secret())
}

// Expiration returns the expiration for the token.
func (t Token) Expiration() time.Duration {
	var issued time.Time
	std := t.Claims().Standard()
	if std.ExpiresAt == nil {
		return 0
	}
	if std.IssuedAt != nil {
		issued = std.IssuedAt.Time
	}
	return std.ExpiresAt.Sub(issued)
}

// NewSignedData returns a signed jwt token for the Map.
func NewSignedData(c config.JWT, d types.Map) (string, error) {
	t := jwt.New(jwt.GetSigningMethod(c.Algorithm))
	for k, v := range d {
		t.Header[k] = v
	}
	return t.SignedString([]byte(c.Secret))
}
