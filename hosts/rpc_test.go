package hosts

import (
	"github.com/jrapoport/gothic/hosts/rpc/admin"
	"testing"
	"time"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func adminClient(t *testing.T, h core.Hosted) admin.AdminClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return admin.NewAdminClient(cc)
	}).(admin.AdminClient)
}

func TestRPCHost(t *testing.T) {
	a, _, _ := tcore.API(t, false)
	// create an rcp-web host
	h := NewRPCHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err := h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, time.Second, 10*time.Millisecond)
	test := a.Settings()
	// unauthenticated call
	ctx := context.Background()
	ac := adminClient(t, h)
	res, err := ac.Settings(ctx, &admin.SettingsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, test.Status, res.Status)
	assert.Equal(t, test.Signup.Provider.Internal, res.Signup.Provider.Internal)
	// shut down
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, time.Second, 10*time.Millisecond)
}
