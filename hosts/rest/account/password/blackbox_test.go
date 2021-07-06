package password_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/password"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	resetEndpoint   = password.Endpoint + password.Reset
	confirmEndpoint = password.Endpoint + password.Confirm
)

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		password.RegisterServer,
	}, true)
	c := srv.Config()
	c.Validation.UsernameRegex = ""
	c.Signup.Username = false
	c.Signup.AutoConfirm = false
	c.Mail.SendLimit = 0
	t.Cleanup(func() {
		web.Close()
	})
	err := srv.API.LoadConfig(c)
	require.NoError(t, err)
	return srv, web, smtp
}

func testUser(t *testing.T, srv *rest.Host) *user.User {
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	em := tutils.RandomEmail()
	ctx := context.Background()
	ctx.SetProvider(srv.Provider())
	u, err := srv.Signup(ctx, em, "", testPass, nil)
	require.NoError(t, err)
	require.NotNil(t, u)
	return u
}

func TestPasswordServer_SendResetPassword(t *testing.T) {
	srv, web, smtp := testServer(t)
	u := testUser(t, srv)
	// invalid req
	_, err := thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, []byte("\n"))
	assert.Error(t, err)
	// empty email
	req := new(password.Request)
	_, err = thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, req)
	assert.Error(t, err)
	// bad email
	req = &password.Request{
		Email: "bad",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, req)
	assert.Error(t, err)
	// not found
	req = &password.Request{
		Email: "i-dont-exist@example.com",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, req)
	assert.NoError(t, err)
	var tok string
	act := template.ResetPasswordAction
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(act, email)
	})
	req = &password.Request{
		Email: u.Email,
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 100*time.Millisecond)
}

func TestPasswordServer_SendResetPassword_RateLimit(t *testing.T) {
	srv, web, smtp := testServer(t)
	srv.Config().Mail.SendLimit = 5 * time.Minute
	srv.Config().Signup.AutoConfirm = true
	u := testUser(t, srv)
	var sent string
	smtp.AddHook(t, func(email string) {
		sent = email
	})
	for i := 0; i < 2; i++ {
		sent = ""
		req := &password.Request{
			Email: u.Email,
		}
		_, err := thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, req)
		if i == 0 {
			assert.NoError(t, err)
			assert.Eventually(t, func() bool {
				return sent != ""
			}, 1*time.Second, 100*time.Millisecond)
		} else {
			msg := thttp.FmtError(http.StatusTooEarly).Error()
			assert.EqualError(t, err, msg)
			assert.Never(t, func() bool {
				return sent != ""
			}, 1*time.Second, 100*time.Millisecond)
		}
	}
}

func TestPasswordServer_ConfirmPasswordChange(t *testing.T) {
	srv, web, smtp := testServer(t)
	const newPass = "sxjAm7QJ4?3dH!aN8T3F5P!oNnpXbaRy#gtx#8jG"
	// invalid req
	_, err := thttp.DoRequest(t, web, http.MethodPost, confirmEndpoint, nil, []byte("\n"))
	assert.Error(t, err)
	// empty token
	req := new(password.Request)
	_, err = thttp.DoRequest(t, web, http.MethodPost, confirmEndpoint, nil, req)
	assert.Error(t, err)
	// bad token
	req = &password.Request{
		Token: "bad",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, confirmEndpoint, nil, req)
	assert.Error(t, err)
	// first get the change token
	u := testUser(t, srv)
	assert.False(t, u.IsConfirmed())
	var tok string
	act := template.ResetPasswordAction
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(act, email)
	})
	req = &password.Request{
		Email: u.Email,
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, resetEndpoint, nil, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 100*time.Millisecond)
	// now use the token to change the password
	req = &password.Request{
		Password: newPass,
		Token:    tok,
	}
	res, err := thttp.DoRequest(t, web, http.MethodPost, confirmEndpoint, nil, req)
	assert.NoError(t, err)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	_, claims := tsrv.MarshalTokenResponse(t, srv.Config().JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	err = u.Authenticate(newPass)
	assert.NoError(t, err)
}
