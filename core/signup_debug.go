//go:build !release
// +build !release

package core

import (
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

func (a *API) debugSignup(tx *store.Connection, email, code string) bool {
	if !utils.IsDebugPIN(code) {
		return false
	}
	if a.config.Signup.AutoConfirm {
		return true
	}
	u, err := users.GetUserWithEmail(tx, email)
	if err != nil {
		a.log.Error(err)
		return false
	}
	u.Status = user.Restricted
	u.ConfirmedAt = nil
	err = tx.Save(u).Error
	if err != nil {
		a.log.Error(err)
		return false
	}
	return true
}
