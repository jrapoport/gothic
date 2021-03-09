package token

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

func init() {
	store.AddAutoMigrationWithIndexes("4000-refresh_tokens",
		RefreshToken{}, AccessTokenIndexes)
}

// RefreshToken holds a refresh token.
type RefreshToken struct {
	AccessToken
}

var _ Token = (*RefreshToken)(nil)

// NewRefreshToken generates a new token. The type is inferred from uses.
func NewRefreshToken(userID uuid.UUID) *RefreshToken {
	t := *NewAccessToken(utils.SecureToken(), InfiniteUse, NoExpiration)
	t.UserID = userID
	return &RefreshToken{t}
}

// Class returns the class of the refresh token.
func (t RefreshToken) Class() Class {
	return Refresh
}

// Usable returns true if the token is usable.
func (t RefreshToken) Usable() bool {
	if t.CreatedAt.IsZero() {
		return false
	}
	return t.AccessToken.Usable()
}

// HasToken returns true if the refresh token is found.
func (t RefreshToken) HasToken(tx *store.Connection) (bool, error) {
	if t.Token == "" {
		return false, errors.New("invalid token")
	}
	return tx.Has(&t, "token = ?", t.Token)
}
