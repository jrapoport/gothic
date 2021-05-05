package system

import (
	"context"
	"github.com/jrapoport/gothic/test/tsrv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemServer_GetUser(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	srv.Config().MaskEmails = false
	ctx := context.Background()
	// no id or email
	req := &system.UserAccountRequest{}
	_, err := srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// bad id
	req.Id = &system.UserAccountRequest_UserId{UserId: "1"}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// id not found
	req.Id = &system.UserAccountRequest_UserId{
		UserId: uuid.New().String(),
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// success
	u, _ := tcore.TestUser(t, srv.API, "", false)
	req.Id = &system.UserAccountRequest_UserId{
		UserId: u.ID.String(),
	}
	res, err := srv.GetUserAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), res.Id)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.Username, res.Username)
	// bad email
	req = &system.UserAccountRequest{}
	req.Id = &system.UserAccountRequest_Email{
		Email: "@",
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// email not found
	req.Id = &system.UserAccountRequest_Email{
		Email: tutils.RandomEmail(),
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// success
	req.Id = &system.UserAccountRequest_Email{
		Email: u.Email,
	}
	res, err = srv.GetUserAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), res.Id)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.Username, res.Username)
}

func TestSystemServer_NotifyUser(t *testing.T) {
	const (
		testSubject = "Test Subject"
		testHTML    = "<html>test_html_notification</html>"
	)
	var testPlain = "test_plain_notification"
	t.Parallel()
	srv := testServer(t)
	ctx := context.Background()
	// no id or email
	req := &system.NotificationRequest{}
	_, err := srv.NotifyUser(ctx, req)
	assert.Error(t, err)
	// bad id
	req = &system.NotificationRequest{UserId: "1"}
	_, err = srv.NotifyUser(ctx, req)
	assert.Error(t, err)
	// offline
	req = &system.NotificationRequest{
		UserId: uuid.New().String(),
	}
	res, err := srv.NotifyUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.False(t, res.Sent)
	// mail online
	s, mock := tsrv.RPCServer(t, true)
	s.Config().Signup.AutoConfirm = true
	s.Config().Mail.SMTP.SpamProtection = false
	srv = newSystemServer(s)
	// id not found
	req = &system.NotificationRequest{
		UserId: uuid.New().String(),
	}
	_, err = srv.NotifyUser(ctx, req)
	assert.Error(t, err)
	// create a user
	u, _ := tcore.TestUser(t, srv.API, "", false)
	require.NotNil(t, u)
	req = &system.NotificationRequest{
		UserId: u.ID.String(),
	}
	_, err = srv.NotifyUser(ctx, req)
	assert.Error(t, err)
	// success
	req = &system.NotificationRequest{
		UserId:  u.ID.String(),
		Subject: testSubject,
		Html:    testHTML,
		Plain:   &testPlain,
	}
	var recv string
	var mu sync.Mutex
	mock.AddHook(t, func(email string) {
		mu.Lock()
		defer mu.Unlock()
		recv = email
	})
	res, err = srv.NotifyUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.True(t, res.Sent)
	assert.Eventually(t, func() bool {
		if !strings.Contains(recv, testSubject) {
			return false
		}
		if !strings.Contains(recv, testHTML) {
			return false
		}
		if !strings.Contains(recv, testPlain) {
			return false
		}
		return true
	}, 1*time.Second, 10*time.Millisecond)
}
