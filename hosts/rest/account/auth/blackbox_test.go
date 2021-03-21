package auth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/auth"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testURL = "http://example.com/auth?client_id=&response_type=code&state="

func providerPath(p provider.Name) string {
	return auth.Endpoint + "/" + p.String()
}

func getToken(t *testing.T, authURL string) string {
	au, err := url.Parse(authURL)
	require.NoError(t, err)
	return au.Query().Get(key.State)
}

func DoProviderURLRequest(t *testing.T, web *httptest.Server, p provider.Name) string {
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(web.Config.Handler.ServeHTTP)
	req := thttp.Request(t, http.MethodGet, providerPath(p),
		"", nil, nil)
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusFound, rec.Code)
	return rec.Header().Get("location")
}

func TestAuthServer_GetAuthorizationURL(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		auth.RegisterServer,
	}, false)
	_, mock := tconf.MockedProvider(t, srv.Config(), "")
	srv.Providers().UseProviders(mock)
	// empty provider
	_, err := thttp.DoRequest(t, web, http.MethodGet, auth.Endpoint, nil, nil)
	assert.Error(t, err)
	// bad provider
	_, err = thttp.DoRequest(t, web, http.MethodGet, providerPath("bad"), nil, nil)
	assert.Error(t, err)
	// disabled provider
	_, err = thttp.DoRequest(t, web, http.MethodGet, providerPath(provider.Google), nil, nil)
	assert.Error(t, err)
	// valid provider
	assert.HTTPRedirect(t, web.Config.Handler.ServeHTTP, http.MethodGet, providerPath(mock.PName()), nil)
	urlReq1 := DoProviderURLRequest(t, web, mock.PName())
	assert.True(t, strings.HasPrefix(urlReq1, testURL))
	tok1 := getToken(t, urlReq1)
	urlReq2 := DoProviderURLRequest(t, web, mock.PName())
	assert.True(t, strings.HasPrefix(urlReq2, testURL))
	tok2 := getToken(t, urlReq2)
	assert.NotEqual(t, tok1, tok2)
}

func TestAuthServer_AuthorizeUser(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		auth.RegisterServer,
	}, false)
	srv.Config().Signup.Default.Color = true
	var callbackURL = web.URL + auth.Endpoint + auth.Callback
	// state not found
	assert.HTTPError(t, web.Config.Handler.ServeHTTP, http.MethodPost, callbackURL, nil)
	_, mock := tconf.MockedProvider(t, srv.Config(), callbackURL)
	srv.Providers().UseProviders(mock)
	urlReq := DoProviderURLRequest(t, web, mock.PName())
	urlReq = strings.Replace(urlReq, web.URL, "", 1)
	urlReq += "&test-value=expected"
	res, err := thttp.DoRequest(t, web, http.MethodPost, urlReq, nil, nil)
	assert.NoError(t, err)
	tr, claims := tsrv.UnmarshalUserResponse(t, srv.Config().JWT, res)
	assert.EqualValues(t, tokens.Bearer, tr.Token.Type)
	uid, err := uuid.Parse(claims.Subject)
	assert.NoError(t, err)
	u, err := srv.GetUser(uid)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), claims.Subject)
	au, err := srv.GetAuthenticatedUser(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
}
