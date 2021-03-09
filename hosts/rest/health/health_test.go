package health

import (
	"net/http"
	"testing"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExpectedResponse(t *testing.T, a *core.API) string {
	test := a.HealthCheck()
	b, err := json.Marshal(test)
	require.NoError(t, err)
	return string(b)
}

func TestHealthServer_HealthCheck(t *testing.T) {
	s, _ := tsrv.RESTServer(t, false)
	srv := newHealthServer(s)
	test := ExpectedResponse(t, s.API)
	assert.HTTPBodyContains(t, srv.HealthCheck,
		http.MethodGet, Endpoint, nil, test)
}
