package login

import (
	"net/http"
	"testing"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"

func testUser(t *testing.T, srv *rest.Host) *user.User {
	rk := srv.Config().Recaptcha.Key
	em := tutils.RandomEmail()
	srv.Config().Recaptcha.Key = ""
	ctx := context.Background()
	ctx.SetProvider(srv.Provider())
	u, err := srv.Signup(ctx, em, "", testPass, nil)
	srv.Config().Recaptcha.Key = rk
	require.NoError(t, err)
	require.NotNil(t, u)
	return u
}

func TestLoginServer_Logout(t *testing.T) {
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		RegisterServer,
	}, false)
	j := srv.Config().JWT
	// not authorized
	_, err := thttp.DoAuthRequest(t, web, http.MethodGet, Logout, "", nil, nil)
	assert.Error(t, err)
	// invalid id
	bad := thttp.BadToken(t, j)
	_, err = thttp.DoAuthRequest(t, web, http.MethodGet, Logout, bad, nil, nil)
	assert.Error(t, err)
	// user not found
	tok := thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodGet, Logout, tok, nil, nil)
	// login then logout
	u := testUser(t, srv)
	req := struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}{u.Email, testPass}
	res, err := thttp.DoRequest(t, web, http.MethodPost, Endpoint, nil, req)
	assert.NoError(t, err)
	ur, _ := tsrv.MarshalUserResponse(t, srv.Config().JWT, res)
	assert.EqualValues(t, tokens.Bearer, ur.Token.Type)
	_, err = thttp.DoAuthRequest(t, web, http.MethodGet, Logout, ur.Token.Access, nil, nil)
	assert.NoError(t, err)
	_, err = srv.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
}
