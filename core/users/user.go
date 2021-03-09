package users

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

// GetUser get a user with the matching id.
func GetUser(conn *store.Connection, userID uuid.UUID) (*user.User, error) {
	if userID == user.SystemID {
		return nil, errors.New("invalid user")
	}
	u := new(user.User)
	err := conn.First(u, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetActiveUser get an active user with the matching id.
func GetActiveUser(conn *store.Connection, userID uuid.UUID) (*user.User, error) {
	u, err := GetUser(conn, userID)
	if err != nil {
		return nil, err
	}
	if !u.IsActive() {
		err = fmt.Errorf("invalid user: %s", userID)
		return nil, err
	}
	return u, nil
}

// GetActiveUserWithRole get an active user with the minimum role.
func GetActiveUserWithRole(conn *store.Connection, userID uuid.UUID, r user.Role) (*user.User, error) {
	u, err := GetActiveUser(conn, userID)
	if err != nil {
		return nil, err
	}
	if u.Role < r {
		err = fmt.Errorf("insuffienct role %s", u.Role)
		return nil, err
	}
	return u, nil
}

// GetAuthenticatedUser get an authenticated user with the matching id.
func GetAuthenticatedUser(conn *store.Connection, userID uuid.UUID) (*user.User, error) {
	var u *user.User
	if userID == user.SystemID {
		return nil, errors.New("invalid user")
	}
	err := conn.Transaction(func(tx *store.Connection) error {
		has, err := tokens.HasUsableRefreshToken(tx, userID)
		if err != nil {
			return err
		}
		if !has {
			err = fmt.Errorf("valid refresh token not found: %s",
				userID.String())
			return err
		}
		u, err = GetActiveUser(tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetUserWithEmail get a user with the matching email.
func GetUserWithEmail(conn *store.Connection, email string) (*user.User, error) {
	e, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}
	u := new(user.User)
	err = conn.First(u, "email = ?", e.Address).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

// HasUserWithEmail get a user with the matching email. Unlike
// GetUserWithEmail, no error is returned if not found.
func HasUserWithEmail(conn *store.Connection, email string) (*user.User, error) {
	e, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}
	u := new(user.User)
	has, err := conn.Has(u, "email = ?", e.Address)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return u, nil
}

// IsEmailTaken returns true if the email address is already taken
func IsEmailTaken(conn *store.Connection, email string) (bool, error) {
	if email == "" {
		return false, errors.New("invalid email")
	}
	return conn.Has(new(user.User), "email = ?", email)
}

// IsUsernameTaken returns true if the username is already taken
func IsUsernameTaken(conn *store.Connection, username string) (bool, error) {
	if username == "" {
		return false, errors.New("invalid username")
	}
	return conn.Has(new(user.User), "username = ?", username)
}

// RandomUsername returns a random username of max random length (28).
func RandomUsername(conn *store.Connection, unique bool) (string, error) {
	username := ""
	for {
		username = utils.RandomUsername()
		if !unique {
			break
		}
		taken, err := IsUsernameTaken(conn, username)
		if err != nil {
			return "", err
		}
		if !taken {
			break
		}
	}
	return username, nil
}
