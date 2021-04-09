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
	c.Network.AdminAddress = "127.0.0.1:0"
	c.Network.RESTAddress = "127.0.0.1:0"
	c.Network.RPCAddress = "127.0.0.1:0"
	c.Network.RPCWebAddress = "127.0.0.1:0"
	c.Network.HealthAddress = "127.0.0.1:0"
	go func() {
		err := Main(c)
		assert.NoError(t, err)
	}()
	assert.Eventually(t, func() bool {
		return hosts.Running()
	}, 5*time.Second, 100*time.Millisecond)
	healthURI := func() string {
		require.NotEmpty(t, c.Network.HealthAddress)
		return "http://" + c.Network.HealthAddress + config.HealthCheck
	}
	_, err := http.Get(healthURI())
	assert.NoError(t, err)
	err = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !hosts.Running()
	}, 5*time.Second, 100*time.Millisecond)
}
