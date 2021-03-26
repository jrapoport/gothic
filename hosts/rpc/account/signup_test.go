package account

import (
	"testing"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/protobuf/grpc/rpc/account"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func testCase(t *testing.T) (*account.SignupRequest, *rpc.UserResponse) {
	em := tutils.RandomEmail()
	un := utils.RandomUsername()
	data, err := structpb.NewStruct(types.Map{
		"foo":   "bar",
		"tasty": "salad",
	})
	require.NoError(t, err)
	req := &account.SignupRequest{
		Email:    em,
		Password: testPass,
		Username: un,
		Data:     data,
	}
	res := &rpc.UserResponse{
		Role:     user.RoleUser.String(),
		Email:    em,
		Username: un,
		Data:     data,
	}
	return req, res
}

func TestAccountServer_Signup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	srv.Config().MaskEmails = false
	// invalid req
	_, err := srv.Signup(ctx, nil)
	assert.Error(t, err)
	// no email
	_, err = srv.Signup(ctx, &account.SignupRequest{})
	assert.Error(t, err)
	// bad email
	_, err = srv.Signup(ctx, &account.SignupRequest{
		Email: "bad",
	})
	assert.Error(t, err)
	req, test := testCase(t)
	res, err := srv.Signup(ctx, req)
	assert.NoError(t, err)
	assertUserResponse(t, srv, test, res)
	// email taken
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	// json success (masked)
	srv.Config().MaskEmails = true
	em := tutils.RandomEmail()
	req.Email = em
	res, err = srv.Signup(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, utils.MaskEmail(em), res.Email)
}

func TestAccountServer_Signup_Confirm(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	srv.Config().Signup.AutoConfirm = false
	srv.Config().MaskEmails = false
	req, _ := testCase(t)
	res, err := srv.Signup(ctx, req)
	assert.NoError(t, err)
	u, err := srv.GetUserWithEmail(res.Email)
	assert.NoError(t, err)
	assert.False(t, u.IsConfirmed())
}

func TestAccountServer_Signup_AutoConfirm(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	req, _ := testCase(t)
	res, err := srv.Signup(ctx, req)
	assert.NoError(t, err)
	u, err := srv.GetUserWithEmail(res.Email)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
}

func TestAccountServer_Signup_Disabled(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	srv.Config().Signup.Disabled = true
	req, _ := testCase(t)
	_, err := srv.Signup(ctx, req)
	assert.Error(t, err)
}

func TestAccountServer_Signup_EmailDisabled(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	srv.Config().UseInternal = false
	req, _ := testCase(t)
	_, err := srv.Signup(ctx, req)
	assert.Error(t, err)
}

func TestAccountServer_Signup_Recaptcha(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	srv.Config().Recaptcha.Key = validate.ReCaptchaDebugKey
	req, _ := testCase(t)
	// invalid client ip
	_, err := srv.Signup(ctx, req)
	assert.Error(t, err)
	// no token
	ctx.SetIPAddress("127.0.0.1")
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	// invalid token
	ctx.SetReCaptcha("invalid")
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	// token
	ctx.SetReCaptcha(validate.ReCaptchaDebugToken)
	_, err = srv.Signup(ctx, req)
	assert.NoError(t, err)
}

func TestAccountServer_Signup_SignupCode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	// no code required
	srv.Config().Signup.Code = false
	// no code
	req, _ := testCase(t)
	_, err := srv.Signup(ctx, req)
	assert.NoError(t, err)
	// code required
	srv.Config().Signup.Code = true
	// missing code
	req, _ = testCase(t)
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	// bad code
	ctx.SetCode("bad")
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	// good code
	pin, err := srv.CreateSignupCode(ctx, code.SingleUse)
	assert.NoError(t, err)
	ctx.SetCode(pin)
	_, err = srv.Signup(ctx, req)
	assert.NoError(t, err)
	// reuse code
	req, _ = testCase(t)
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
}

func TestAccountServer_Signup_Password(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	srv := testServer(t)
	const passRegex = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	srv.Config().Validation.PasswordRegex = passRegex
	// good pw
	req, _ := testCase(t)
	_, err := srv.Signup(ctx, req)
	assert.NoError(t, err)
	// missing pw
	req, _ = testCase(t)
	req.Password = ""
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	// bad password
	req, _ = testCase(t)
	req.Password = "bad"
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	srv.Config().Validation.PasswordRegex = ""
	// blank  ok
	req, _ = testCase(t)
	req.Password = ""
	_, err = srv.Signup(ctx, req)
	assert.NoError(t, err)
	// custom password
	srv.Config().Validation.PasswordRegex = "^[a-z]"
	req, _ = testCase(t)
	req.Password = "12345678"
	_, err = srv.Signup(ctx, req)
	assert.Error(t, err)
	req.Password = "password"
	_, err = srv.Signup(ctx, req)
	assert.NoError(t, err)
}
