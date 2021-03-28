package account_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account"
	"github.com/jrapoport/gothic/hosts/rest/account/confirm"
	"github.com/jrapoport/gothic/hosts/rest/account/login"
	"github.com/jrapoport/gothic/hosts/rest/account/password"
	"github.com/jrapoport/gothic/hosts/rest/account/signup"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		account.RegisterServer,
	}, true)
	c := srv.Config()
	c.Mail.SendLimit = 0
	c.Security.MaskEmails = false
	c.Signup.AutoConfirm = false
	c.Signup.Default.Color = true
	c.Signup.Default.Username = true
	c.Signup.Username = false
	c.Validation.UsernameRegex = ""
	t.Cleanup(func() {
		web.Close()
	})
	err := srv.API.LoadConfig(c)
	require.NoError(t, err)
	return srv, web, smtp
}

func TestAccountServer(t *testing.T) {
	t.Parallel()
	const (
		testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
		newPass  = "gj8#xtg#yrabxpnno!p5f3t8na!hd3?4jq7majxs"
	)
	srv, web, smtp := testServer(t)
	c := srv.Config()
	em := tutils.RandomEmail()
	// signup
	sr := &signup.Request{
		Email:    em,
		Password: testPass,
	}
	route := account.Account + signup.Signup
	res, err := thttp.DoRequest(t, web, http.MethodPost, route, nil, sr)
	assert.NoError(t, err)
	ur := rest.UserResponse{}
	err = json.Unmarshal([]byte(res), &ur)
	assert.NoError(t, err)
	assert.Equal(t, em, ur.Email)
	// check user
	u, err := srv.GetUserWithEmail(em)
	assert.NoError(t, err)
	assert.Equal(t, em, u.Email)
	// user is not confirmed
	assert.False(t, u.IsConfirmed())
	// resend confirmation email
	act := template.ConfirmUserAction
	var confirmToken string
	smtp.AddHook(t, func(email string) {
		confirmToken = tconf.GetEmailToken(act, email)
	})
	cr := &confirm.Request{
		Email: em,
	}
	route = account.Account + confirm.Confirm + confirm.Send
	_, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, cr)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return confirmToken != ""
	}, 1*time.Second, 10*time.Millisecond)
	// confirm user email
	route = account.Account + confirm.Confirm
	cr = &confirm.Request{
		Token: confirmToken,
	}
	res, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, cr)
	assert.NoError(t, err)
	require.NotEmpty(t, res)
	// we were logged in
	_, claims := tsrv.UnmarshalTokenResponse(t, c.JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject())
	// check user again
	u, err = srv.GetUserWithEmail(em)
	assert.NoError(t, err)
	assert.Equal(t, em, u.Email)
	// user is now confirmed
	assert.True(t, u.IsConfirmed())
	// request a password reset
	var passToken string
	act = template.ResetPasswordAction
	smtp.AddHook(t, func(email string) {
		passToken = tconf.GetEmailToken(act, email)
	})
	route = account.Account + password.Password + password.Reset
	pr := &password.Request{
		Email: em,
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, pr)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return passToken != ""
	}, 1*time.Second, 10*time.Millisecond)
	// reset the password
	route = account.Account + password.Password
	pr = &password.Request{
		Token:    passToken,
		Password: newPass,
	}
	res, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, pr)
	assert.NoError(t, err)
	// check that we were logged again
	_, claims = tsrv.UnmarshalTokenResponse(t, c.JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject())
	// login with the >new< password
	route = account.Account + login.Login
	lr := login.Request{
		Email:    em,
		Password: newPass,
	}
	res, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, lr)
	assert.NoError(t, err)
	// check that we were logged correctly
	_, claims = tsrv.UnmarshalUserResponse(t, c.JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject())
	// verify the old password fails
	lr.Password = testPass
	_, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, lr)
	assert.Error(t, err)
}

func TestAccountServer_RateLimit(t *testing.T) {
	t.Parallel()
	srv, web, _ := testServer(t)
	c := srv.Config()
	c.Signup.AutoConfirm = true
	var err error
	for i := 0; i < 200; i++ {
		req := &signup.Request{
			Email:    tutils.RandomEmail(),
			Password: tutils.RandomEmail(),
		}
		route := account.Account + signup.Signup
		_, err = thttp.DoRequest(t, web, http.MethodPost, route, nil, req)
		if err != nil {
			break
		}
	}
	assert.Eventually(t, func() bool {
		return err != nil
	}, 1*time.Second, 10*time.Millisecond)
	assert.EqualError(t, err,
		thttp.FmtError(http.StatusTooManyRequests).Error())
}
