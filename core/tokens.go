package core

import (
	"errors"
	"fmt"
	"github.com/jrapoport/gothic/config"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// GrantBearerToken issues a bearer token for the user.
func (a *API) GrantBearerToken(ctx context.Context, u *user.User) (*tokens.BearerToken, error) {
	if u == nil || (!u.IsActive() && !u.IsRestricted()) {
		err := errors.New("invalid user")
		return nil, a.logError(err)
	}
	var bt *tokens.BearerToken
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		bt, err = tokens.GrantBearerToken(a.conn, a.config.JWT, u)
		if err != nil {
			return err
		}
		return audit.LogTokenGranted(ctx, tx, bt)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	return bt, nil
}

// RefreshBearerToken refreshes the bearer token for the user.
func (a *API) RefreshBearerToken(ctx context.Context, refreshToken string) (*tokens.BearerToken, error) {
	var bt *tokens.BearerToken
	err := a.conn.Transaction(func(tx *store.Connection) error {
		rt, err := tokens.GetUsableRefreshToken(tx, refreshToken)
		if err != nil {
			return err
		}
		u, err := users.GetUser(tx, rt.UserID)
		if err != nil {
			return err
		}
		if !u.IsActive() {
			err = fmt.Errorf("inactive user %s", u.ID)
			return err
		}
		bt, err = tokens.RefreshBearerToken(tx, a.config.JWT, u, rt.String())
		if err != nil {
			return err
		}
		return audit.LogTokenRefreshed(ctx, tx, bt)
	})
	if err != nil {
		return nil, a.logError(err)
	}
	return bt, nil
}

type (
	checkSenderFunc func(u *user.User) (bool, error)
	sendConfirmFunc func(to string, ct *token.ConfirmToken) error
)

func (a *API) sendConfirmToken(ctx context.Context, userID uuid.UUID, check checkSenderFunc, send sendConfirmFunc) error {
	if a.mail == nil || a.mail.IsOffline() {
		a.log.Warn("mail not found")
		return nil
	}
	if ctx.GetProvider().IsExternal() {
		// log the error but hide it so we don't leak information
		err := fmt.Errorf("invalid provider: %s", ctx.GetProvider())
		return err
	}
	err := a.conn.Transaction(func(tx *store.Connection) error {
		u, err := users.GetUser(tx, userID)
		if err != nil {
			return err
		}
		// double check that this is an internal user provider
		if u.Provider.IsExternal() {
			// log the error but hide it so we don't leak information
			err = fmt.Errorf("invalid provider: %s", u.Provider)
			return err
		}
		ok, err := check(u)
		if !ok || err != nil {
			return err
		}
		ct, err := tokens.GrantConfirmToken(tx, u.ID, a.config.Mail.Expiration)
		if err != nil {
			return err
		}
		limit := a.config.Mail.SendLimit
		if ct.SentAt != nil && time.Now().UTC().Before(ct.SentAt.Add(limit)) {
			a.log.Warnf("rate limit exceeded for user: %s", u.ID)
			return config.ErrRateLimitExceeded
		}
		err = send(u.EmailAddress().String(), ct)
		if err != nil {
			return err
		}
		err = tokens.ConfirmTokenSent(tx, ct)
		if err != nil {
			return err
		}
		return audit.LogConfirmSent(ctx, tx, ct)

	})
	return a.logError(err)
}
