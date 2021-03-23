package tconf

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/azureadv2"
	"github.com/markbates/goth/providers/faux"
	"github.com/stretchr/testify/require"
)

// MockedProvider returns a mocked provider for tests.
func MockedProvider(t *testing.T, c *config.Config, callback string) (*config.Config, goth.Provider) {
	const (
		testClientKey = "provider-test-client-key"
		testSecret    = "provider-test-secret"
		testCallback  = "http://auth.exmaple.com/test/callback"
	)
	mp := newMockProvider(t, callback)
	p := provider.Name(mp.Name())
	provider.AddExternal(p)
	t.Cleanup(func() {
		delete(provider.External, p)
	})
	if callback == "" {
		callback = testCallback
	}
	c.Authorization.Providers[p] = config.Provider{
		ClientKey:   testClientKey,
		Secret:      testSecret,
		CallbackURL: callback,
	}
	return c, mp
}

// MockOpenIDConnect mocks an openid connect response.
func MockOpenIDConnect(t *testing.T) string {
	const discovery = `{
		"issuer": "https://example.com/",
		"authorization_endpoint": "https://example.com/authorize",
		"token_endpoint": "https://example.com/token",
		"userinfo_endpoint": "https://example.com/userinfo",
		"jwks_uri": "https://example.com/.well-known/jwks.json",
		"scopes_supported": [
			"pets_read",
			"pets_write",
			"admin"
		],
		"response_types_supported": [
			"code",
			"id_token",
			"token id_token"
		],
		"token_endpoint_auth_methods_supported": [
			"client_secret_basic"
		]
	}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(discovery))
		require.NoError(t, err)
	}))
	t.Cleanup(func() {
		srv.Close()
	})
	return srv.URL + "/.well-known/openid-configuration"
}

func setEnv(t *testing.T, key, value string) {
	err := os.Setenv(config.ENVPrefix+"_"+key, value)
	require.NoError(t, err)
}

// ProvidersConfig returns an external provider config.
func ProvidersConfig(t *testing.T) *config.Config {
	const (
		testClientKey = "provider-test-client-key"
		testSecret    = "provider-test-secret"
		testCallback  = "http://auth.exmaple.com/test/callback"
	)
	setEnv(t, config.Auth0DomainEnv, "example.com")
	setEnv(t, config.OpenIDConnectURLEnv, MockOpenIDConnect(t))
	setEnv(t, config.AzureADTenantEnv, string(azureadv2.CommonTenant))
	var testScopes = []string{"test-scope-1", "test-scope-2"}
	c := Config(t)
	a := config.Authorization{
		Providers: map[provider.Name]config.Provider{},
	}
	a.UseInternal = true
	for name := range provider.External {
		a.Providers[name] = config.Provider{
			ClientKey:   testClientKey,
			Secret:      testSecret,
			CallbackURL: testCallback,
			Scopes:      testScopes,
		}
	}
	c.Authorization = a
	return c
}

type mockSession struct {
	goth.Session
	Role string
	t    *testing.T
}

type mockProvider struct {
	faux      faux.Provider
	Role      string
	AccountID string
	Email     string
	Username  string
	Callback  string
	t         *testing.T
}

var _ goth.Provider = (*mockProvider)(nil)

const testRole = "mock-user"

func newMockProvider(t *testing.T, callback string) *mockProvider {
	return &mockProvider{
		faux:      faux.Provider{},
		Role:      testRole,
		AccountID: uuid.NewString(),
		Username:  utils.RandomUsername(),
		Email:     tutils.RandomEmail(),
		Callback:  callback,
		t:         t,
	}
}

func (fp mockProvider) Name() string {
	return fp.faux.Name()
}

// BeginAuth satisfies the goth.Provider
func (fp mockProvider) BeginAuth(state string) (goth.Session, error) {
	s, err := fp.faux.BeginAuth(state)
	if err != nil {
		return nil, err
	}
	authURL, err := s.GetAuthURL()
	if err != nil {
		return nil, err
	}
	authURL += fmt.Sprintf("&%s=%s", key.Role, fp.Role)
	if fp.Callback != "" {
		au, err := url.Parse(authURL)
		if err != nil {
			return nil, err
		}
		authURL = fp.Callback + "?" + au.RawQuery
	}
	return &mockSession{
		Session: &faux.Session{
			ID:      fp.AccountID,
			Name:    fp.Username,
			Email:   fp.Email,
			AuthURL: authURL,
		},
		t: fp.t,
	}, nil
}

func (fp mockProvider) SetName(name string) {
	fp.faux.SetName(name)
}

func (fp mockProvider) UnmarshalSession(s string) (goth.Session, error) {
	sess, err := fp.faux.UnmarshalSession(s)
	if err != nil {
		return nil, err
	}
	return &mockSession{
		sess,
		fp.Role,
		fp.t,
	}, nil
}

func (fp mockProvider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*mockSession)
	return fp.faux.FetchUser(sess.Session)
}

func (fp mockProvider) Debug(b bool) {
	fp.faux.Debug(b)
}

func (fp mockProvider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return fp.faux.RefreshToken(refreshToken)
}

func (fp mockProvider) RefreshTokenAvailable() bool {
	return fp.faux.RefreshTokenAvailable()
}

// Authorize is used only for testing.
func (s *mockSession) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	tok := params.Get(key.Role)
	require.Equal(s.t, s.Role, tok)
	return s.Session.Authorize(provider, params)
}

// ToMockProvider returns the mock provider for tests
func ToMockProvider(p goth.Provider) *mockProvider {
	return p.(*mockProvider)
}
