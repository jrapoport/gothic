package core

import (
	"net/url"
	"testing"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func authToken(t *testing.T, a *API, p provider.Name) string {
	authURL, err := a.GetAuthorizationURL(context.Background(), p)
	require.NoError(t, err)
	au, err := url.Parse(authURL)
	require.NoError(t, err)
	tok := au.Query().Get(key.State)
	return tok
}

func TestAPI_GetAuthorizationURL(t *testing.T) {
	ctx := context.Background()
	a := apiWithTempDB(t)
	_, mock := tconf.MockedProvider(t, a.config, "")
	a.ext.UseProviders(mock)
	// no context
	_, err := a.GetAuthorizationURL(ctx, provider.Unknown)
	assert.Error(t, err)
	// no request context
	_, err = a.GetAuthorizationURL(context.Background(), provider.Unknown)
	assert.Error(t, err)
	// bad provider
	ctx.SetProvider("bad")
	_, err = a.GetAuthorizationURL(ctx, "bad")
	assert.Error(t, err)
	// internal provider
	ctx.SetProvider(a.Provider())
	_, err = a.GetAuthorizationURL(ctx, a.Provider())
	assert.Error(t, err)
	// disabled provider
	ctx.SetProvider(provider.BitBucket)
	_, err = a.GetAuthorizationURL(ctx, provider.BitBucket)
	assert.Error(t, err)
	// valid external provider
	ctx.SetProvider(mock.PName())
	authURL, err := a.GetAuthorizationURL(ctx, mock.PName())
	assert.NoError(t, err)
	_, err = url.Parse(authURL)
	assert.NoError(t, err)
}

func TestAPI_AuthorizeUser(t *testing.T) {
	a := apiWithTempDB(t)
	_, mock := tconf.MockedProvider(t, a.config, "")
	a.ext.UseProviders(mock)
	// no token
	_, err := a.AuthorizeUser(nil, "", nil)
	assert.Error(t, err)
	// bad token
	_, err = a.AuthorizeUser(nil, "bad", nil)
	assert.Error(t, err)
	// bad provider
	at, err := tokens.GrantAuthToken(a.conn, "bad", 0)
	require.NoError(t, err)
	_, err = a.AuthorizeUser(nil, at.Token, nil)
	assert.Error(t, err)
	// provider not found
	at, err = tokens.GrantAuthToken(a.conn, provider.Google, 0)
	require.NoError(t, err)
	_, err = a.AuthorizeUser(nil, at.Token, nil)
	assert.Error(t, err)
	// invalid session
	at, err = tokens.GrantAuthToken(a.conn, mock.PName(), 0)
	require.NoError(t, err)
	_, err = a.AuthorizeUser(nil, at.Token, nil)
	assert.Error(t, err)
	// create
	tok := authToken(t, a, mock.PName())
	u, err := a.AuthorizeUser(context.Background(), tok, nil)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsConfirmed())
	assert.True(t, u.IsActive())
	assert.Equal(t, mock.PName(), u.Provider)
	assert.Equal(t, mock.Username, u.Username)
	assert.Equal(t, mock.Email, u.Email)
	// update
	username := u.Username
	tok = authToken(t, a, mock.PName())
	u, err = a.AuthorizeUser(context.Background(), tok, nil)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.NoError(t, err)
	assert.Equal(t, mock.PName(), u.Provider)
	assert.Equal(t, username, u.Username)
	assert.Equal(t, mock.Email, u.Email)
}
