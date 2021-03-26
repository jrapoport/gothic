package account

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"gorm.io/gorm"
)

func init() {
	var accountIndexes = []string{
		"idx_provider_account_id",
	}
	store.AddAutoMigrationWithIndexes("4500-linked-account",
		Account{}, accountIndexes)
}

// Account holds a linked account
type Account struct {
	gorm.Model
	Type      Type          `json:"type" gorm:"<-:create"`
	Provider  provider.Name `json:"provider" gorm:"<-:create;uniqueIndex:idx_provider_account_id;type:varchar(255)"`
	AccountID string        `json:"account_id" gorm:"<-:create;uniqueIndex:idx_provider_account_id;type:varchar(320)"`
	Email     string        `json:"email" gorm:"type:varchar(320)"`
	Data      types.Map     `json:"data"`
	UserID    uuid.UUID     `json:"user_id" gorm:"<-:create;type:char(36)"`
}

func NewAccount(p provider.Name, accountID, email string, data types.Map) *Account {
	if data == nil {
		data = types.Map{}
	}
	return &Account{
		Type:      providerType(p),
		Provider:  p,
		AccountID: accountID,
		Email:     email,
		Data:      data,
	}
}

// BeforeSave runs before create or update.
func (la *Account) BeforeCreate(*gorm.DB) error {
	if la.Data == nil {
		la.Data = types.Map{}
	}
	return la.Valid()
}

// Valid returns nil if the linked account is valid.
func (la *Account) Valid() error {
	if la.Type == None {
		return errors.New("invalid type")
	}
	if la.Provider == provider.Unknown {
		return errors.New("invalid provider")
	}
	if la.AccountID == "" {
		return errors.New("invalid account id")
	}
	if !la.CreatedAt.IsZero() && la.UserID == uuid.Nil {
		return errors.New("invalid user id")
	}
	return nil
}

func (la Account) HasType(t Type) bool { return la.Type.Has(t) }

// providerType returns the type for the provider (if known).
func providerType(p provider.Name) Type {
	t := None
	switch p {
	case provider.Unknown:
		break
	case provider.PayPal, provider.WePay:
		t = t.Set(Wallet)
		fallthrough
	case provider.Stripe:
		t = t.Set(Payment)
		fallthrough
	default:
		t = t.Set(Auth)
		break
	}
	return t
}
