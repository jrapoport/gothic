package core

import (
	"errors"
	"regexp"
	"time"

	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/providers"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/utils"
)

// GetAuthorizationURL get the auth url for a configured provider.
func (a *API) GetAuthorizationURL(ctx context.Context, p provider.Name) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx.SetProvider(p)
	var au *providers.AuthURL
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		au, err = providers.GrantAuthURL(a.conn, p, 60*time.Minute)
		if err != nil {
			return err
		}
		return audit.LogTokenGranted(ctx, tx, au.Token)
	})
	if err != nil {
		return "", a.logError(err)
	}
	return au.URL, nil
}

// AuthorizeUser authorizes a user with a configured provider.
func (a *API) AuthorizeUser(ctx context.Context, tok string, data types.Map) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var u *user.User
	err := a.conn.Transaction(func(tx *store.Connection) error {
		au, err := providers.AuthorizeUser(tx, tok, data)
		if err != nil {
			return err
		}
		data = types.Map{
			key.Name:        au.Name,
			key.FirstName:   au.FirstName,
			key.LastName:    au.LastName,
			key.Nickname:    au.NickName,
			key.Description: au.Description,
			key.AvatarURL:   au.AvatarURL,
		}
		for k, v := range data {
			if v == "" {
				delete(data, k)
			}
		}
		p := provider.Name(au.Provider)
		ctx.SetProvider(p)
		u, err = users.HasLinkedUser(tx, p, au.UserID)
		if err != nil {
			return err
		}
		if u == nil {
			raw := types.DataFromMap(au.RawData)
			u, err = a.externalCreate(ctx, tx, p, au.UserID, au.Email, data, raw)
		} else {
			if !u.IsActive() {
				return errors.New("user account is not active")
			}
			// TODO: We should login the user here w/ the login API
			err = a.externalUpdate(tx, u, au.Email, data)
			if err != nil {
				return err
			}
			err = audit.LogUserUpdated(ctx, tx, u.ID)
		}
		return err
	})
	if err != nil {
		return nil, a.logError(err)
	}
	if u.IsConfirmed() {
		return u, nil
	}
	err = a.autoConfirm(ctx, u)
	if err != nil {
		return nil, a.logError(err)
	}
	return u, nil
}

func (a *API) externalCreate(ctx context.Context, conn *store.Connection,
	p provider.Name, accountID, email string, data, raw types.Map) (*user.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var u *user.User
	err := conn.Transaction(func(tx *store.Connection) (err error) {
		pw := utils.SecureToken()
		username := getUsername(data)
		a.log.Debugf("external provider create: %s %s %s %v (%v)",
			email, username, pw, data, ctx)
		u, err = a.userSignup(ctx, tx, p, email, username, pw, data)
		if err != nil {
			return err
		}
		return a.linkAccount(ctx, tx, u, accountID, raw)
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (a *API) linkAccount(ctx context.Context, conn *store.Connection,
	u *user.User, accountID string, raw types.Map) error {
	if raw == nil {
		raw = types.Map{}
	}
	ip := ctx.GetIPAddress()
	raw[key.IPAddress] = ip
	la := &account.Linked{
		Type:      account.Auth,
		Provider:  u.Provider,
		AccountID: accountID,
		Email:     u.Email,
		Data:      raw,
	}
	return conn.Transaction(func(tx *store.Connection) error {
		err := users.LinkAccount(tx, u, la)
		if err != nil {
			return err
		}
		return audit.LogLinked(ctx, tx, u.ID, la)
	})
}

func (a *API) externalUpdate(conn *store.Connection, u *user.User, email string, data types.Map) error {
	return conn.Transaction(func(tx *store.Connection) (err error) {
		if u.Email != email {
			err = users.ChangeEmail(tx, u, email)
		}
		if err != nil {
			return err
		}
		return users.Update(tx, u, nil, data)
	})
}

func getUsername(data types.Map) string {
	var name string
	if n, ok := data[key.Name].(string); ok && n != "" {
		name = n
	} else if n, ok = data[key.Nickname].(string); ok && n != "" {
		name = n
	}
	// only unicode letters & numbers
	rx := regexp.MustCompile(`[^\p{L}0-9]+`)
	return rx.ReplaceAllString(name, "")
}
