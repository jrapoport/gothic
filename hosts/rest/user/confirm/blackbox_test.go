package confirm_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/user/confirm"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const send = confirm.Endpoint + confirm.Root

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		confirm.RegisterServer,
	}, true)
	c := srv.Config()
	c.Signup.AutoConfirm = false
	return srv, web, smtp
}

func TestConfirmServer_SendConfirmUser(t *testing.T) {
	srv, web, smtp := testServer(t)
	srv.Config().Mail.SendLimit = 0
	j := srv.Config().JWT
	// not authorized
	_, err := thttp.DoRequest(t, web, http.MethodPost, send, nil, nil)
	assert.Error(t, err)
	// no user id
	bad := thttp.BadToken(t, j)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, send, bad, nil, nil)
	assert.Error(t, err)
	// user not found
	bad = thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, send, bad, nil, nil)
	assert.Error(t, err)
	// invalid req
	bad = thttp.UserToken(t, j, true, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, send, bad, nil, []byte("\n"))
	assert.Error(t, err)
	// user not confirmed
	bad = thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, send, bad, nil, nil)
	assert.Error(t, err)
	// empty email
	u, _ := tcore.TestUser(t, srv.API, "", false)
	bt, err := srv.GrantBearerToken(context.Background(), u)
	require.NoError(t, err)
	require.NotNil(t, bt)
	// success
	var tok string
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(template.ConfirmUserAction, email)
	})
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, send, bt.Token, nil, nil)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 100*time.Millisecond)
}

func TestConfirmServer_SendConfirmUser_RateLimit(t *testing.T) {
	srv, web, smtp := testServer(t)
	srv.Config().Mail.SendLimit = 5 * time.Minute
	var sent string
	smtp.AddHook(t, func(email string) {
		sent = email
	})
	// sent initial
	u, _ := tcore.TestUser(t, srv.API, "", false)
	bt, err := srv.GrantBearerToken(context.Background(), u)
	require.NoError(t, err)
	require.NotNil(t, bt)
	assert.False(t, u.IsConfirmed())
	assert.Eventually(t, func() bool {
		return sent != ""
	}, 1*time.Second, 100*time.Millisecond)
	// resend
	sent = ""
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, send, bt.Token, nil, nil)
	assert.EqualError(t, err, thttp.FmtError(http.StatusTooEarly).Error())
	assert.Never(t, func() bool {
		return sent != ""
	}, 1*time.Second, 100*time.Millisecond)
}
