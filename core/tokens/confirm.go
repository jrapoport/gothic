package tokens

import (
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/store"
)

// GrantConfirmToken gets or creates a confirmation token for the provided user.
func GrantConfirmToken(conn *store.Connection, userID uuid.UUID, exp time.Duration) (*token.ConfirmToken, error) {
	t, err := grantToken(conn, userID, func() token.Token {
		return token.NewConfirmToken(userID, exp)
	})
	if err != nil {
		return nil, err
	}
	return t.(*token.ConfirmToken), nil
}

// GetConfirmToken returns the confirmation token for the token string if found.
func GetConfirmToken(conn *store.Connection, tok string) (*token.ConfirmToken, error) {
	var ct token.ConfirmToken
	err := conn.First(&ct, "token = ?", tok).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

// ConfirmTokenSent marks a confirmation token as sent.
func ConfirmTokenSent(conn *store.Connection, ct *token.ConfirmToken) error {
	now := time.Now().UTC()
	ct.SentAt = &now
	return conn.Model(ct).Update("sent_at", ct.SentAt).Error
}

// GetLastConfirmTokenSent returns the last confirmation token that was sent.
func GetLastConfirmTokenSent(conn *store.Connection, userID uuid.UUID) (*token.ConfirmToken, error) {
	ct := new(token.ConfirmToken)
	has, err := conn.HasLast(ct, "user_id = ? AND sent_at NOT NULL", userID)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return ct, nil
}
