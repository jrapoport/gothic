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
}

var _ Claims = (*UserClaims)(nil)

// NewUserClaims returns a new set of claims for the user.
func NewUserClaims(u *user.User) *UserClaims {
	if u == nil {
		return nil
	}
	c := &UserClaims{
		*NewStandardClaims(u.ID.String()),
	}
	_ = c.Set(ProviderKey, u.Provider.String())
	_ = c.Set(AdminKey, u.IsAdmin())
	_ = c.Set(RestrictedKey, u.IsRestricted())
	_ = c.Set(ConfirmedKey, u.IsConfirmed())
	_ = c.Set(VerifiedKey, u.IsVerified())
	return c
}

// UserID returns the jwt subject as a uuid.
func (c UserClaims) UserID() uuid.UUID {
	uid, err := uuid.Parse(c.Subject())
	if err != nil {
		return uuid.Nil
	}
	return uid
}

// Provider returns the provider name
func (c UserClaims) Provider() provider.Name {
	name, ok := c.Get(ProviderKey)
	if !ok {
		return provider.Unknown
	}
	return provider.Name(name.(string))
}

// Admin returns true if admin
func (c UserClaims) Admin() bool {
	v, _ := c.Get(AdminKey)
	return v.(bool)
}

// Restricted returns true if restricted
func (c UserClaims) Restricted() bool {
	v, _ := c.Get(RestrictedKey)
	return v.(bool)
}

// Confirmed returns true if confirmed
func (c UserClaims) Confirmed() bool {
	v, _ := c.Get(ConfirmedKey)
	return v.(bool)
}

// Verified returns true if verified
func (c UserClaims) Verified() bool {
	v, _ := c.Get(VerifiedKey)
	return v.(bool)
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
