package providers

import (
	"net/url"
	"testing"
	"time"

	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getToken(t *testing.T, authURL string) string {
	au, err := url.Parse(authURL)
	require.NoError(t, err)
	return au.Query().Get(key.State)
}

func TestGrantAuthURL(t *testing.T) {
	const badProvider = "bad"
	providers := NewProviders()
	conn, c := tconn.TempConn(t)
	c.Authorization = tconf.ProvidersConfig(t).Authorization
	err := providers.LoadProviders(c)
	require.NoError(t, err)
	for p := range provider.External {
		// these providers will attempt a live connection
		if p == provider.Twitter || p == provider.Xero || p == provider.Tumblr {
			// assert.Contains(t, err.Error(), "401 Unauthorized")
			continue
		}
		auth, err := providers.GrantAuthURL(conn, p, 60*time.Minute)
		assert.NoError(t, err)
		_, err = url.Parse(auth.URL)
		assert.NoError(t, err)
	}
	_, err = providers.GrantAuthURL(conn, provider.Unknown, 0)
	assert.Error(t, err)
	_, err = providers.GrantAuthURL(conn, c.Provider(), 0)
	assert.Error(t, err)
	_, err = providers.GrantAuthURL(conn, badProvider, 0)
	assert.Error(t, err)
	// confirm different tokens are returned
	_, mock := tconf.MockedProvider(t, c, "")
	providers.UseProviders(mock)
	auth1, err := providers.GrantAuthURL(conn, mock.PName(), 60*time.Minute)
	assert.NoError(t, err)
	auth2, err := providers.GrantAuthURL(conn, mock.PName(), 60*time.Minute)
	assert.NoError(t, err)
	tok1 := getToken(t, auth1.URL)
	assert.Equal(t, auth1.Token.String(), tok1)
	tok2 := getToken(t, auth2.URL)
	assert.Equal(t, auth2.Token.String(), tok2)
	assert.NotEqual(t, tok1, tok2)
}

func TestAuthorizeUser(t *testing.T) {
	providers := NewProviders()
	conn, c := tconn.TempConn(t)
	_, mock := tconf.MockedProvider(t, c, "")
	providers.UseProviders(mock)
	authURL, err := providers.GrantAuthURL(conn, mock.PName(), 0)
	require.NoError(t, err)
	tok := getToken(t, authURL.URL)
	_, err = providers.AuthorizeUser(conn, tok, nil)
	require.NoError(t, err)
	// cannot reuse token
	_, err = providers.AuthorizeUser(conn, tok, nil)
	require.Error(t, err)
	// empty token
	_, err = providers.AuthorizeUser(conn, "", nil)
	require.Error(t, err)
	// bad token
	_, err = providers.AuthorizeUser(conn, utils.SecureToken(), nil)
	require.Error(t, err)
	// bad provider
	at, err := tokens.GrantAuthToken(conn, "bad", 0)
	require.NoError(t, err)
	_, err = providers.AuthorizeUser(conn, at.String(), nil)
	assert.Error(t, err)
	// provider not found
	at, err = tokens.GrantAuthToken(conn, provider.Google, 0)
	require.NoError(t, err)
	_, err = providers.AuthorizeUser(conn, at.String(), nil)
	assert.Error(t, err)
	// invalid session
	at, err = tokens.GrantAuthToken(conn, mock.PName(), 0)
	require.NoError(t, err)
	_, err = providers.AuthorizeUser(conn, at.String(), nil)
	assert.Error(t, err)
}
