package health_test

import (
	"net/http"
	"testing"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/health"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
)

func TestHealthServer_HealthCheck(t *testing.T) {
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		health.RegisterServer,
	}, false)
	res, err := thttp.DoRequest(t, web, http.MethodGet, health.Endpoint, nil, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, health.ExpectedResponse(t, srv.API), res)
}
