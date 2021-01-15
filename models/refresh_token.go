package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/crypto"
	"github.com/jrapoport/gothic/storage"
	"gorm.io/gorm"
)

func init() {
	storage.AddMigration(&RefreshToken{})
}

// RefreshToken is the database model for refresh tokens.
type RefreshToken struct {
	gorm.Model
	Token   string    `gorm:"index:refresh_tokens_token_idx;type:varchar(255)"`
	UserID  uuid.UUID `gorm:"index:user_id_idx;"`
	Revoked bool
}

// GrantAuthenticatedUser creates a refresh token for the provided user.
func GrantAuthenticatedUser(tx *storage.Connection, user *User) (*RefreshToken, error) {
	return createRefreshToken(tx, user)
}

// GrantRefreshTokenSwap swaps a refresh token for a new one, revoking the provided token.
func GrantRefreshTokenSwap(tx *storage.Connection, user *User, token *RefreshToken) (*RefreshToken, error) {
	var newToken *RefreshToken
	err := tx.Transaction(func(rtx *storage.Connection) error {
		var terr error
		if terr = NewAuditLogEntry(rtx, user, TokenRevokedAction, nil); terr != nil {
			terr = fmt.Errorf("%w granting refresh token", terr)
			return terr
		}

		token.Revoked = true
		if terr = rtx.Model(&token).Select("revoked").Updates(token).Error; terr != nil {
			return terr
		}
		newToken, terr = createRefreshToken(rtx, user)
		return terr
	})
	return newToken, err
}

// Logout deletes all refresh tokens for a user.
func Logout(tx *storage.Connection, uid uuid.UUID) error {
	return tx.Where("user_id = ?", uid).Delete(&RefreshToken{}).Error
}

func createRefreshToken(tx *storage.Connection, user *User) (*RefreshToken, error) {
	t := &RefreshToken{
		UserID: user.ID,
		Token:  crypto.SecureToken(),
	}

	if err := tx.Create(t).Error; err != nil {
		err = fmt.Errorf("%w creating refresh token", err)
		return nil, err
	}
	return t, nil
}
