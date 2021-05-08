package system

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
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
	req.User = &system.UserAccountRequest_UserId{UserId: "1"}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// id not found
	req.User = &system.UserAccountRequest_UserId{
		UserId: uuid.New().String(),
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// success
	u, _ := tcore.TestUser(t, srv.API, "", false)
	req.User = &system.UserAccountRequest_UserId{
		UserId: u.ID.String(),
	}
	res, err := srv.GetUserAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), res.GetId())
	assert.Equal(t, u.Email, res.GetEmail())
	assert.Equal(t, u.Username, res.GetUsername())
	// bad email
	req = &system.UserAccountRequest{}
	req.User = &system.UserAccountRequest_Email{
		Email: "@",
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// email not found
	req.User = &system.UserAccountRequest_Email{
		Email: tutils.RandomEmail(),
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// success
	req.User = &system.UserAccountRequest_Email{
		Email: u.Email,
	}
	res, err = srv.GetUserAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), res.Id)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.Username, res.Username)
}
