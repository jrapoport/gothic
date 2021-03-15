package account

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"gorm.io/gorm"
)

func init() {
	var linkedIndexes = []string{
		"idx_type_provider_account_id",
	}
	store.AddAutoMigrationWithIndexes("4500-linked",
		LinkedAccount{}, linkedIndexes)
}

// Type is the type of linked account
type Type int

const (
	// Auth account type.
	Auth Type = iota

	// Payment account type.
	Payment
)

func (t Type) String() string {
	switch t {
	case Auth:
		return "auth"
	case Payment:
		return "payment"
	default:
		return ""
	}
}

// LinkedAccount holds a linked account
type LinkedAccount struct {
	gorm.Model
	UserID    uuid.UUID     `json:"user_id" gorm:"type:char(36)"`
	Type      Type          `json:"type" gorm:"uniqueIndex:idx_type_provider_account_id"`
	Provider  provider.Name `json:"provider" gorm:"uniqueIndex:idx_type_provider_account_id;type:varchar(255)"`
	AccountID string        `json:"account_id" gorm:"uniqueIndex:idx_type_provider_account_id;type:varchar(320)"`
	Email     string        `json:"email" gorm:"type:varchar(320)"`
	Data      types.Map     `json:"data"`
}

// Valid returns nil if the linked account is valid.
func (l LinkedAccount) Valid() error {
	if l.Type.String() == "" {
		return errors.New("invalid type")
	}
	if l.Provider == provider.Unknown {
		return errors.New("invalid provider")
	}
	if l.AccountID == "" {
		return errors.New("invalid account id")
	}
	return nil
}
