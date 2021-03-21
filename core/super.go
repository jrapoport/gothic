package core

import (
	"fmt"

	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
)

// CreateSuperAdmin creates a built-in super admin account
func (a *API) CreateSuperAdmin() error {
	pw := a.config.RootPassword
	su := user.NewSuperAdmin(pw)
	su.Provider = a.superProvider()
	err := a.conn.FirstOrCreate(&su).Error
	if err != nil {
		return a.logError(err)
	}
	if su.Provider != a.superProvider() {
		err = fmt.Errorf("invalid provider: %s", su.Provider)
		return a.logError(err)
	}
	err = su.Authenticate(pw)
	return a.logError(err)
}

// GetSuperAdmin gets the built-in super admin account
func (a *API) GetSuperAdmin(pw string) (*user.User, error) {
	su, err := users.GetUser(a.conn, user.SuperAdminID)
	if err != nil {
		return nil, a.logError(err)
	}
	if su.Provider != a.superProvider() {
		err = fmt.Errorf("invalid provider: %s", su.Provider)
		return nil, a.logError(err)
	}
	if pw == "" {
		pw = a.config.RootPassword
	}
	err = su.Authenticate(pw)
	if err != nil {
		return nil, a.logError(err)
	}
	return su, nil
}

func (a *API) superProvider() provider.Name {
	return provider.NormalizeName(a.config.Name)
}
