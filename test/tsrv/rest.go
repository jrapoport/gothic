package tsrv

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/require"
)

// RESTServer a rest server for tests.
func RESTServer(t *testing.T, smtp bool) (*rest.Server, *tconf.SMTPMock) {
	s, mock := tcore.Server(t, smtp)
	return rest.NewServer(s), mock
}

// RESTHost a rest host for tests.
func RESTHost(t *testing.T, reg []rest.RegisterServer, smtp bool) (*rest.Host, *httptest.Server, *tconf.SMTPMock) {
	a, c, mock := tcore.API(t, smtp)
	rt := rest.NewRouter(c)
	web := httptest.NewServer(rt)
	for i, r := range reg {
		reg[i] = func(s *http.Server, srv *rest.Server) {
			r(web.Config, srv)
		}
	}
	t.Cleanup(func() {
		web.Close()
	})
	s := rest.NewHost(a, "test-rest", c.RESTAddress, reg)
	require.NotNil(t, s)
	return s, web, mock
}

// UnmarshalTokenResponse extracts the token from a token response.
func UnmarshalTokenResponse(t *testing.T, c config.JWT, res string) (*rest.BearerResponse, *jwt.UserClaims) {
	r := new(rest.BearerResponse)
	err := json.Unmarshal([]byte(res), r)
	require.NoError(t, err)
	claims, err := jwt.ParseUserClaims(c, r.Access)
	require.NoError(t, err)
	require.NotNil(t, claims)
	return r, claims
}

// UnmarshalUserResponse extracts the token from a token response.
func UnmarshalUserResponse(t *testing.T, c config.JWT, res string) (*rest.UserResponse, *jwt.UserClaims) {
	ur := new(rest.UserResponse)
	err := json.Unmarshal([]byte(res), ur)
	require.NoError(t, err)
	if ur.Token == nil {
		return ur, nil
	}
	claims, err := jwt.ParseUserClaims(c, ur.Token.Access)
	require.NoError(t, err)
	require.NotNil(t, claims)
	return ur, claims
}
