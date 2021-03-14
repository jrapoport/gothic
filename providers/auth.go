package providers

import (
	"errors"
	"time"

	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/markbates/goth"
)

// AuthURL holds a authorization url and token.
type AuthURL struct {
	URL   string
	Token *token.AuthToken
}

// GrantAuthURL returns the auth url for a named provider.
func GrantAuthURL(conn *store.Connection, p provider.Name, exp time.Duration) (*AuthURL, error) {
	gp, err := GetProvider(p)
	if err != nil {
		return nil, err
	}
	var au = &AuthURL{}
	err = conn.Transaction(func(tx *store.Connection) error {
		au.Token, err = tokens.GrantAuthToken(tx, p, exp)
		if err != nil {
			return err
		}
		s, err := gp.BeginAuth(au.Token.String())
		if err != nil {
			return err
		}
		err = tx.Model(au.Token).Update(key.Data, types.Map{
			key.Session: s.Marshal(),
		}).Error
		if err != nil {
			return err
		}
		au.URL, err = s.GetAuthURL()
		return err
	})
	if err != nil {
		return nil, err
	}
	return au, nil
}

// AuthorizeUser checks the token and turns the oauth authorized user.
func AuthorizeUser(conn *store.Connection, tok string, data types.Map) (*goth.User, error) {
	var u goth.User
	err := conn.Transaction(func(tx *store.Connection) error {
		t, err := tokens.GetAuthToken(tx, tok)
		if err != nil {
			return err
		}
		gp, err := GetProvider(t.Provider)
		if err != nil {
			return err
		}
		sd, ok := t.Data[key.Session].(string)
		if !ok {
			return errors.New("invalid session")
		}
		s, err := gp.UnmarshalSession(sd)
		if err != nil {
			return err
		}
		_, err = s.Authorize(gp, data)
		if err != nil {
			return err
		}
		u, err = gp.FetchUser(s)
		if err != nil {
			return err
		}
		return tokens.UseToken(tx, t)
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}
