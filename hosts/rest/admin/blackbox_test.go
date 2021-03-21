package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/admin"
	"github.com/jrapoport/gothic/hosts/rest/admin/settings"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	srv, web, smtp := tsrv.RESTHost(t, []rest.RegisterServer{
		admin.RegisterServer,
	}, true)
	c := srv.Config()
	t.Cleanup(func() {
		web.Close()
	})
	err := srv.API.LoadConfig(c)
	require.NoError(t, err)
	return srv, web, smtp
}

func testResponse(t *testing.T, s *rest.Host) string {
	test := s.Settings()
	b, err := json.Marshal(test)
	require.NoError(t, err)
	return string(b)
}

func TestAdminServer_Config(t *testing.T) {
	t.Parallel()
	const settings = admin.Endpoint + settings.Endpoint
	s, web, _ := testServer(t)
	j := s.Config().JWT
	// bad token
	bad := thttp.BadToken(t, j)
	_, err := thttp.DoAuthRequest(t, web, http.MethodGet, settings, bad, nil, nil)
	assert.Error(t, err)
	// not admin
	tok := thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, web, http.MethodGet, settings, tok, nil, nil)
	assert.Error(t, err)
	// admin
	tok = thttp.UserToken(t, j, false, true)
	res, err := thttp.DoAuthRequest(t, web, http.MethodGet, settings, tok, nil, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, testResponse(t, s), res)
}
