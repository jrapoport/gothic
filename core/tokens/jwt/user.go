package jwt

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
)

// UserClaims jwt keys
const (
	ProviderKey   = "pvd"
	AdminKey      = "adm"
	RestrictedKey = "rst"
	ConfirmedKey  = "cnf"
	VerifiedKey   = "vrd"
)

// UserClaims is a struct to hold extended jwt claims
type UserClaims struct {
	StandardClaims
	Provider   provider.Name `json:"pvd,omitempty"`
	Admin      bool          `json:"adm,omitempty"`
	Restricted bool          `json:"rst,omitempty"`
	Confirmed  bool          `json:"cnf,omitempty"`
	Verified   bool          `json:"vrd,omitempty"`
}

var _ Claims = (*UserClaims)(nil)

// NewUserClaims returns a new set of claims for the user.
func NewUserClaims(u *user.User) *UserClaims {
	c := &UserClaims{
		StandardClaims: *NewStandardClaims(""),
	}
	if u == nil {
		return c
	}
	c.SetSubject(u.ID.String())
	c.Provider = u.Provider
	c.Admin = u.IsAdmin()
	c.Restricted = u.IsRestricted()
	c.Confirmed = u.IsConfirmed()
	c.Verified = u.IsVerified()
	return c
}

// ParseToken handles the parsed values coming back from a token
// TODO: consider using the token directly here instead.
func (c *UserClaims) ParseToken(tok *Token) {
	c.StandardClaims.ParseToken(tok)
	if v, ok := c.Get(ProviderKey); ok {
		c.Provider = provider.Name(v.(string))
	}
	if v, ok := c.Get(AdminKey); ok {
		c.Admin = v.(bool)
	}
	if v, ok := c.Get(RestrictedKey); ok {
		c.Restricted = v.(bool)
	}
	if v, ok := c.Get(ConfirmedKey); ok {
		c.Confirmed = v.(bool)
	}
	if v, ok := c.Get(VerifiedKey); ok {
		c.Verified = v.(bool)
	}
}

// UserID returns the jwt subject as a uuid.
func (c UserClaims) UserID() uuid.UUID {
	uid, err := uuid.Parse(c.Subject())
	if err != nil {
		return uuid.Nil
	}
	return uid
}

// ParseUserClaims parses and returns a set of UserClaims form a token.
func ParseUserClaims(c config.JWT, token string) (*UserClaims, error) {
	claims := &UserClaims{}
	err := ParseClaims(c, token, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// NewUserToken returns a new Token for the user with UserClaims.
func NewUserToken(c config.JWT, u *user.User) *Token {
	return NewToken(c, NewUserClaims(u))
}
