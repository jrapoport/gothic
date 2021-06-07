package hosts

import (
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func configClient(t *testing.T, h core.Hosted) admin.AdminClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return admin.NewAdminClient(cc)
	}).(admin.AdminClient)
}

func TestAdminHost(t *testing.T) {
	t.Parallel()
	a, c, _ := tcore.API(t, false)
	// create an rcp-web host
	h := NewAdminHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err := h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, 1*time.Second, 10*time.Millisecond)
	test := a.Settings()
	// unauthenticated call
	ctx := metadata.NewOutgoingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, c.RootPassword))
	ac := configClient(t, h)
	res, err := ac.Settings(ctx, &admin.SettingsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, test.Status, res.Status)
	assert.Equal(t, test.Signup.Provider.Internal, res.Signup.Provider.Internal)
	// shut down
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, 1*time.Second, 10*time.Millisecond)
}
