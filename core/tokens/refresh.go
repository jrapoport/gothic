package tokens

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/store"
)

// GrantRefreshToken gets or creates a refresh token for the provided user.
func GrantRefreshToken(conn *store.Connection, userID uuid.UUID) (*token.RefreshToken, error) {
	rt, err := grantToken(conn, userID, func() token.Token {
		return token.NewRefreshToken(userID)
	})
	if err != nil {
		return nil, err
	}
	return rt.(*token.RefreshToken), nil
}

// SwapRefreshToken swaps a refresh token for a new one, revoking the previous token.
func SwapRefreshToken(conn *store.Connection, userID uuid.UUID, tok string) (*token.RefreshToken, error) {
	var rt token.Token
	err := conn.Transaction(func(tx *store.Connection) (err error) {
		rt, err = GetUsableRefreshToken(tx, tok)
		if err != nil {
			return err
		}
		// soft delete
		err = tx.Delete(rt).Error
		if err != nil {
			return err
		}
		if rt.IssuedTo() != userID {
			return nil
		}
		rt, err = GrantRefreshToken(tx, rt.IssuedTo())
		return err
	})
	if err != nil {
		return nil, err
	}
	return rt.(*token.RefreshToken), nil
}

// RevokeAllRefreshTokens revokes (deletes) all refresh tokens for a user id.
func RevokeAllRefreshTokens(conn *store.Connection, userID uuid.UUID) error {
	rt := new(token.RefreshToken)
	rt.UserID = userID
	return revokeAll(conn, rt, userID)
}

// HasUsableRefreshToken returns true if a usable refresh token is found.
func HasUsableRefreshToken(conn *store.Connection, userID uuid.UUID) (bool, error) {
	var rt token.RefreshToken
	has, err := conn.Has(&rt, "user_id = ?", userID)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return rt.Usable(), nil
}

// GetUsableRefreshToken returns a usable refresh token for the token string if one exists.
func GetUsableRefreshToken(conn *store.Connection, tok string) (*token.RefreshToken, error) {
	rt := new(token.RefreshToken)
	err := conn.First(rt, "token = ?", tok).Error
	if err != nil {
		return nil, err
	}
	if !rt.Usable() {
		err = errors.New("invalid token")
		return nil, err
	}
	return rt, nil
}
