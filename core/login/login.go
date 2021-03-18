package login

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// UserLogin authorizes a user and returns a bearer token
func UserLogin(conn *store.Connection, p provider.Name, email, pw string) (*user.User, error) {
	email, err := validate.Email(email)
	if err != nil {
		return nil, err
	}
	var u *user.User
	err = conn.Transaction(func(tx *store.Connection) error {
		u, err = users.GetUserWithEmail(tx, email)
		if err != nil {
			return err
		}
		if !u.IsActive() {
			return errors.New("inactive user")
		}
		if u.Provider != p {
			return errors.New("invalid provider")
		}
		err = u.Authenticate(pw)
		if err != nil {
			err = fmt.Errorf("incorrect password %w", err)
			return err
		}
		now := time.Now().UTC()
		u.LoginAt = &now
		return tx.Model(u).Update("login_at", u.LoginAt).Error
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UserLogout logs a user out and revokes all their refresh tokens.
func UserLogout(conn *store.Connection, userID uuid.UUID) error {
	return tokens.RevokeAllRefreshTokens(conn, userID)
}
