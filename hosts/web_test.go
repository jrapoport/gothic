package hosts

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCWebHost(t *testing.T) {
	a, c, _ := tcore.API(t, false)
	c.Security.MaskEmails = false
	h := NewRPCWebHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err := h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, time.Second, 10*time.Millisecond)
	// non-auth call
	testRPCCall(t, h)
	// auth call
	testRPCAuthCall(t, h.(*rpc.Host))
	// shut down
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, time.Second, 10*time.Millisecond)
}
