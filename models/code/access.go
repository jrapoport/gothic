package code

import (
	"errors"
	"time"

	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/utils"
	"gorm.io/gorm"
)

// Format is the format of the signup code.
type Format uint8

// Access Code formats
const (
	Invite Format = iota
	PIN
)

// NoExpiration indicates the code will not expire.
const NoExpiration = token.NoExpiration

const (
	// InfiniteUse indicates this code may be used an infinite
	// number of times. This implies type Infinite.
	InfiniteUse = token.InfiniteUse

	// SingleUse indicates this code may be used once.
	// This implies type Single.
	SingleUse = token.SingleUse
)

// AccessCode can be used for codes.
type AccessCode struct {
	token.AccessToken
	Format Format `json:"format"`
}

// NewAccessCode generates a new code of the format. The type is inferred
// from the number of uses since it is not clear they are different.
func NewAccessCode(f Format, uses int, exp time.Duration) *AccessCode {
	var c string
	switch f {
	case PIN:
		c = utils.PINCode()
	case Invite:
		fallthrough
	default:
		c = utils.SecureToken()
	}
	t := *token.NewAccessToken(c, uses, exp)
	return &AccessCode{t, f}
}

// BeforeCreate runs before create.
func (ac AccessCode) BeforeCreate(db *gorm.DB) error {
	if ac.Format == PIN && utils.IsDebugPIN(ac.Code()) {
		return errors.New("invalid code")
	}
	return ac.AccessToken.BeforeCreate(db)
}

// Usable returns true if the code is usable.
func (ac AccessCode) Usable() bool {
	if ac.Format == PIN && utils.IsDebugPIN(ac.Code()) {
		// a debug pin only works with debug code
		return true
	}
	if ac.CreatedAt.IsZero() {
		return false
	}
	return ac.AccessToken.Usable()
}

// Code returns the access code as a string.
func (ac AccessCode) Code() string {
	return ac.String()
}
