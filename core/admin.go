package core

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

//////////////////////////////////
//////////////////////////////////
// NOTE:
// APIs in this file require admin
// or super admin user permissions
//
//

// CreateUser creates a new user
// NOTE: This API requires admin user permissions
func (a *API) CreateUser(ctx context.Context, email, username, pw string, data types.Map, admin bool) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !ctx.IsAdmin() {
		err := errors.New("admin user required")
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
		aid := ctx.AdminID()
		role, err := a.validateAdmin(tx, aid)
		if err != nil {
			return err
		}
		if admin && role != user.RoleSuper {
			err = fmt.Errorf("super admin required: %s", aid)
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

// ChangeRole changes a user's role
// NOTE: This API requires admin user permissions
func (a *API) ChangeRole(ctx context.Context, userID uuid.UUID, role user.Role) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !ctx.IsAdmin() {
		err := errors.New("admin user required")
		return nil, a.logError(err)
	}
	if userID == uuid.Nil || userID == user.SystemID || userID == user.SuperAdminID {
		err := errors.New("user id required")
		return nil, a.logError(err)
	}
	if role == user.InvalidRole || role == user.RoleSuper {
		err := fmt.Errorf("invalid role change: %s", role.String())
		return nil, a.logError(err)
	}
	a.log.Debugf("change user %s to %s", userID, role.String())
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) error {
		aid := ctx.AdminID()
		adminRole, err := a.validateAdmin(tx, aid)
		if err != nil {
			return err
		}
		u, err = users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		if u.Role == role {
			return nil
		}
		switch u.Role {
		case user.RoleUser:
			// only super admins can promote other admins
			if role == user.RoleAdmin && adminRole < user.RoleSuper {
				err = fmt.Errorf("super admin required: %s", aid)
				return err
			}
			break
		case user.RoleAdmin:
			// only super admins can demote other admins
			if role == user.RoleUser && adminRole < user.RoleSuper {
				err = fmt.Errorf("super admin required: %s", aid)
				return err
			}
			break
		case user.RoleSuper:
			fallthrough
		default:
			err = fmt.Errorf("%s user role change forbidden", u.Role.String())
			return err
		}
		return a.changeRole(ctx, tx, u, role)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("changed user %s to %s", userID, role.String())
	return u, nil
}

// PromoteUser promotes a user to an admin
// NOTE: This API requires admin user permissions
func (a *API) PromoteUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !ctx.IsAdmin() {
		err := errors.New("admin user required")
		return nil, a.logError(err)
	}
	if userID == uuid.Nil {
		err := errors.New("user id required")
		return nil, a.logError(err)
	}
	a.log.Debugf("promote user to admin: %s", userID)
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) error {
		aid := ctx.AdminID()
		role, err := a.validateAdmin(tx, aid)
		if err != nil {
			return err
		}
		if role != user.RoleSuper {
			err = fmt.Errorf("super admin required: %s", aid)
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

// DeleteUser deletes a user
// NOTE: This API requires admin user permissions
func (a *API) DeleteUser(ctx context.Context, userID uuid.UUID, hard bool) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if !ctx.IsAdmin() {
		err := errors.New("admin user required")
		return a.logError(err)
	}
	if userID == uuid.Nil {
		err := errors.New("user id required")
		return a.logError(err)
	}
	a.log.Debugf("delete user: %s", userID)
	err := a.conn.Transaction(func(tx *store.Connection) error {
		role, err := a.validateAdmin(tx, ctx.AdminID())
		if err != nil {
			return err
		}
		u, err := users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		if u.IsAdmin() && role != user.RoleSuper {
			err = fmt.Errorf("super admin required to delete admin: %s", userID)
			return err
		}
		err = users.DeleteUser(tx, u, hard)
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

// ValidateAdmin validates the user id as an admin and returns the role
func (a *API) ValidateAdmin(aid uuid.UUID) (user.Role, error) {
	return a.validateAdmin(a.conn, aid)
}

func (a *API) validateAdmin(tx *store.Connection, aid uuid.UUID) (user.Role, error) {
	if aid == user.SuperAdminID {
		return user.RoleSuper, nil
	}
	adm, err := users.GetAuthenticatedUser(tx, aid)
	if err != nil {
		return user.InvalidRole, err
	}
	if !adm.IsAdmin() {
		err = fmt.Errorf("admin required: %s", aid)
		return user.InvalidRole, err
	}
	return adm.Role, nil
}
