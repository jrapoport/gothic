package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// GetUser returns a user for the id if found.
func (a *API) GetUser(userID uuid.UUID) (*user.User, error) {
	u, err := users.GetUser(a.conn, userID)
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

// GetAuthenticatedUser returns an authenticated user for the id if found.
func (a *API) GetAuthenticatedUser(userID uuid.UUID) (*user.User, error) {
	u, err := users.GetAuthenticatedUser(a.conn, userID)
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

// GetUserWithEmail returns a user for the email if found.
func (a *API) GetUserWithEmail(email string) (*user.User, error) {
	u, err := users.GetUserWithEmail(a.conn, email)
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

// SearchUsers searches for users.
func (a *API) SearchUsers(ctx context.Context, f store.Filters, page *store.Pagination) ([]*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return users.SearchUsers(a.conn, ctx.GetSort(), f, page)
}

// ChangePassword changes a user password.
func (a *API) ChangePassword(ctx context.Context, userID uuid.UUID, old, pw string) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := validate.Password(a.config, pw)
	if err != nil {
		return nil, a.logError(err)
	}
	var u *user.User
	err = a.conn.Transaction(func(tx *store.Connection) (err error) {
		u, err = users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		if u.IsLocked() {
			err = fmt.Errorf("invalid user %s", u.ID)
			return err
		}
		err = u.Authenticate(old)
		if err != nil {
			err = fmt.Errorf("incorrect password %w", err)
			return err
		}
		err = users.ChangePassword(tx, u, pw)
		if err != nil {
			return err
		}
		return audit.LogPasswordChange(ctx, tx, u.ID)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

// ChangeRole changes a user role.
func (a *API) ChangeRole(ctx context.Context, userID uuid.UUID, r user.Role) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	a.log.Debugf("change user to %s: %s", r.String(), userID)
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		u, err = users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		return a.changeRole(ctx, tx, u, r)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("changed user to %s: %s", r.String(), userID)
	return u, nil
}

func (a *API) changeRole(ctx context.Context, conn *store.Connection, u *user.User, r user.Role) error {
	if u == nil || u.IsLocked() {
		err := errors.New("invalid user")
		return err
	}
	if u.Role == r {
		return nil
	}
	switch r {
	case user.RoleUser:
		break
	case user.RoleAdmin:
		if u.Role >= r {
			a.log.Warnf("user is already %s: %s",
				r.String(), u.ID)
			return nil
		}
	default:
		err := fmt.Errorf("invalid role: %s",
			r.String())
		return err
	}
	a.log.Debugf("change user to %s: %s",
		r.String(), u.ID)
	err := conn.Transaction(func(tx *store.Connection) (err error) {
		err = users.ChangeRole(tx, u, r)
		if err != nil {
			return err
		}
		return audit.LogChangeRole(ctx, tx, u.ID, u.Role)
	})
	if err != nil {
		return err
	}
	a.log.Debugf("changed user to %s: %s",
		r.String(), u.ID)
	return nil
}

// ConfirmUser confirms a user account.
func (a *API) ConfirmUser(ctx context.Context, tok string) (*user.User, error) {
	return a.confirmUserWithChanges(ctx, tok, nil)
}

// ConfirmPasswordChange confirms a user pw change & account (if needed).
func (a *API) ConfirmPasswordChange(ctx context.Context, tok string, pw string) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := validate.Password(a.config, pw)
	if err != nil {
		return nil, a.logError(err)
	}
	return a.confirmUserWithChanges(ctx, tok,
		func(tx *store.Connection, ct *token.ConfirmToken, u *user.User) error {
			err = users.ChangePassword(tx, u, pw)
			if err != nil {
				return err
			}
			return audit.LogPasswordChange(ctx, tx, u.ID)
		})
}

// ConfirmChangeEmail confirms a user email change & account (if needed).
func (a *API) ConfirmChangeEmail(ctx context.Context, tok string, email string) (*user.User, error) {
	var err error
	email, err = a.ValidateEmail(email)
	if err != nil {
		return nil, a.logError(err)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return a.confirmUserWithChanges(ctx, tok,
		func(tx *store.Connection, ct *token.ConfirmToken, u *user.User) error {
			err = users.ChangeEmail(tx, u, email)
			if err != nil {
				return err
			}
			return audit.LogEmailChange(ctx, tx, u.ID)
		})
}

// UpdateUser updates a user.
func (a *API) UpdateUser(ctx context.Context, userID uuid.UUID, username *string, data types.Map) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if username != nil {
		err := validate.Username(a.config, *username)
		if err != nil {
			return nil, a.logError(err)
		}
	}
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		u, err = users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		if !u.IsActive() {
			err = fmt.Errorf("invalid user: %s", userID)
			return err
		}
		err = users.Update(tx, u, username, data)
		if err != nil {
			return err
		}
		return audit.LogUserUpdated(ctx, tx, u.ID)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

type changesFunc func(tx *store.Connection, ct *token.ConfirmToken, u *user.User) error

func (a *API) confirmUserWithChanges(ctx context.Context, token string, changes changesFunc) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx.SetProvider(a.Provider())
	ip := ctx.GetIPAddress()
	var u *user.User
	var confirmed bool
	err := a.conn.Transaction(func(tx *store.Connection) error {
		ct, err := tokens.GetConfirmToken(tx, token)
		if err != nil {
			return err
		}
		if !ct.Usable() {
			return errors.New("invalid token")
		}
		u, err = users.GetUser(tx, ct.UserID)
		if err != nil {
			return err
		}
		confirmed, err = users.ConfirmIfNeeded(tx, ct, u)
		if err != nil {
			return err
		}
		if confirmed {
			err = audit.LogConfirmed(ctx, tx, u.ID)
			if err != nil {
				return err
			}
		}
		if changes != nil {
			err = changes(tx, ct, u)
		}
		return err
	})
	if err != nil {
		return nil, a.logError(err)
	}
	if confirmed {
		// was the user was confirmed for the first time?
		a.dispatchEvent(events.Confirmed, types.Map{
			key.Provider:  u.Provider,
			key.UserID:    u.ID,
			key.IPAddress: ip,
			key.Timestamp: time.Now().UTC(),
		})
	}
	return u, nil
}

// BanUser bans a user.
func (a *API) BanUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if userID == uuid.Nil {
		err := errors.New("user id required")
		return nil, a.logError(err)
	}
	a.log.Debugf("ban user: %s", userID)
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		u, err = users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		err = users.BanUser(tx, u)
		if err != nil {
			return err
		}
		return audit.LogBanned(ctx, tx, userID)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("banned user: %s", userID)
	return u, nil
}

// DeleteUser deletes a user.
func (a *API) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if userID == uuid.Nil {
		err := errors.New("user id required")
		return a.logError(err)
	}
	a.log.Debugf("delete user: %s", userID)
	err := a.conn.Transaction(func(tx *store.Connection) error {
		u, err := users.GetUser(tx, userID)
		if err != nil {
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
