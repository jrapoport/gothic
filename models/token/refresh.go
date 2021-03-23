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
	at := *NewAccessToken(utils.SecureToken(), InfiniteUse, NoExpiration)
	at.UserID = userID
	return &RefreshToken{at}
}

// Class returns the class of the refresh token.
func (rt RefreshToken) Class() Class {
	return Refresh
}

// Usable returns true if the token is usable.
func (rt RefreshToken) Usable() bool {
	if rt.CreatedAt.IsZero() {
		return false
	}
	return rt.AccessToken.Usable()
}

// HasToken returns true if the refresh token is found.
func (rt RefreshToken) HasToken(tx *store.Connection) (bool, error) {
	if rt.Token == "" {
		return false, errors.New("invalid token")
	}
	return tx.Has(&rt, "token = ?", rt.Token)
}
