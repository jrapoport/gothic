package hosts

import (
	"testing"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
