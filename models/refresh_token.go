package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jrapoport/gothic/crypto"
	"github.com/jrapoport/gothic/storage"
	"github.com/pkg/errors"
)

func init() {
	storage.AddMigration(&RefreshToken{})
}

// RefreshToken is the database model for refresh tokens.
type RefreshToken struct {
	ID int64 `gorm:"primaryKey"`

	Token  string    `gorm:"index:refresh_tokens_token_idx;type:varchar(255) DEFAULT NULL"`
	UserID uuid.UUID `gorm:"index:user_id_idx;type:varchar(255) DEFAULT NULL"`

	Revoked   bool      `gorm:"type:tinyint(1) DEFAULT NULL"`
	CreatedAt time.Time `gorm:"type:timestamp NULL DEFAULT NULL"`
	UpdatedAt time.Time `gorm:"type:timestamp NULL DEFAULT NULL"`
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
			return errors.Wrap(terr, "error creating audit log entry")
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
	token := &RefreshToken{
		UserID: user.ID,
		Token:  crypto.SecureToken(),
	}

	if err := tx.Create(token).Error; err != nil {
		return nil, errors.Wrap(err, "error creating refresh token")
	}
	return token, nil
}
