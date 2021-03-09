package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/stretchr/testify/assert"
)

func TestProtect(t *testing.T) {
	c := tconf.Config(t)
	r := NewRouter(c)
	s := httptest.NewServer(r)
	j := c.JWT
	r = r.Authenticate(j)
	resOk := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	r.Get("/", resOk)
	// not authorized
	_, err := thttp.DoAuthRequest(t, s, http.MethodGet, "/", "", nil, nil)
	assert.Error(t, err)
	// invalid subject
	bad := thttp.BadToken(t, j)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/", bad, nil, nil)
	assert.Error(t, err)
	// authorized
	tok := thttp.UserToken(t, j, false, false)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/", tok, nil, nil)
	assert.NoError(t, err)
	// not confirmed
	r.Confirmed().Get("/confirmed", resOk)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/confirmed", tok, nil, nil)
	assert.Error(t, err)
	// confirmed
	tok = thttp.UserToken(t, j, true, false)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/confirmed", tok, nil, nil)
	assert.NoError(t, err)
	// not admin
	r.Admin().Get("/admin", resOk)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/admin", tok, nil, nil)
	assert.Error(t, err)
	tok = thttp.UserToken(t, j, true, true)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/admin", tok, nil, nil)
	assert.NoError(t, err)
	r.Admin().Confirmed().Get("/confirmed-admin", resOk)
	_, err = thttp.DoAuthRequest(t, s, http.MethodGet, "/confirmed-admin", tok, nil, nil)
	assert.NoError(t, err)
}
