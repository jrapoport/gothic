package auth

import (
	"errors"
	"time"

	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/markbates/goth"
)

// URL holds a authorization url and token.
type URL struct {
	URL   string
	Token *token.AuthToken
}

// GrantAuthURL returns the auth url for a named provider.
func (pv *Providers) GrantAuthURL(conn *store.Connection, name provider.Name, exp time.Duration) (*URL, error) {
	p, err := pv.GetProvider(name)
	if err != nil {
		return nil, err
	}
	var au = &URL{}
	err = conn.Transaction(func(tx *store.Connection) error {
		au.Token, err = tokens.GrantAuthToken(tx, name, exp)
		if err != nil {
			return err
		}
		s, err := p.BeginAuth(au.Token.String())
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
func (pv *Providers) AuthorizeUser(conn *store.Connection, tok string, data types.Map) (*goth.User, error) {
	var u goth.User
	err := conn.Transaction(func(tx *store.Connection) error {
		t, err := tokens.GetAuthToken(tx, tok)
		if err != nil {
			return err
		}
		p, err := pv.GetProvider(t.Provider)
		if err != nil {
			return err
		}
		sd, ok := t.Data[key.Session].(string)
		if !ok {
			return errors.New("invalid session")
		}
		s, err := p.UnmarshalSession(sd)
		if err != nil {
			return err
		}
		_, err = s.Authorize(p, data)
		if err != nil {
			return err
		}
		u, err = p.FetchUser(s)
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
