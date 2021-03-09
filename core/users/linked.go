package users

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// LinkAccount links an account to a user.
func LinkAccount(conn *store.Connection, u *user.User, la *account.Linked) error {
	if u == nil || u.IsLocked() {
		return errors.New("invalid user")
	}
	if !la.CreatedAt.IsZero() {
		return errors.New("account already linked")
	}
	if la.UserID != uuid.Nil {
		return errors.New("user id must be nil")
	}
	if err := la.Valid(); err != nil {
		return err
	}
	return conn.Model(u).Association("Linked").Append(la)
}

// HasLinkedUser returns the user linked to the account.
func HasLinkedUser(conn *store.Connection, p provider.Name, accountID string) (*user.User, error) {
	var u *user.User
	err := conn.Transaction(func(tx *store.Connection) error {
		var la account.Linked
		const query = "provider = ? AND account_id = ?"
		has, err := tx.Has(&la, query, p, accountID)
		if err != nil || !has {
			return err
		}
		u, err = GetUser(tx, la.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}
