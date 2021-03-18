package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/login"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
)

// Login is the endpoint for logging in a user
func (a *API) Login(ctx context.Context, email, pw string) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ip := ctx.GetIPAddress()
	recaptcha := ctx.GetReCaptcha()
	p := a.Provider()
	ctx.SetProvider(p)
	a.log.Debugf("login: %s %s (%s %s %s)",
		email, pw, p, ip, recaptcha)
	err := a.ext.IsEnabled(p)
	if err != nil {
		return nil, a.logError(err)
	}
	email, err = a.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	if a.config.Recaptcha.Login {
		// if recaptcha is disabled this is a no-op
		err = validate.ReCaptcha(a.config, ip, recaptcha)
		if err != nil {
			return nil, a.logError(err)
		}
	}
	var u *user.User
	err = a.conn.Transaction(func(tx *store.Connection) error {
		u, err = login.UserLogin(tx, p, email, pw)
		if err != nil {
			return err
		}
		return audit.LogLogin(ctx, tx, u.ID)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	a.dispatchEvent(events.Login, types.Map{
		key.Provider:  p,
		key.IPAddress: ip,
		key.UserID:    u.ID,
		key.Timestamp: time.Now().UTC(),
	})
	return u, nil
}

// Logout revokes all refresh token for a user id.
func (a *API) Logout(ctx context.Context, userID uuid.UUID) error {
	if ctx == nil {
		ctx = context.Background()
	}
	p := ctx.GetProvider()
	ip := ctx.GetIPAddress()
	err := a.conn.Transaction(func(tx *store.Connection) error {
		err := login.UserLogout(tx, userID)
		if err != nil {
			return err
		}
		return audit.LogLogout(ctx, tx, userID)
	})
	if err != nil {
		return a.logError(err)
	}
	a.dispatchEvent(events.Logout, types.Map{
		key.Provider:  p,
		key.IPAddress: ip,
		key.UserID:    userID,
		key.Timestamp: time.Now().UTC(),
	})
	return nil
}
