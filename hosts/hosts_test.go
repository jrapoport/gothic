package hosts

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/health"
	"github.com/jrapoport/gothic/hosts/rpc/user"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestStart(t *testing.T) {
	c := tconf.TempDB(t)
	c.Network.REST = "127.0.0.1:0"
	c.Network.RPC = "127.0.0.1:0"
	c.Network.RPCWeb = "127.0.0.1:0"
	c.Network.Health = "127.0.0.1:0"
	a, err := core.NewAPI(c)
	require.NoError(t, err)
	err = Start(a, c)
	assert.NoError(t, err)
	t.Cleanup(func() {
		err = Shutdown()
		assert.NoError(t, err)
	})
}

func testRESTCall(t *testing.T, h core.Hosted) {
	healthURI := func() string {
		require.NotEmpty(t, h.Address())
		return "http://" + h.Address() + config.HealthEndpoint
	}
	res, err := http.Get(healthURI())
	require.NoError(t, err)
	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	err = res.Body.Close()
	require.NoError(t, err)
	hc := h.(*rest.Host).HealthCheck()
	test, err := json.Marshal(hc)
	require.NoError(t, err)
	assert.JSONEq(t, string(test), string(b))
}

func testRPCCall(t *testing.T, h core.Hosted) {
	hc := tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return health.NewHealthClient(cc)
	}).(health.HealthClient)
	res, err := hc.HealthCheck(context.Background(), &health.HealthCheckRequest{})
	require.NoError(t, err)
	test := h.(*rpc.Host).HealthCheck()
	assert.Equal(t, test.Name, res.Name)
	assert.Equal(t, test.Version, res.Version)
	assert.Equal(t, test.Status, res.Status)
}

func testRPCAuthCall(t *testing.T, h core.Hosted) {
	uc := tsrv.RPCClient(t, h.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return user.NewUserClient(cc)
	}).(user.UserClient)
	rh := h.(*rpc.Host)
	test, tok := tcore.TestUser(t, rh.API, "", false)
	ctx := tsrv.RPCAuthContext(t, rh.Config(), tok)
	res, err := uc.GetUser(ctx, &user.GetUserRequest{})
	require.NoError(t, err)
	assert.Equal(t, test.Email, res.Email)
	assert.Equal(t, test.Username, res.Username)
}
