package users

import (
	"errors"

	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/providers"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/utils"
)

// CreateUser creates a user.
func CreateUser(conn *store.Connection, p provider.Name, email, username, pw string, data, meta types.Map) (*user.User, error) {
	err := providers.IsEnabled(p)
	if err != nil {
		return nil, errors.New("invalid provider")
	}
	email, err = validate.Email(email)
	if err != nil {
		return nil, err
	}
	var u *user.User
	err = conn.Transaction(func(tx *store.Connection) (err error) {
		var taken bool
		// has this user has already signed up?
		taken, err = IsEmailTaken(tx, email)
		if err != nil {
			return err
		}
		if taken {
			return errors.New("email taken")
		}
		u, err = createUser(tx, p, email, username, pw, data, meta)
		return err
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func createUser(conn *store.Connection, p provider.Name, email, username, pw string, data, sys types.Map) (*user.User, error) {
	hashed, err := utils.HashPassword(pw)
	if err != nil {
		return nil, err
	}
	u := user.NewUser(p, user.RoleUser, email, username, hashed, data, sys)
	err = conn.Create(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}