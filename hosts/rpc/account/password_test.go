package account

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestAccountServer_SendResetPassword(t *testing.T) {
	t.Parallel()
	s, smtp := tsrv.RPCServer(t, true)
	srv := newAccountServer(s)
	ctx := context.Background()
	// invalid req
	_, err := srv.SendResetPassword(ctx, nil)
	assert.Error(t, err)
	// empty email
	req := &account.ResetPasswordRequest{}
	_, err = srv.SendResetPassword(ctx, req)
	assert.Error(t, err)
	// bad email
	req.Email = "bad"
	_, err = srv.SendResetPassword(ctx, req)
	assert.Error(t, err)
	// not found
	req.Email = "i-dont-exist@example.com"
	_, err = srv.SendResetPassword(ctx, req)
	assert.Error(t, err)
	// success
	var tok string
	act := template.ResetPasswordAction
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(act, email)
	})
	u := testUser(t, srv)
	req.Email = u.Email
	_, err = srv.SendResetPassword(ctx, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
}

func TestAccountServer_SendResetPassword_RateLimit(t *testing.T) {
	s, smtp := tsrv.RPCServer(t, true)
	srv := newAccountServer(s)
	ctx := context.Background()
	srv.Config().Mail.SendLimit = 5 * time.Minute
	srv.Config().Signup.AutoConfirm = true
	u := testUser(t, srv)
	var sent string
	smtp.AddHook(t, func(email string) {
		sent = email
	})
	for i := 0; i < 2; i++ {
		sent = ""
		req := &account.ResetPasswordRequest{
			Email: u.Email,
		}
		_, err := srv.SendResetPassword(ctx, req)
		if i == 0 {
			assert.NoError(t, err)
			assert.Eventually(t, func() bool {
				return sent != ""
			}, 1*time.Second, 10*time.Millisecond)
		} else {
			test := s.RPCError(codes.DeadlineExceeded,
				config.ErrRateLimitExceeded)
			require.NotNil(t, test)
			assert.EqualError(t, err, test.Error())
			assert.Never(t, func() bool {
				return sent != ""
			}, 1*time.Second, 10*time.Millisecond)
		}
	}
}

func TestAccountServer_ConfirmResetPassword(t *testing.T) {
	t.Parallel()
	const newPass = "sxjAm7QJ4?3dH!aN8T3F5P!oNnpXbaRy#gtx#8jG"
	s, smtp := tsrv.RPCServer(t, true)
	srv := newAccountServer(s)
	srv.Config().Signup.AutoConfirm = false
	srv.Config().Mail.SendLimit = 0
	ctx := context.Background()
	// invalid req
	_, err := srv.ConfirmResetPassword(ctx, nil)
	assert.Error(t, err)
	// empty token
	req := &account.ConfirmPasswordRequest{}
	_, err = srv.ConfirmResetPassword(ctx, req)
	assert.Error(t, err)
	// bad token
	req.Token = "bad"
	_, err = srv.ConfirmResetPassword(ctx, req)
	assert.Error(t, err)
	// first get the change token
	u := testUser(t, srv)
	assert.False(t, u.IsConfirmed())
	var tok string
	act := template.ResetPasswordAction
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(act, email)
	})
	pw := &account.ResetPasswordRequest{
		Email: u.Email,
	}
	_, err = srv.SendResetPassword(ctx, pw)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
	// now use the token to change the password
	req = &account.ConfirmPasswordRequest{
		Token:    tok,
		Password: newPass,
	}
	res, err := srv.ConfirmResetPassword(ctx, req)
	assert.NoError(t, err)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Access)
	assert.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, u.ID.String(), claims.Subject())
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	err = u.Authenticate(newPass)
	assert.NoError(t, err)
}
