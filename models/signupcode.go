package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/utils"
	"gorm.io/gorm"
)

func init() {
	storage.AddMigration(&SignupCode{})
}

// Format is the format of the signup code.
type Format uint8

const (
	// PINFormat is a PIN style signup code.
	PINFormat Format = iota
)

// Type is the type of code (e.g. single use)
type Type uint8

const (
	// SingleUse is a code that can be used once.
	SingleUse Type = iota

	// MultiUse is a code that can be used repeatedly.
	MultiUse
)

const (
	// InfiniteUses indicates this code may be used an infinite number of times.
	// infinite MaxUses is only valid for MultiUse codes
	InfiniteUses = -1
)

// maxPIN is the max length of a PIN (6).
const maxPINCode = 6

// debugPINCode is bypass code that is only valid in debug builds.
const debugPINCode = "000000"

type UserIDs struct {
	gorm.Model
	UserID string
}

// SignupCode can be used for signup signup codes.
type SignupCode struct {
	gorm.Model
	Code       string    `json:"code" gorm:"index"`
	Format     Format    `json:"format"`
	Type       Type      `json:"type"`
	ReferralID uuid.UUID `json:"referral_id" gorm:"type:char(36)"`
	MaxUses    int       `json:"max_uses"`
	// NOTE: Uses is kept in sync with len(Used) purely by convention
	// in UseSignupCode since there is no concept of "un"-using a code.
	Used   int        `json:"used"`
	UsedAt *time.Time `json:"used_at"`
	Users  []User     `json:"users"`
}

func (ac SignupCode) Usable() bool {
	if ac.Type == MultiUse && ac.MaxUses == InfiniteUses {
		return true
	} else if ac.Used >= ac.MaxUses {
		return false
	}
	return !ac.DeletedAt.Valid
}

// NewSignupCode generates a new code of the format and type.
func NewSignupCode(fmt Format, typ Type) *SignupCode {
	var code string
	switch fmt {
	case PINFormat:
		for {
			code = utils.RandomPIN(maxPINCode)
			// make sure this isn't the debug code
			if code != debugPINCode {
				break
			}
		}
	default:
		return nil
	}
	maxUses := 0
	switch typ {
	case SingleUse:
		maxUses = 1
	case MultiUse:
		maxUses = InfiniteUses
		break
	}
	return &SignupCode{
		Code:       code,
		Format:     fmt,
		Type:       typ,
		ReferralID: SystemUserUUID,
		MaxUses:    maxUses,
	}
}

// NewUniqueSignupCode returns a new unique signup code.
func NewUniqueSignupCode(tx *storage.Connection, f Format, t Type) (*SignupCode, error) {
	code := NewSignupCode(f, t)
	has, err := HasSignupCode(tx, code.Code)
	if err != nil {
		return nil, err
	} else if has {
		return NewUniqueSignupCode(tx, f, t)
	}
	return code, nil
}

// CreateSignupCode creates a unique signup code.
func CreateSignupCode(tx *storage.Connection, f Format, t Type, unique bool) (code *SignupCode, err error) {
	err = tx.Transaction(func(tx *storage.Connection) error {
		if unique {
			code, err = NewUniqueSignupCode(tx, f, t)
			if err != nil {
				return err
			}
		} else {
			code = NewSignupCode(f, t)
		}
		if err = tx.Create(code).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

// HasSignupCode returns whether the code exists.
func HasSignupCode(tx *storage.Connection, code string) (bool, error) {
	if code == "" {
		return false, errors.New("invalid signup code")
	}
	return storage.Has(tx.Where("code = ?", code), new(SignupCode))
}

// GetSignupCode finds a user with the matching code.
func GetSignupCode(tx *storage.Connection, code string) (*SignupCode, error) {
	if code == "" {
		return nil, errors.New("invalid signup code")
	}
	c := new(SignupCode)
	err := tx.Where("code = ?", code).First(c).Error //.Preload("Used").First(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

// CanUseSignupCode
func CanUseSignupCode(tx *storage.Connection, code string) (bool, error) {
	if conf.Debug && code == debugPINCode {
		return true, nil
	}
	ac, err := GetSignupCode(tx, code)
	if err != nil {
		return false, err
	}
	return ac.Usable(), nil
}

func UseSignupCode(tx *storage.Connection, code string, user *User) error {
	if conf.Debug && code == debugPINCode {
		return nil
	}
	if code == "" {
		return errors.New("invalid signup code")
	}
	ac, err := GetSignupCode(tx, code)
	if err != nil {
		return err
	}
	if !ac.Usable() {
		return errors.New("signup code cannot be used")
	}
	now := time.Now()
	ac.Users = append(ac.Users, *user)
	ac.Used++
	ac.UsedAt = &now
	return tx.Save(ac).Error
}
