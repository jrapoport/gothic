package accounts

import (
	"dario.cat/mergo"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"

	// make sure the fk gets migrated correctly
	_ "github.com/jrapoport/gothic/models/code"
)

// GetAccount returns a linked account.
func GetAccount(conn *store.Connection, p provider.Name, accountID string) (*account.Account, error) {
	const linkQuery = key.Provider + " = ? AND " + key.AccountID + " = ?"
	var la account.Account
	err := conn.First(&la, linkQuery, p, accountID).Error
	if err != nil {
		return nil, err
	}
	return &la, nil
}

// HasAccount returns true if the linked account exists.
func HasAccount(conn *store.Connection, p provider.Name, accountID string) (bool, error) {
	const linkQuery = key.Provider + " = ? AND " + key.AccountID + " = ?"
	has, err := conn.Has(new(account.Account), linkQuery, p, accountID)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

// UpdateAccount updates an external account
func UpdateAccount(conn *store.Connection,
	la *account.Account, email *string, data types.Map) (bool, error) {
	if email != nil && la.Email != *email {
		la.Email = *email
	} else if data == nil || len(data) <= 0 {
		return false, nil
	}
	err := mergo.Map(&la.Data, data, mergo.WithOverride)
	if err != nil {
		return false, err
	}
	err = conn.Model(la).Select(key.Email, key.Data).Updates(la).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
