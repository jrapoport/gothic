package email_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/user/email"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	change   = email.Email + email.Change
	unmask   = email.Email + rest.Root
	confirm  = email.Email + email.Confirm
)

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		email.RegisterServer,
	}, true)
	c := srv.Config()
	c.Validation.UsernameRegex = ""
	c.Signup.Username = false
	c.Signup.AutoConfirm = true
	c.Mail.SendLimit = 0
	t.Cleanup(func() {
		web.Close()
	})
	err := srv.API.LoadConfig(c)
	require.NoError(t, err)
	return srv, web, smtp
}

func TestEmailServer_ConfirmChangeEmail(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	var newEmail = tutils.RandomEmail()
	// invalid req
	_, err := thttp.DoRequest(t, web, http.MethodPost, confirm, nil, []byte("\n"))
	assert.Error(t, err)
	// empty token
	req := new(email.Request)
	_, err = thttp.DoRequest(t, web, http.MethodPost, confirm, nil, req)
	assert.Error(t, err)
	// bad token
	req = &email.Request{
		Token: "bad",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, confirm, nil, req)
	assert.Error(t, err)
	// first get the change token
	u, bt := tcore.TestUser(t, srv.API, "", false)
	var tok string
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(template.ChangeEmailAction, email)
	})
	req = &email.Request{
		Email:    newEmail,
		Password: testPass,
	}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bt, nil, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
	// now use the token to change the email
	req = &email.Request{
		Token: tok,
	}
	res, err := thttp.DoAuthRequest(t, web, http.MethodPost, confirm, bt, nil, req)
	assert.NoError(t, err)
	_, claims := tsrv.UnmarshalTokenResponse(t, srv.Config().JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject)
	updated, err := srv.GetUser(u.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, newEmail, u.Email)
	assert.Equal(t, newEmail, updated.Email)
}

func TestEmailServer_SendChangeEmail(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	u, bt := tcore.TestUser(t, srv.API, "", false)
	j := srv.Config().JWT
	// not authorized
	_, err := thttp.DoRequest(t, web, http.MethodPost, change, nil, nil)
	assert.Error(t, err)
	// no user id
	bad := thttp.BadToken(t, j)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bad, nil, nil)
	assert.Error(t, err)
	// user not found
	bad = thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bad, nil, nil)
	assert.Error(t, err)
	// invalid req
	bad = thttp.UserToken(t, j, true, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bad, nil, []byte("\n"))
	assert.Error(t, err)
	// user not confirmed
	bad = thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bad, nil, nil)
	assert.Error(t, err)
	// empty email
	req := new(email.Request)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bt, nil, req)
	assert.Error(t, err)
	// bad email
	req = &email.Request{
		Email: "bad",
	}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bt, nil, req)
	assert.Error(t, err)
	// bad password
	req = &email.Request{
		Email:    u.Email,
		Password: "bad",
	}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bt, nil, req)
	assert.Error(t, err)
	// success
	var tok string
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(template.ChangeEmailAction, email)
	})
	req = &email.Request{
		Email:    u.Email,
		Password: testPass,
	}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, change, bt, nil, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
}

func TestEmailServer_SendChangeEmail_RateLimit(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	srv.Config().Mail.SendLimit = 5 * time.Minute
	u, bt := tcore.TestUser(t, srv.API, "", false)
	var sent string
	smtp.AddHook(t, func(email string) {
		sent = email
	})
	for i := 0; i < 2; i++ {
		sent = ""
		req := &email.Request{
			Email:    u.Email,
			Password: testPass,
		}
		_, err := thttp.DoAuthRequest(t, web, http.MethodPost, change, bt, nil, req)
		if i == 0 {
			assert.NoError(t, err)
			assert.Eventually(t, func() bool {
				return sent != ""
			}, 1*time.Second, 10*time.Millisecond)
		} else {
			msg := thttp.FmtError(http.StatusTooEarly).Error()
			assert.EqualError(t, err, msg)
			assert.Never(t, func() bool {
				return sent != ""
			}, 1*time.Second, 10*time.Millisecond)
		}
	}
}
