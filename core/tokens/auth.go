package tokens

import (
	"time"

	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/store"
)

// GrantAuthToken gets or creates an auth token for the provider.
func GrantAuthToken(conn *store.Connection, p provider.Name, exp time.Duration) (*token.AuthToken, error) {
	var t token.Token
	err := conn.Transaction(func(tx *store.Connection) error {
		t = token.NewAuthToken(p, exp)
		return conn.Create(t).Error
	})
	if err != nil {
		return nil, err
	}
	return t.(*token.AuthToken), nil
}

// GetAuthToken returns the auth token for the token string if found.
func GetAuthToken(conn *store.Connection, tok string) (*token.AuthToken, error) {
	var ct token.AuthToken
	err := conn.First(&ct, "token = ?", tok).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}
