package jwt

import "github.com/lestrrat-go/jwx/jwt"

// Claims interface for jwt.
type Claims interface {
	jwt.Token
	parseToken(*Token)
}

// StandardClaims holds standard jwt claims.
type StandardClaims struct {
	jwt.Token `json:"-"`
}

var _ Claims = (*StandardClaims)(nil)

// NewStandardClaims returns a new set of standard jwt claims.
func NewStandardClaims(sub string) *StandardClaims {
	std := &StandardClaims{jwt.New()}
	std.SetSubject(sub)
	return std
}

// SetSubject sets the subject for the token
func (c *StandardClaims) SetSubject(sub string) {
	// we can safely ignore an error here
	// because the it is strongly typed
	_ = c.Set(jwt.SubjectKey, sub)
}

// ParseToken handles the parsed values coming back from a token
func (c *StandardClaims) parseToken(tok *Token) {
	c.Token = tok.Token
}
