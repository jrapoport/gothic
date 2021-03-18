package token

import (
	"errors"
	"time"

	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

func init() {
	indexes := append(AccessTokenIndexes, "idx_provider")
	store.AddAutoMigrationWithIndexes("3000-auth_tokens",
		AuthToken{}, indexes)
}

// AuthToken holds an auth token.
type AuthToken struct {
	AccessToken
	Provider provider.Name `json:"provider" gorm:"index:idx_provider;type:char(255)"`
}

var _ Token = (*AuthToken)(nil)

// NewAuthToken generates a new token. The type is inferred from uses.
func NewAuthToken(p provider.Name, exp time.Duration) *AuthToken {
	t := *NewAccessToken(utils.SecureToken(), SingleUse, exp)
	t.UserID = p.ID()
	return &AuthToken{AccessToken: t, Provider: p}
}

// Class returns the class of the auth token.
func (t AuthToken) Class() Class {
	return Auth
}

// HasToken returns true if the auth token is found.
func (t AuthToken) HasToken(tx *store.Connection) (bool, error) {
	if t.Token == "" {
		return false, errors.New("invalid token")
	}
	return tx.Has(&t, "token = ?", t.Token)
}
