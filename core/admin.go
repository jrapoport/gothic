package core

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
)

// AdminCreateUser creates a new user
func (a *API) AdminCreateUser(ctx context.Context, email, username, pw string, data types.Map, admin bool) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	uid := ctx.GetAdminID()
	if uid == uuid.Nil {
		err := errors.New("admin user id required")
		return nil, a.logError(err)
	}
	email, err := a.ValidateEmail(email)
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("admin user create: %s", email)
	var u *user.User
	err = a.conn.Transaction(func(tx *store.Connection) error {
		username, err = a.validateUsername(tx, username)
		if err != nil {
			return err
		}
		var adm *user.User
		adm, err = users.GetAuthenticatedUser(tx, uid)
		if err != nil {
			return err
		}
		if !adm.IsAdmin() {
			err = fmt.Errorf("admin required: %s", uid)
			return err
		}
		if admin && adm.Role != user.RoleSuper {
			err = fmt.Errorf("super admin required: %s", uid)
			return err
		}
		u, err = a.userSignup(ctx, tx, a.Provider(), email, username, pw, data)
		if err != nil {
			return err
		}
		if admin {
			err = a.changeRole(ctx, tx, u, user.RoleAdmin)
		}
		return err
	})
	if err != nil {
		return nil, a.logError(err)
	}
	err = a.autoConfirm(ctx, u)
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("admin user signed up:r: %s", u.ID)
	return u, nil
}

// AdminPromoteUser promotes a user to an admin
func (a *API) AdminPromoteUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	uid := ctx.GetAdminID()
	if uid == uuid.Nil {
		err := errors.New("admin user id required")
		return nil, a.logError(err)
	}
	if userID == uuid.Nil {
		err := errors.New("user id required")
		return nil, a.logError(err)
	}
	a.log.Debugf("promote user to admin: %s", userID)
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) error {
		adm, err := users.GetAuthenticatedUser(tx, uid)
		if err != nil {
			return err
		}
		if adm.Role != user.RoleSuper {
			err = fmt.Errorf("super admin required: %s", uid)
			return err
		}
		u, err = users.GetActiveUser(tx, userID)
		if err != nil {
			return err
		}
		return a.changeRole(ctx, tx, u, user.RoleAdmin)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("promoted user: %s", userID)
	return u, nil
}

// AdminDeleteUser deletes a user.
func (a *API) AdminDeleteUser(ctx context.Context, userID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	uid := ctx.GetAdminID()
	if uid == uuid.Nil {
		err := errors.New("context user id required")
		return a.logError(err)
	}
	if userID == uuid.Nil {
		err := errors.New("user id required")
		return a.logError(err)
	}
	a.log.Debugf("delete user: %s", userID)
	err := a.conn.Transaction(func(tx *store.Connection) error {
		adm, err := users.GetAuthenticatedUser(tx, uid)
		if err != nil {
			return err
		}
		if !adm.IsAdmin() {
			err = fmt.Errorf("admin required: %s", uid)
			return err
		}
		u, err := users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		if u.IsAdmin() && adm.Role != user.RoleSuper {
			err = fmt.Errorf("super admin required to delete admin: %s", userID)
			return err
		}
		err = users.DeleteUser(tx, u)
		if err != nil {
			return err
		}
		return audit.LogDeleted(ctx, tx, userID)
	})
	if err != nil {
		return a.logError(err)
	}
	a.log.Debugf("deleted user: %s", userID)
	return nil
}
