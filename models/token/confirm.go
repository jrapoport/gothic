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
	t := *NewAccessToken(utils.SecureToken(), SingleUse, exp)
	t.UserID = userID
	return &ConfirmToken{AccessToken: t}
}

// Class returns the class of the confirmation token.
func (t ConfirmToken) Class() Class {
	return Confirm
}

// Usable returns true if the token is usable.
func (t ConfirmToken) Usable() bool {
	if t.CreatedAt.IsZero() {
		return false
	}
	return t.AccessToken.Usable()
}

// HasToken returns true if the refresh token is found.
func (t ConfirmToken) HasToken(tx *store.Connection) (bool, error) {
	if t.Token == "" {
		return false, errors.New("invalid token")
	}
	return tx.Has(&t, "token = ?", t.Token)
}
