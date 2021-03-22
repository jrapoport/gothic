package core

import (
	"regexp"
	"time"

	"github.com/jrapoport/gothic/core/accounts"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/auth"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// GetAuthorizationURL get the auth url for a configured provider.
func (a *API) GetAuthorizationURL(ctx context.Context, p provider.Name) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx.SetProvider(p)
	var au *auth.AuthURL
	err := a.conn.Transaction(func(tx *store.Connection) (err error) {
		au, err = a.ext.GrantAuthURL(a.conn, p, 60*time.Minute)
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
		au, err := a.ext.AuthorizeUser(tx, tok, data)
		if err != nil {
			return err
		}
		p := provider.Name(au.Provider)
		raw := types.DataFromMap(au.RawData)
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
		ctx.SetProvider(p)
		has, err := accounts.HasAccount(tx, p, au.UserID)
		if err != nil {
			return err
		}
		if has {
			u, err = a.externalLogin(ctx, tx, p, au.UserID, au.Email, data, raw)
		} else { // create a new user
			u, err = a.externalSignup(ctx, tx, p, au.UserID, au.Email, data, raw)
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
