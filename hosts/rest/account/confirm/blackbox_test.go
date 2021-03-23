package confirm_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/confirm"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
)

const (
	confirmRt = confirm.Confirm + rest.Root
	sendRt    = confirm.Confirm + confirm.Send
)

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		confirm.RegisterServer,
	}, true)
	c := srv.Config()
	c.Signup.AutoConfirm = false
	return srv, web, smtp
}

func TestConfirmServer_ConfirmUser(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	// invalid req
	_, err := thttp.DoRequest(t, web, http.MethodPost, confirmRt, nil, []byte("\n"))
	assert.Error(t, err)
	// empty token
	req := new(confirm.Request)
	_, err = thttp.DoRequest(t, web, http.MethodPost, confirmRt, nil, req)
	assert.Error(t, err)
	// bad token
	req = &confirm.Request{
		Token: "bad",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, confirmRt, nil, req)
	assert.Error(t, err)
	var tok string
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(template.ConfirmUserAction, email)
	})
	u, _ := tcore.TestUser(t, srv.API, "", false)
	assert.False(t, u.IsConfirmed())
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
	req = &confirm.Request{
		Token: tok,
	}
	res, err := thttp.DoRequest(t, web, http.MethodPost, confirmRt, nil, req)
	assert.NoError(t, err)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	_, claims := tsrv.UnmarshalTokenResponse(t, srv.Config().JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject)
}

func TestConfirmServer_SendConfirmUser(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	// invalid req
	_, err := thttp.DoRequest(t, web, http.MethodPost, sendRt, nil, []byte("\n"))
	assert.Error(t, err)
	// empty email
	req := new(confirm.Request)
	_, err = thttp.DoRequest(t, web, http.MethodPost, sendRt, nil, req)
	assert.Error(t, err)
	// bad email
	req = &confirm.Request{
		Email: "bad",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, sendRt, nil, req)
	assert.Error(t, err)
	// email not found
	req = &confirm.Request{
		Email: "i-dont-exist@example.com",
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, sendRt, nil, req)
	assert.NoError(t, err)
	srv.Config().Mail.SendLimit = 0
	var tok1 string
	smtp.AddHook(t, func(email string) {
		tok1 = tconf.GetEmailToken(template.ConfirmUserAction, email)
	})
	u, _ := tcore.TestUser(t, srv.API, "", false)
	assert.False(t, u.IsConfirmed())
	assert.Eventually(t, func() bool {
		return tok1 != ""
	}, 1*time.Second, 10*time.Millisecond)
	var tok2 string
	smtp.AddHook(t, func(email string) {
		tok2 = tconf.GetEmailToken(template.ConfirmUserAction, email)
	})
	req = &confirm.Request{
		Email: u.Email,
	}
	_, err = thttp.DoRequest(t, web, http.MethodPost, sendRt, nil, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok2 != ""
	}, 1*time.Second, 10*time.Millisecond)
	assert.Equal(t, tok1, tok2)
}

func TestConfirmServer_SendConfirmUser_RateLimit(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	srv.Config().Mail.SendLimit = 5 * time.Minute
	var sent string
	smtp.AddHook(t, func(email string) {
		sent = email
	})
	// sent initial
	u, _ := tcore.TestUser(t, srv.API, "", false)
	assert.False(t, u.IsConfirmed())
	assert.Eventually(t, func() bool {
		return sent != ""
	}, 1*time.Second, 10*time.Millisecond)
	// resend
	sent = ""
	req := &confirm.Request{
		Email: u.Email,
	}
	_, err := thttp.DoRequest(t, web, http.MethodPost, sendRt, nil, req)
	assert.EqualError(t, err, thttp.FmtError(http.StatusTooEarly).Error())
	assert.Never(t, func() bool {
		return sent != ""
	}, 1*time.Second, 10*time.Millisecond)
}
