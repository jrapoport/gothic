package account

import (
	"context"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/protobuf/grpc/rpc/account"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestAccountServer_SendConfirmUser(t *testing.T) {
	t.Parallel()
	s, smtp := tsrv.RPCServer(t, true)
	srv := newAccountServer(s)
	srv.Config().Signup.AutoConfirm = false
	srv.Config().Mail.SendLimit = 0
	ctx := context.Background()
	// invalid req
	_, err := srv.SendConfirmUser(ctx, nil)
	assert.Error(t, err)
	// empty email
	req := &account.SendConfirmRequest{}
	_, err = srv.SendConfirmUser(ctx, req)
	assert.Error(t, err)
	// bad email
	req.Email = "bad"
	_, err = srv.SendConfirmUser(ctx, req)
	assert.Error(t, err)
	// not found
	req.Email = "i-dont-exist@example.com"
	_, err = srv.SendConfirmUser(ctx, req)
	assert.Error(t, err)
	// success
	var tok string
	act := template.ConfirmUserAction
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(act, email)
	})
	u := testUser(t, srv)
	req.Email = u.Email
	_, err = srv.SendConfirmUser(ctx, req)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
}

func TestAccountServer_SendConfirmUser_RateLimit(t *testing.T) {
	t.Parallel()
	s, smtp := tsrv.RPCServer(t, true)
	srv := newAccountServer(s)
	srv.Config().Signup.AutoConfirm = false
	srv.Config().Mail.SendLimit = 5 * time.Minute
	ctx := context.Background()
	var sent string
	smtp.AddHook(t, func(email string) {
		sent = email
	})
	// sent initial
	u := testUser(t, srv)
	assert.False(t, u.IsConfirmed())
	assert.Eventually(t, func() bool {
		return sent != ""
	}, 1*time.Second, 10*time.Millisecond)
	// resend
	sent = ""
	req := &account.SendConfirmRequest{
		Email: u.Email,
	}
	_, err := srv.SendConfirmUser(ctx, req)
	test := s.RPCError(codes.DeadlineExceeded,
		config.ErrRateLimitExceeded)
	require.NotNil(t, test)
	assert.EqualError(t, err, test.Error())
	assert.Never(t, func() bool {
		return sent != ""
	}, 1*time.Second, 10*time.Millisecond)
}

func TestAccountServer_ConfirmUser(t *testing.T) {
	t.Parallel()
	s, smtp := tsrv.RPCServer(t, true)
	srv := newAccountServer(s)
	srv.Config().Signup.AutoConfirm = false
	srv.Config().Mail.SendLimit = 0
	ctx := context.Background()
	// invalid req
	_, err := srv.ConfirmUser(ctx, nil)
	assert.Error(t, err)
	// empty token
	req := &account.ConfirmUserRequest{}
	_, err = srv.ConfirmUser(ctx, req)
	assert.Error(t, err)
	// bad token
	req.Token = "bad"
	_, err = srv.ConfirmUser(ctx, req)
	assert.Error(t, err)
	var tok string
	smtp.AddHook(t, func(email string) {
		const act = template.ConfirmUserAction
		tok = tconf.GetEmailToken(act, email)
	})
	u, _ := tcore.TestUser(t, srv.API, "", false)
	assert.False(t, u.IsConfirmed())
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
	req = &account.ConfirmUserRequest{
		Token: tok,
	}
	res, err := srv.ConfirmUser(ctx, req)
	assert.NoError(t, err)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Access)
	assert.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, u.ID.String(), claims.Subject)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
}
