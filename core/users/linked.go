package users

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/accounts"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
)

// LinkAccount links an account to a user.
func LinkAccount(conn *store.Connection, userID uuid.UUID, link *account.Account) error {
	if link == nil {
		return errors.New("account required")
	}
	err := link.Valid()
	if err != nil {
		return err
	}
	return conn.Transaction(func(tx *store.Connection) error {
		u, err := GetUser(tx, userID)
		if err != nil {
			return err
		}
		if u == nil || u.IsLocked() {
			return errors.New("invalid user")
		}
		has, err := accounts.HasAccount(tx, link.Provider, link.AccountID)
		if err != nil {
			return err
		}
		if has {
			return errors.New("account already linked")
		}
		return tx.Model(u).Association("Linked").Append(link)
	})
}

// GetLinkedAccounts returned the linked accounts for a user.
func GetLinkedAccounts(conn *store.Connection,
	userID uuid.UUID, t account.Type, f store.Filters) ([]*account.Account, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user id required")
	}
	filters := f.Copy()
	flt := store.Filter{
		Filters:   filters,
		DataField: key.Data,
		Fields: []string{
			key.AccountID,
			key.Email,
			key.Provider,
		},
	}
	var linked []*account.Account
	err := conn.Transaction(func(tx *store.Connection) error {
		u, err := GetUser(tx, userID)
		if err != nil {
			return err
		}
		if u == nil || u.IsLocked() {
			return errors.New("invalid user")
		}
		db := tx.Model(new(account.Account))
		db.Where(key.UserID+" = ?", userID)
		if t != account.None && t != account.All {
			db.Where(key.Type+" & ? != 0", t)
		}
		err = store.Search(db, &linked, store.Descending, flt, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return linked, err
}
