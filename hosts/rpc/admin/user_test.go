package admin

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestSignupServer_CreateUser(t *testing.T) {
	t.Parallel()
	var (
		testUsername = utils.RandomUsername()
		testPassword = utils.SecureToken()
	)
	req := &admin.CreateUserRequest{}
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := context.Background()
	_, err := srv.CreateUser(ctx, nil)
	assert.Error(t, err)
	_, tok := tcore.TestUser(t, srv.API, "", false)
	ctx = tsrv.RPCAuthContext(t, srv.Config(), tok)
	_, err = srv.CreateUser(ctx, req)
	assert.Error(t, err)
	_, tok = tcore.TestUser(t, srv.API, "", true)
	ctx = tsrv.RPCAuthContext(t, srv.Config(), tok)
	_, err = srv.CreateUser(ctx, req)
	assert.Error(t, err)
	req.Email = tutils.RandomEmail()
	req.Admin = true
	_, err = srv.CreateUser(ctx, req)
	assert.Error(t, err)
	req.Admin = false
	res, err := srv.CreateUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, req.Email, res.Email)
	assert.Equal(t, user.RoleUser.String(), res.Role)
	req.Admin = true
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, "bad"))
	_, err = srv.CreateUser(ctx, req)
	assert.Error(t, err)
	root := s.Config().RootPassword
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, root))
	res, err = srv.CreateUser(ctx, req)
	assert.Error(t, err)
	req.Email = tutils.RandomEmail()
	req.Password = testPassword
	req.Username = &testUsername
	req.Data, err = structpb.NewStruct(types.Map{
		"tasty": "salad",
	})
	require.NoError(t, err)
	res, err = srv.CreateUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, req.Email, res.Email)
	assert.Equal(t, user.RoleAdmin.String(), res.Role)
	u, err := srv.API.GetUserWithEmail(res.Email)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, res.Email, u.Email)
	assert.Equal(t, res.Role, u.Role.String())
	assert.Equal(t, testUsername, u.Username)
	assert.Equal(t, "salad", u.Data["tasty"])
	err = u.Authenticate(testPassword)
	assert.NoError(t, err)
}
