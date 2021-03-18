package signup_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/signup"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCase(t *testing.T, srv *rest.Host, web *httptest.Server) (url.Values, *rest.UserResponse) {
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	addr, err := url.Parse(web.URL)
	require.NoError(t, err)
	em := tutils.RandomEmail()
	un := utils.RandomUsername()
	data := types.Map{
		"foo":         "bar",
		"tasty":       "salad",
		key.IPAddress: addr.Host,
	}
	d, err := data.JSON()
	require.NoError(t, err)
	v := url.Values{}
	v.Set(key.Email, em)
	v.Set(key.Username, un)
	v.Set(key.Password, testPass)
	v.Set(key.Data, string(d))
	if srv.Config().MaskEmails {
		em = utils.MaskEmail(em)
	}
	ur := &rest.UserResponse{
		Role:     user.RoleUser.String(),
		Email:    em,
		Username: un,
		Data:     data,
	}
	return v, ur
}

func assertResponse(t *testing.T, h *rest.Host, test *rest.UserResponse, res string) {
	ur, _ := tsrv.UnmarshalUserResponse(t, h.Config().JWT, res)
	assert.NotNil(t, ur.Token)
	assert.NotEmpty(t, ur.Token.Access)
	ur.Token = nil
	assert.Equal(t, test, ur)
}

func TestSignupServer_Signup(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	// json success (unmasked)
	srv.Config().MaskEmails = false
	v, test := testCase(t, srv, web)
	res, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	assertResponse(t, srv, test, res)
	u, err := srv.GetUserWithEmail(v.Get(key.Email))
	assert.NoError(t, err)
	assert.Equal(t, user.RoleUser, u.Role)
	assert.Equal(t, v.Get(key.Email), u.Email)
	assert.Equal(t, v.Get(key.Username), u.Username)
	// email taken
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// json success (masked)
	srv.Config().MaskEmails = true
	v, test = testCase(t, srv, web)
	res, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	assertResponse(t, srv, test, res)
	// form encoded success
	v, test = testCase(t, srv, web)
	res, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, v, v)
	assert.NoError(t, err)
	assertResponse(t, srv, test, res)
	// query string success
	v, test = testCase(t, srv, web)
	res, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, v, nil)
	assert.NoError(t, err)
	assertResponse(t, srv, test, res)
}

func TestSignupServer_Signup_Confirm(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	srv.Config().Signup.AutoConfirm = false
	v, _ := testCase(t, srv, web)
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	em := v.Get(key.Email)
	u, err := srv.GetUserWithEmail(em)
	assert.NoError(t, err)
	assert.False(t, u.IsConfirmed())
}

func TestSignupServer_Signup_AutoConfirm(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	//
	srv.Config().Signup.AutoConfirm = true
	v, _ := testCase(t, srv, web)
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, v, nil)
	assert.NoError(t, err)
	em := v.Get(key.Email)
	u, err := srv.GetUserWithEmail(em)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
}

func TestSignupServer_Signup_Disabled(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	srv.Config().Signup.Disabled = true
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, nil)
	assert.Error(t, err)
}

func TestSignupServer_Signup_EmailDisabled(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	srv.Config().UseInternal = false
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, nil)
	assert.Error(t, err)
}

func TestSignupServer_Signup_Recaptcha(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	srv.Config().Recaptcha.Key = validate.ReCaptchaDebugKey
	v, _ := testCase(t, srv, web)
	// invalid client ip
	assert.HTTPError(t, web.Config.Handler.ServeHTTP, http.MethodPost, signup.Endpoint, v)
	// no token
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// invalid token
	v.Set(key.ReCaptcha, "invalid")
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// token
	v.Set(key.ReCaptcha, validate.ReCaptchaDebugToken)
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
}

func TestSignupServer_Signup_SignupCode(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	// no code
	srv.Config().Signup.Code = false
	v, _ := testCase(t, srv, web)
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	srv.Config().Signup.Code = true
	// missing code
	v, _ = testCase(t, srv, web)
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// bad code
	v.Set(key.Code, "bad")
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// good code
	pin, err := srv.CreateSignupCode(context.Background(), code.SingleUse)
	assert.NoError(t, err)
	v.Set(key.Code, pin)
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	// reuse code
	v, _ = testCase(t, srv, web)
	v.Set(key.Code, pin)
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
}

func TestSignupServer_Signup_Password(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		signup.RegisterServer,
	}, false)
	const passRegex = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	srv.Config().Validation.PasswordRegex = passRegex
	// good pw
	v, _ := testCase(t, srv, web)
	_, err := thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	// missing pw
	v, _ = testCase(t, srv, web)
	v.Del(key.Password)
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// bad password
	v.Set(key.Password, "nope")
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	// blank password ok
	srv.Config().Validation.PasswordRegex = ""
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
	// custom password
	srv.Config().Validation.PasswordRegex = "^[a-z]"
	v, _ = testCase(t, srv, web)
	v.Set(key.Password, "12345678")
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.Error(t, err)
	v.Set(key.Password, "password")
	_, err = thttp.DoRequest(t, web, http.MethodPost, signup.Endpoint, nil, v)
	assert.NoError(t, err)
}
