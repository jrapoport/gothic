package invite_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/modules/invite"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		invite.RegisterServer,
	}, true)
	c := srv.Config()
	c.Signup.AutoConfirm = true
	c.Signup.Invites = config.Users
	c.Validation.UsernameRegex = ""
	c.Signup.Username = false
	c.Mail.SendLimit = 0
	t.Cleanup(func() {
		web.Close()
	})
	err := srv.API.LoadConfig(c)
	require.NoError(t, err)
	return srv, web, smtp
}

func TestInviteServer_SendInviteUser(t *testing.T) {
	t.Parallel()
	srv, web, smtp := testServer(t)
	var inviteTok string
	smtp.AddHook(t, func(email string) {
		inviteTok = tconf.GetEmailToken(template.InviteUserAction, email)
	})
	u, tok := tcore.TestUser(t, srv.API, "", true)
	req := &invite.Request{
		Email: u.Email,
	}
	_, err := thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return inviteTok != ""
	}, 1*time.Second, 10*time.Millisecond)
	data, err := jwt.ParseData(srv.Config().JWT, inviteTok)
	assert.NoError(t, err)
	assert.Equal(t, req.Email, data.Get(key.Email))
	in := data.Get(key.Token)
	assert.NotEmpty(t, in)
	_, err = srv.CheckSignupCode(in)
	assert.NoError(t, err)
}

func TestInviteServer_SendInviteUser_Error(t *testing.T) {
	t.Parallel()
	srv, web, _ := testServer(t)
	srv.Config().Signup.Invites = config.Users
	tok := thttp.UserToken(t, srv.Config().JWT, false, true)
	// bad request
	_, err := thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, []byte("\n"))
	assert.Error(t, err)
	// no email
	req := &invite.Request{Email: ""}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.Error(t, err)
	// bad email
	req = &invite.Request{Email: "@"}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.Error(t, err)
	// user not found
	req = &invite.Request{Email: tutils.RandomEmail()}
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.Error(t, err)
	// only admins
	srv.Config().Signup.Invites = config.Admins
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.Error(t, err)
	// invites disabled
	srv.Config().Signup.Invites = config.Disabled
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.Error(t, err)
	// signups disabled
	srv.Config().Signup.Disabled = true
	_, err = thttp.DoAuthRequest(t, web, http.MethodPost, invite.Invite, tok, nil, req)
	assert.Error(t, err)
}
