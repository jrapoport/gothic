package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"

func testServer(t *testing.T) *userServer {
	srv, _ := tsrv.RPCServer(t, false)
	srv.Config().Signup.AutoConfirm = true
	return newUserServer(srv)
}

func TestUserServer_GetUser(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	srv.Config().MaskEmails = false
	req := &GetUserRequest{}
	ctx := context.Background()
	// no id
	_, err := srv.GetUser(ctx, req)
	assert.Error(t, err)
	// unmasked
	u, tok := tcore.TestUser(t, srv.API, "", false)
	ctx = tsrv.RPCAuthContext(t, srv.Config(), tok)
	res, err := srv.GetUser(ctx, req)
	assert.NoError(t, err)
	test, err := rpc.NewUserResponse(u)
	assert.NoError(t, err)
	assert.Equal(t, test.Email, res.Email)
	assert.Equal(t, test.Username, res.Username)
	assert.Equal(t, test.Data.AsMap(), res.Data.AsMap())
	// masked
	srv.Config().MaskEmails = true
	test.Email = utils.MaskEmail(u.Email)
	res, err = srv.GetUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, test.Email, res.Email)
	assert.Equal(t, test.Username, res.Username)
	assert.Equal(t, test.Data.AsMap(), res.Data.AsMap())
}

func TestUserServer_UpdateUser(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	u, tok := tcore.TestUser(t, srv.API, "", false)
	ctx := tsrv.RPCAuthContext(t, srv.Config(), tok)
	data, err := structpb.NewStruct(types.Map{
		"foo":   "bar",
		"tasty": "salad",
	})
	require.NoError(t, err)
	req := &UpdateUserRequest{
		Username: "peaches",
		Data:     data,
	}
	// invalid username
	srv.Config().Validation.UsernameRegex = "0"
	_, err = srv.UpdateUser(ctx, req)
	assert.Error(t, err)
	// unmasked
	srv.Config().MaskEmails = false
	srv.Config().Validation.UsernameRegex = ""
	res, err := srv.UpdateUser(ctx, req)
	assert.NoError(t, err)
	err = mergo.Map(&u.Data, req.Data.AsMap())
	assert.NoError(t, err)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, req.Username, res.Username)
	assert.EqualValues(t, u.Data, res.Data.AsMap())
	// masked
	srv.Config().MaskEmails = true
	data, err = structpb.NewStruct(types.Map{
		"quack":   99.0,
		"peaches": "happy",
	})
	require.NoError(t, err)
	req = &UpdateUserRequest{
		Username: "mario",
		Data:     data,
	}
	res, err = srv.UpdateUser(ctx, req)
	assert.NoError(t, err)
	err = mergo.Map(&u.Data, req.Data.AsMap())
	assert.NoError(t, err)
	assert.Equal(t, utils.MaskEmail(u.Email), res.Email)
	assert.Equal(t, req.Username, res.Username)
	assert.EqualValues(t, u.Data, res.Data.AsMap())
}

func TestUserServer_SendConfirmUser(t *testing.T) {
	t.Parallel()
	s, smtp := tsrv.RPCServer(t, true)
	s.Config().Signup.AutoConfirm = false
	s.Config().Mail.SendLimit = 0
	srv := newUserServer(s)
	u, _ := tcore.TestUser(t, srv.API, "", false)
	bt, err := srv.GrantBearerToken(context.Background(), u)
	require.NoError(t, err)
	require.NotNil(t, bt)
	ctx := tsrv.RPCAuthContext(t, srv.Config(), bt.Token)
	// success
	var tok string
	act := template.ConfirmUserAction
	smtp.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(act, email)
	})
	_, err = srv.SendConfirmUser(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 200*time.Millisecond, 10*time.Millisecond)
	srv.Config().Mail.SendLimit = 5 * time.Minute
	_, err = srv.SendConfirmUser(ctx, &emptypb.Empty{})
	assert.Error(t, err)
	u, err = srv.ConfirmUser(rpc.RequestContext(ctx), tok)
	require.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsConfirmed())
	_, err = srv.SendConfirmUser(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
}

func TestUserServer_ChangePassword(t *testing.T) {
	t.Parallel()
	const newPassword = "gj8#xtg#yrabxpnno!p5f3t8na!hd3?4jq7majxs"
	srv := testServer(t)
	u, tok := tcore.TestUser(t, srv.API, "", false)
	ctx := tsrv.RPCAuthContext(t, srv.Config(), tok)
	req := &ChangePasswordRequest{
		Password: newPassword,
	}
	// invalid password
	srv.Config().Validation.PasswordRegex = "0"
	_, err := srv.ChangePassword(ctx, req)
	assert.Error(t, err)
	// wrong old password
	srv.Config().Validation.PasswordRegex = ""
	_, err = srv.ChangePassword(ctx, req)
	assert.Error(t, err)
	// success
	req.OldPassword = testPass
	res, err := srv.ChangePassword(ctx, req)
	assert.NoError(t, err)
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Access)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, u.ID.String(), claims.Subject)
	u, err = srv.API.GetUser(u.ID)
	assert.NoError(t, err)
	err = u.Authenticate(newPassword)
	assert.NoError(t, err)
}

func TestRequestErrors(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	ctx := context.Background()
	_, err := srv.GetUser(ctx, &GetUserRequest{})
	assert.Error(t, err)
	_, err = srv.UpdateUser(ctx, &UpdateUserRequest{})
	assert.Error(t, err)
	_, err = srv.ChangePassword(ctx, &ChangePasswordRequest{})
	assert.Error(t, err)
	claims := jwt.UserClaims{}
	claims.Subject = uuid.New().String()
	ctx = context.WithContext(rpc.WithClaims(ctx, claims))
	_, err = srv.GetUser(ctx, &GetUserRequest{})
	assert.Error(t, err)
	_, err = srv.UpdateUser(ctx, &UpdateUserRequest{})
	assert.Error(t, err)
	_, err = srv.ChangePassword(ctx, &ChangePasswordRequest{})
	assert.Error(t, err)
}
