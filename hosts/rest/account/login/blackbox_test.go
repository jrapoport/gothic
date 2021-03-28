package login_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/login"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
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

type testRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func TestLoginServer_Login(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		login.RegisterServer,
	}, false)
	// invalid req
	_, err := thttp.DoRequest(t, web, http.MethodPost, login.Login, nil, []byte("\n"))
	assert.Error(t, err)
	// empty email
	req := new(testRequest)
	_, err = thttp.DoRequest(t, web, http.MethodPost, login.Login, nil, req)
	assert.Error(t, err)
	// bad email
	req = &testRequest{
		Email: "bad",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, login.Login, nil, req)
	assert.Error(t, err)
	// not found
	req = &testRequest{
		Email: "bad@example.com",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, login.Login, nil, req)
	assert.Error(t, err)
	u := testUser(t, srv)
	// bad password
	req = &testRequest{
		Email:    u.Email,
		Password: "",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, login.Login, nil, req)
	assert.Error(t, err)
	// login
	_, err = srv.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	req = &testRequest{
		Email:    u.Email,
		Password: testPass,
	}
	res, err := thttp.DoRequest(t, web, http.MethodPost, login.Login, nil, req)
	assert.NoError(t, err)
	ur, claims := tsrv.UnmarshalUserResponse(t, srv.Config().JWT, res)
	assert.EqualValues(t, tokens.Bearer, ur.Token.Type)
	assert.Equal(t, u.ID.String(), claims.Subject())
	assert.Equal(t, ur.Email, utils.MaskEmail(u.Email))
	au, err := srv.GetAuthenticatedUser(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
}

func TestLoginServer_Login_Recaptcha(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		login.RegisterServer,
	}, false)
	srv.Config().Recaptcha.Key = validate.ReCaptchaDebugKey
	v := url.Values{}
	v.Set(key.Email, tutils.RandomEmail())
	v.Set(key.Password, testPass)
	// invalid client ip
	assert.HTTPError(t, web.Config.Handler.ServeHTTP, http.MethodPost, login.Login, v)
	// no token
	_, err := thttp.DoRequest(t, web, http.MethodPost, login.Login, v, nil)
	assert.Error(t, err)
	// invalid token
	v.Set(key.ReCaptcha, "invalid")
	_, err = thttp.DoRequest(t, web, http.MethodPost, login.Login, v, nil)
	assert.Error(t, err)
	// token
	u := testUser(t, srv)
	v.Set(key.Email, u.Email)
	v.Set(key.ReCaptcha, validate.ReCaptchaDebugToken)
	_, err = thttp.DoRequest(t, web, http.MethodPost, login.Login, v, nil)
	assert.NoError(t, err)
}
