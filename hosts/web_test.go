package hosts

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc/user"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func accountClient(t *testing.T, h core.Hosted) account.AccountClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return account.NewAccountClient(cc)
	}).(account.AccountClient)
}

func userClient(t *testing.T, h core.Hosted) user.UserClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return user.NewUserClient(cc)
	}).(user.UserClient)
}

func TestRPCWebHost(t *testing.T) {
	t.Parallel()
	a, c, _ := tcore.API(t, false)
	// create an rcp-web host
	h := NewRPCWebHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err := h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, 1*time.Second, 10*time.Millisecond)
	// create a test user
	c.Signup.AutoConfirm = true
	c.Security.MaskEmails = false
	const pass = "1234567890asdfghjkl"
	test, _ := tcore.TestUser(t, a, pass, false)
	// unauthenticated call
	ctx := context.Background()
	ac := accountClient(t, h)
	ur, err := ac.Login(ctx, &account.LoginRequest{
		Email:    test.Email,
		Password: pass,
	})
	assert.NoError(t, err)
	assert.Equal(t, test.Email, ur.Email)
	require.NotEmpty(t, ur.Token)
	// authenticated call (error)
	uc := userClient(t, h)
	_, err = uc.GetUser(ctx, &user.UserRequest{})
	assert.Error(t, err)
	// authenticated call (success)
	ctx = tsrv.RPCAuthContext(t, c, ur.Token.Access)
	res, err := uc.GetUser(ctx, &user.UserRequest{})
	assert.NoError(t, err)
	assert.Equal(t, test.Email, res.Email)
	// shut down
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, 1*time.Second, 10*time.Millisecond)
}
