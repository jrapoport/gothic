package auth_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/auth"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthServer_RefreshBearerToken(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		auth.RegisterServer,
	}, false)
	srv.Config().Signup.AutoConfirm = true
	u, _ := tcore.TestUser(t, srv.API, "", false)
	bt, err := srv.GrantBearerToken(context.Background(), u)
	require.NoError(t, err)
	// bad request
	_, err = thttp.DoRequest(t, web, http.MethodPost, auth.Endpoint, nil, []byte("\n"))
	assert.Error(t, err)
	// no token
	v := url.Values{}
	res, err := thttp.DoRequest(t, web, http.MethodPost, auth.Endpoint, v, nil)
	assert.Error(t, err)
	// token
	v[key.Token] = []string{bt.RefreshToken.Token}
	res, err = thttp.DoRequest(t, web, http.MethodPost, auth.Endpoint, v, nil)
	assert.NoError(t, err)
	tr, claims := tsrv.UnmarshalTokenResponse(t, srv.Config().JWT, res)
	assert.EqualValues(t, tokens.Bearer, tr.Type)
	uid, err := uuid.Parse(claims.Subject)
	assert.NoError(t, err)
	u2, err := srv.GetUser(uid)
	assert.NoError(t, err)
	assert.Equal(t, u2.ID.String(), claims.Subject)
	au, err := srv.GetAuthenticatedUser(u2.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
	res, err = thttp.DoRequest(t, web, http.MethodPost, auth.Endpoint, v, nil)
	assert.NoError(t, err)
	assert.Empty(t, res)
}
