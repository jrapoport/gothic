package hosts

import (
	"context"
	"testing"
	"time"

	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func systemClient(t *testing.T, h core.Hosted) system.SystemClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return system.NewSystemClient(cc)
	}).(system.SystemClient)
}

func TestRPCHost(t *testing.T) {
	t.Parallel()
	a, _, _ := tcore.API(t, false)
	// create an rcp-web host
	h := NewRPCHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err := h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, 1*time.Second, 10*time.Millisecond)
	u, _ := tcore.TestUser(t, a, "", false)
	req := &system.UserAccountRequest{}
	req.Id = &system.UserAccountRequest_UserId{UserId: u.ID.String()}
	sc := systemClient(t, h)
	res, err := sc.GetUserAccount(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, res.Email)
	// shut down
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, 1*time.Second, 10*time.Millisecond)
}
