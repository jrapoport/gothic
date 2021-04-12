package core

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// CreateSignupCode returns a new signup pin code.
func (a *API) CreateSignupCode(ctx context.Context, uses int) (string, error) {
	var sc *code.SignupCode
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		sc, err = codes.CreateSignupCode(a.conn, user.SystemID, code.PIN, uses, true)
		if err != nil {
			return err
		}
		return audit.LogTokenGranted(ctx, tx, sc)
	})
	if err != nil {
		return "", a.logError(err)
	}
	return sc.Code(), nil
}

// CreateSignupCodes returns a list of new signup pin codes.
func (a *API) CreateSignupCodes(ctx context.Context, uses, count int) ([]string, error) {
	aid := ctx.AdminID()
	if aid == uuid.Nil {
		err := errors.New("admin user id required")
		return nil, a.logError(err)
	}
	var list []string
	err := a.conn.Transaction(func(tx *store.Connection) error {
		cl, err := codes.CreateSignupCodes(tx, aid, code.PIN, uses, count)
		if err != nil {
			return err
		}
		for _, sc := range cl {
			list = append(list, sc.Code())
			err = audit.LogTokenGranted(ctx, tx, sc)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, a.logError(err)
	}
	return list, nil
}

// CheckSignupCode returns nil if the code is valid and usable.
func (a *API) CheckSignupCode(tok string) (*code.SignupCode, error) {
	return codes.GetUsableSignupCode(a.conn, tok)
}

// DeleteSignupCode returns nil if the code is valid and usable.
func (a *API) DeleteSignupCode(tok string) error {
	return codes.DeleteSignupCode(a.conn, tok)
}

type (
	createCodeFunc func(tx *store.Connection) (*code.SignupCode, error)
	sendCodeFunc   func(from, to string, sc *code.SignupCode) error
)

func (a *API) sendSignupCode(ctx context.Context, userID uuid.UUID, to string, create createCodeFunc, send sendCodeFunc) error {
	if a.mail == nil || a.mail.IsOffline() {
		a.log.Warn("mail not found")
		return nil
	}
	var err error
	to, err = a.ValidateEmail(to)
	if err != nil {
		return a.logError(err)
	}
	reqRole := user.ToRole(string(a.config.Signup.Invites))
	return a.conn.Transaction(func(tx *store.Connection) (err error) {
		var from string
		u := user.NewSystemUser()
		if userID != user.SystemID {
			u, err = users.GetActiveUserWithRole(tx, userID, reqRole)
			if err != nil {
				return
			}
			from = u.EmailAddress().String()
		}
		if u.Role == user.RoleUser {
			err = a.rateLimitCode(tx, u.ID)
			if err != nil {
				a.log.Warnf("rate limit exceeded for user: %s", userID)
				return err
			}
		}
		sc, err := create(tx)
		if err != nil {
			return
		}
		err = audit.LogTokenGranted(ctx, tx, sc)
		if err != nil {
			return
		}
		err = send(from, to, sc)
		if err != nil {
			return
		}
		err = codes.SignupCodeSent(tx, sc)
		if err != nil {
			return
		}
		return audit.LogCodeSent(ctx, tx, sc)
	})
}

func (a *API) rateLimitCode(tx *store.Connection, userID uuid.UUID) error {
	last, err := codes.GetLastSentSignupCode(tx, userID)
	if err != nil {
		return err
	}
	if last == nil {
		return nil
	}
	return a.config.Mail.CheckSendLimit(last.SentAt)
}
