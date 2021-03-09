package hosts

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRESTHost(t *testing.T) {
	c := tconf.TempDB(t)
	a, err := core.NewAPI(c)
	require.NoError(t, err)
	h := NewRESTHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err = h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, time.Second, 10*time.Millisecond)
	testRESTCall(t, h)
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, time.Second, 10*time.Millisecond)
}
