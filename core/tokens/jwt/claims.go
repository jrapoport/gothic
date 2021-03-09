package jwt

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/jrapoport/gothic/config"
)

// Claims interface for jwt.
type Claims interface {
	jwt.Claims
	Secret() []byte
	Method() jwt.SigningMethod
	Standard() StandardClaims
}

// StandardClaims holds standard jwt claims.
type StandardClaims struct {
	jwt.StandardClaims
	method jwt.SigningMethod
	secret []byte
}

// NewStandardClaims returns a new set of standard jwt claims.
func NewStandardClaims(c config.JWT) StandardClaims {
	std := StandardClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   c.Issuer,
			IssuedAt: jwt.At(time.Now().UTC()),
		},
		method: jwt.GetSigningMethod(c.Algorithm),
		secret: []byte(c.Secret),
	}
	if c.Audience != "" {
		aud := strings.Split(c.Audience, ",")
		std.Audience = aud
	}
	if c.Expiration > 0 {
		exp := std.IssuedAt.Add(c.Expiration)
		std.ExpiresAt = jwt.At(exp)
	}
	return std
}

// Standard returns a the standard jwt claims.
func (s StandardClaims) Standard() StandardClaims {
	return s
}

// Secret returns the jwt secret to use.
func (s StandardClaims) Secret() []byte {
	return s.secret
}

// Method returns the jwt signature method to use.
func (s StandardClaims) Method() jwt.SigningMethod {
	return s.method
}
