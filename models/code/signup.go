package code

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

func init() {
	store.AddAutoMigrationWithIndexes("2000-signup_codes",
		SignupCode{}, token.AccessTokenIndexes)
}

// Signup class
const Signup token.Class = "signup"

// SignupCode can be used for signup signup codes.
type SignupCode struct {
	AccessCode
	SentAt *time.Time  `json:"sent_at"`
	Users  []user.User `json:"users" gorm:"foreignKey:SignupCode"`
}

var _ token.Token = (*SignupCode)(nil)

// NewSignupCode generates a new code of the format and type.
func NewSignupCode(userID uuid.UUID, f Format, uses int) *SignupCode {
	ac := NewAccessCode(f, uses, token.NoExpiration)
	ac.UserID = userID
	return &SignupCode{AccessCode: *ac}
}

// Class returns the class of the signup code.
func (sc SignupCode) Class() token.Class {
	return Signup
}

// HasCode returns true if the code was found.
func (sc SignupCode) HasCode(tx *store.Connection) (bool, error) {
	if sc.Token == "" {
		return false, errors.New("invalid code")
	}
	return tx.Has(&sc, "token = ?", sc.Token)
}

// UseCode uses a signup code.
func (sc *SignupCode) UseCode(tx *store.Connection, u *user.User) error {
	if !validUser(u) {
		return errors.New("invalid user")
	}
	if sc.Format == PIN && utils.IsDebugPIN(sc.Code()) {
		return nil
	}
	sc.Users = append(sc.Users, *u)
	return tokens.UseToken(tx, sc)
}

func validUser(u *user.User) bool {
	return u != nil && u.IsRestricted() && u.SignupCode == nil
}
