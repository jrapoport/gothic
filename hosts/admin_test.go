package hosts

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/protobuf/grpc/rpc/admin/settings"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func configClient(t *testing.T, h core.Hosted) settings.SettingsClient {
	return tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return settings.NewSettingsClient(cc)
	}).(settings.SettingsClient)
}

func TestAdminHost(t *testing.T) {
	t.Parallel()
	a, _, _ := tcore.API(t, false)
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
	ctx := context.Background()
	ac := configClient(t, h)
	res, err := ac.Settings(ctx, &settings.SettingsRequest{})
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
