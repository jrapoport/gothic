package gothic

import (
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Main(t *testing.T) {
	c := tconf.TempDB(t)
	c.Network.REST = "127.0.0.1:0"
	c.Network.RPC = "127.0.0.1:0"
	c.Network.RPCWeb = "127.0.0.1:0"
	c.Network.Health = "127.0.0.1:0"
	go func() {
		err := Main(c)
		assert.NoError(t, err)
	}()
	assert.Eventually(t, func() bool {
		return hosts.Running()
	}, 5*time.Second, 100*time.Millisecond)
	healthURI := func() string {
		require.NotEmpty(t, c.Network.Health)
		return "http://" + c.Network.Health + config.HealthEndpoint
	}
	_, err := http.Get(healthURI())
	assert.NoError(t, err)
	err = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !hosts.Running()
	}, 5*time.Second, 100*time.Millisecond)
}
