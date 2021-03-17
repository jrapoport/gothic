package tokens

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// Bearer is a bearer token
const Bearer token.Class = "bearer"

// BearerToken holds a bearer (and refresh) token
type BearerToken struct {
	*token.AccessToken
	RefreshToken *token.RefreshToken
}

var _ token.Token = (*BearerToken)(nil)

// NewBearerToken generates a new bearer token.
func NewBearerToken(t *jwt.Token) (*BearerToken, error) {
	if t == nil {
		return nil, errors.New("invalid token")
	}
	std := t.Claims().Standard()
	userID, err := uuid.Parse(std.Subject)
	if err != nil {
		return nil, err
	}
	tok, err := t.Bearer()
	if err != nil {
		return nil, err
	}
	at := token.NewAccessToken(tok, token.InfiniteUse, t.Expiration())
	at.UserID = userID
	if std.ExpiresAt != nil {
		at.ExpiredAt = &std.ExpiresAt.Time
	}
	return &BearerToken{AccessToken: at}, nil
}

// Class returns the class for the bearer token.
func (t BearerToken) Class() token.Class {
	return Bearer
}

// GrantBearerToken grants a new bearer token
func GrantBearerToken(conn *store.Connection, c config.JWT, u *user.User) (*BearerToken, error) {
	return RefreshBearerToken(conn, c, u, "")
}

// RefreshBearerToken refreshes a bearer token
func RefreshBearerToken(conn *store.Connection, c config.JWT, u *user.User, tok string) (*BearerToken, error) {
	if u == nil {
		return nil, errors.New("invalid user")
	}
	if !u.IsActive() && !u.IsRestricted() {
		err := fmt.Errorf("inactive user: %s", u.ID)
		return nil, err
	}
	var bt *BearerToken
	err := conn.Transaction(func(tx *store.Connection) (err error) {
		t := jwt.NewUserToken(c, u)
		bt, err = NewBearerToken(t)
		if err != nil {
			return err
		}
		if tok == "" {
			bt.RefreshToken, err = GrantRefreshToken(tx, bt.UserID)
		} else {
			bt.RefreshToken, err = SwapRefreshToken(tx, bt.UserID, tok)
		}
		if err != nil {
			return err
		}
		if bt.UserID != bt.RefreshToken.UserID {
			return errors.New("mismatched user id")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bt, nil
}
