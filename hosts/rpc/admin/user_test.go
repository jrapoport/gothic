package admin

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

func TestAdminServer_CreateUser(t *testing.T) {
	t.Parallel()
	var (
		testUsername = utils.RandomUsername()
		testPassword = utils.SecureToken()
	)
	req := &admin.CreateUserRequest{}
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := rootContext(srv.Config())
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

func TestAdminServer_ChangeUserRole(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := rootContext(srv.Config())
	// nil request
	_, err := srv.ChangeUserRole(ctx, nil)
	assert.Error(t, err)
	// no params
	req := &admin.ChangeUserRoleRequest{
		Role: user.RoleAdmin.String(),
	}
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// no root password
	req.User = &admin.ChangeUserRoleRequest_UserId{UserId: uuid.Nil.String()}
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// bad root password
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, "bad"))
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// nil user id
	root := s.Config().RootPassword
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, root))
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// invalid email
	req.User = &admin.ChangeUserRoleRequest_Email{Email: uuid.Nil.String()}
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// bad user id
	req.User = &admin.ChangeUserRoleRequest_UserId{UserId: uuid.New().String()}
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// bad email
	req.User = &admin.ChangeUserRoleRequest_Email{Email: tutils.RandomEmail()}
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	// success id
	u, _ := tcore.TestUser(t, srv.API, "", false)
	req.User = &admin.ChangeUserRoleRequest_UserId{UserId: u.ID.String()}
	res, err := srv.ChangeUserRole(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, u.ID.String(), res.GetUserId())
	// success email
	u, _ = tcore.TestUser(t, srv.API, "", false)
	req.User = &admin.ChangeUserRoleRequest_Email{Email: u.Email}
	res, err = srv.ChangeUserRole(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, u.ID.String(), res.GetUserId())
	// success hard
	u, _ = tcore.TestUser(t, srv.API, "", false)
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	_, err = srv.API.Signup(nil, u.Email, "", testPass, nil)
	assert.Error(t, err)
	req.User = &admin.ChangeUserRoleRequest_UserId{UserId: u.ID.String()}
	res, err = srv.ChangeUserRole(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	req.Role = user.RoleSuper.String()
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
	_, tok := tcore.TestUser(t, s.API, testPass, true)
	ctx = tsrv.RPCAuthContext(t, s.Config(), tok)
	req.User = &admin.ChangeUserRoleRequest_UserId{UserId: u.ID.String()}
	req.Role = user.RoleAdmin.String()
	_, err = srv.ChangeUserRole(ctx, req)
	assert.Error(t, err)
}

func TestAdminServer_DeleteUser(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := rootContext(srv.Config())
	// nil request
	_, err := srv.DeleteUser(ctx, nil)
	assert.Error(t, err)
	// no params
	req := &admin.DeleteUserRequest{}
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// no root password
	req.User = &admin.DeleteUserRequest_UserId{UserId: uuid.Nil.String()}
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// bad root password
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, "bad"))
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// nil user id
	root := s.Config().RootPassword
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, root))
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// invalid email
	req.User = &admin.DeleteUserRequest_Email{Email: uuid.Nil.String()}
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// bad user id
	req.User = &admin.DeleteUserRequest_UserId{UserId: uuid.New().String()}
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// bad email
	req.User = &admin.DeleteUserRequest_Email{Email: tutils.RandomEmail()}
	_, err = srv.DeleteUser(ctx, req)
	assert.Error(t, err)
	// success id
	u, _ := tcore.TestUser(t, srv.API, "", false)
	req.User = &admin.DeleteUserRequest_UserId{UserId: u.ID.String()}
	res, err := srv.DeleteUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, u.ID.String(), res.GetUserId())
	_, err = srv.GetUser(u.ID)
	assert.Error(t, err)
	// success email
	u, _ = tcore.TestUser(t, srv.API, "", false)
	req.User = &admin.DeleteUserRequest_Email{Email: u.Email}
	res, err = srv.DeleteUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, u.ID.String(), res.GetUserId())
	_, err = srv.GetUser(u.ID)
	assert.Error(t, err)
	// success hard
	u, _ = tcore.TestUser(t, srv.API, "", false)
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	_, err = srv.API.Signup(nil, u.Email, "", testPass, nil)
	assert.Error(t, err)
	req.User = &admin.DeleteUserRequest_UserId{UserId: u.ID.String()}
	req.Hard = true
	res, err = srv.DeleteUser(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, u.ID.String(), res.GetUserId())
	_, err = srv.GetUser(u.ID)
	assert.Error(t, err)
	_, err = srv.API.Signup(nil, u.Email, "", testPass, nil)
	assert.NoError(t, err)
}
