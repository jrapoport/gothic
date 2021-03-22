package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

func init() {
	store.AddAutoMigrationWithIndexes("3000-confirm_tokens",
		ConfirmToken{}, AccessTokenIndexes)
}

// ConfirmToken holds a confirmation token.
type ConfirmToken struct {
	AccessToken
	SentAt *time.Time `json:"sent_at"`
}

var _ Token = (*ConfirmToken)(nil)

// NewConfirmToken generates a new token. The type is inferred from uses.
func NewConfirmToken(userID uuid.UUID, exp time.Duration) *ConfirmToken {
	at := *NewAccessToken(utils.SecureToken(), SingleUse, exp)
	at.UserID = userID
	return &ConfirmToken{AccessToken: at}
}

// Class returns the class of the confirmation token.
func (ct ConfirmToken) Class() Class {
	return Confirm
}

// Usable returns true if the token is usable.
func (ct ConfirmToken) Usable() bool {
	if ct.CreatedAt.IsZero() {
		return false
	}
	return ct.AccessToken.Usable()
}

// HasToken returns true if the refresh token is found.
func (ct ConfirmToken) HasToken(tx *store.Connection) (bool, error) {
	if ct.Token == "" {
		return false, errors.New("invalid token")
	}
	return tx.Has(&ct, "token = ?", ct.Token)
}
