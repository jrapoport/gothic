package jwt

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
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
func NewUserClaims(c config.JWT, u *user.User) UserClaims {
	std := NewStandardClaims(c)
	std.Subject = u.ID.String()
	return UserClaims{
		StandardClaims: std,
		Provider:       u.Provider,
		Admin:          u.IsAdmin(),
		Restricted:     u.IsRestricted(),
		Confirmed:      u.IsConfirmed(),
		Verified:       u.IsVerified(),
	}
}

// NewUserToken returns a new Token for the user with UserClaims.
func NewUserToken(c config.JWT, u *user.User) *Token {
	claims := NewUserClaims(c, u)
	return NewToken(claims)
}

// UserID returns the jwt subject as a uuid.
func (s UserClaims) UserID() uuid.UUID {
	uid, err := uuid.Parse(s.Subject)
	if err != nil {
		return uuid.Nil
	}
	return uid
}

// ParseUserClaims parses and returns a set of UserClaims form a token.
func ParseUserClaims(c config.JWT, token string) (UserClaims, error) {
	claims := UserClaims{}
	err := ParseClaims(c, token, &claims)
	if err != nil {
		return UserClaims{}, err
	}
	return claims, nil
}
