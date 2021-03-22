package health_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/health"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testResponse = `{"name":"gothic","version":"debug","status":"bela lugosi's` +
	` dead","hosts":{"test1":{"name":"test1","address":"` + tsrv.MockAddress +
	`","online":true},"test2":{"name":"mock","address":"` + tsrv.MockAddress +
	`","online":true},"test3":{"name":"test3","online":false}}}`

func TestNewHealthHost(t *testing.T) {
	t.Parallel()
	c := tconf.TempDB(t)
	a, err := core.NewAPI(c)
	require.NoError(t, err)
	h := health.NewHealthHost(a, "127.0.0.1:0", map[string]core.Hosted{
		"test1": &tsrv.MockHost{},
		"test2": tsrv.NewMockHost("mock"),
		"test3": nil,
	})
	require.NotNil(t, h)
	err = h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, 1*time.Second, 10*time.Millisecond)
	assert.JSONEq(t, testResponse, checkHealth(t, h))
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, 1*time.Second, 10*time.Millisecond)
}

func checkHealth(t *testing.T, h core.Hosted) string {
	healthURI := func() string {
		require.NotEmpty(t, h.Address())
		return "http://" + h.Address() + config.HealthCheck
	}
	res, err := http.Get(healthURI())
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	err = res.Body.Close()
	require.NoError(t, err)
	return string(b)
}
