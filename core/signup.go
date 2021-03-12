package core

import (
	"errors"
	"time"

	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/providers"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/utils"
)

// Signup signs-up a new user
func (a *API) Signup(ctx context.Context, email, username, pw string, data types.Map) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	p := a.Provider()
	ctx.SetProvider(p)
	err := a.signupEnabled(p)
	if err != nil {
		return nil, err
	}
	email, err = a.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	err = validate.Password(a.config, pw)
	if err != nil {
		return nil, a.logError(err)
	}
	a.log.Debugf("signup: %s %s %s %v (%v)",
		email, username, pw, data, ctx)
	var u *user.User
	err = a.conn.Transaction(func(tx *store.Connection) error {
		username, err = a.validateUsername(tx, username)
		if err != nil {
			return err
		}
		u, err = a.userSignup(ctx, tx, p, email, username, pw, data)
		return err
	})
	if err != nil {
		return nil, a.logError(err)
	}
	if a.config.Signup.AutoConfirm {
		err = a.autoConfirm(ctx, u)
		if err != nil {
			return nil, a.logError(err)
		}
	}
	if u.IsConfirmed() {
		return u, nil
	}
	err = a.SendConfirmUser(ctx, u.ID)
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

func (a *API) validateUsername(conn *store.Connection, username string) (string, error) {
	var err error
	if username == "" && a.config.Signup.Default.Username {
		username, err = users.RandomUsername(conn, true)
		if err != nil {
			return "", a.logError(err)
		}
	}
	if username == "" {
		if a.config.Signup.Username {
			err = errors.New("username required")
		}
		return username, err
	}
	return username, validate.Username(a.config, username)
}

func (a *API) userSignup(ctx context.Context, conn *store.Connection,
	p provider.Name, email, username, pw string, data types.Map) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := a.signupEnabled(p)
	if err != nil {
		return nil, err
	}
	code := ctx.GetCode()
	if a.config.Signup.Code && code == "" {
		err = errors.New("signup code required")
		return nil, err
	}
	ip := ctx.GetIPAddress()
	recaptcha := ctx.GetReCaptcha()
	// if recaptcha is disabled this is a no-op
	err = validate.ReCaptcha(a.config, ip, recaptcha)
	if err != nil {
		return nil, err
	}
	var u *user.User
	err = conn.Transaction(func(tx *store.Connection) error {
		username, err = a.useDefaultUsername(tx, username)
		if err != nil {
			return err
		}
		data = a.useDefaultColor(data)
		meta := types.Map{key.IPAddress: ip}
		u, err = users.CreateUser(tx, p, email, username, pw, data, meta)
		if err != nil && utils.IsDebugPIN(code) {
			u, err = users.GetUserWithEmail(tx, email)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		err = a.useSignupCode(tx, u, code)
		if err != nil {
			return err
		}
		return audit.LogSignup(ctx, tx, u.ID, u.Role)
	})
	if err != nil {
		return nil, err
	}
	a.dispatchEvent(events.Signup, types.Map{
		key.Provider:  u.Provider,
		key.IPAddress: ip,
		key.UserID:    u.ID,
		key.Email:     u.Email,
		key.Timestamp: time.Now().UTC(),
	})
	return u, nil
}

// autoConfirm automatically confirms a user
func (a *API) autoConfirm(ctx context.Context, u *user.User) error {
	if u == nil || u.IsLocked() {
		return errors.New("invalid user")
	}
	if u.IsConfirmed() {
		return nil
	}
	err := a.conn.Transaction(func(tx *store.Connection) error {
		err := users.ConfirmUser(tx, u, time.Now().UTC())
		if err != nil {
			return err
		}
		return audit.LogConfirmed(ctx, tx, u.ID)
	})
	if err != nil {
		return a.logError(err)
	}
	const typeAuto = "auto"
	a.dispatchEvent(events.Confirmed, types.Map{
		key.Type:      typeAuto,
		key.Provider:  u.Provider,
		key.Email:     u.Email,
		key.Timestamp: time.Now().UTC(),
	})
	return nil
}

func (a *API) signupEnabled(p provider.Name) error {
	if a.config.Signup.Disabled {
		err := errors.New("signup disabled")
		return err
	}
	return providers.IsEnabled(p)
}

func (a *API) useDefaultUsername(tx *store.Connection, username string) (string, error) {
	if !a.config.Signup.Default.Username || username != "" {
		return username, nil
	}
	var err error
	username, err = users.RandomUsername(tx, true)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (a *API) useDefaultColor(data types.Map) types.Map {
	if !a.config.Signup.Default.Color {
		return data
	}
	if data == nil {
		data = types.Map{}
	}
	if _, ok := data[key.Color]; !ok {
		data[key.Color] = utils.RandomColor()
	}
	return data
}

func (a *API) useSignupCode(tx *store.Connection, u *user.User, code string) error {
	if !a.config.Signup.Code {
		return nil
	}
	if code == "" {
		return errors.New("signup code required")
	}
	sc, err := codes.GetUsableCode(tx, code)
	if err != nil {
		return err
	}
	return sc.UseCode(tx, u)
}
