package tconf

import (
	"fmt"
	"github.com/jrapoport/gothic/models/types/key"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/azureadv2"
	"github.com/markbates/goth/providers/faux"
	"github.com/stretchr/testify/require"
)

const testRole = "mock-user"

type MockSession struct {
	goth.Session
	Role string
	t    *testing.T
}

// MockProvider is a provider mock.
type MockProvider struct {
	faux      faux.Provider
	Role      string
	AccountID string
	Email     string
	Username  string
	Callback  string
	t         *testing.T
}

var _ goth.Provider = (*MockProvider)(nil)

func NewMockProvider(t *testing.T, callback string) *MockProvider {
	return &MockProvider{
		faux:      faux.Provider{},
		Role:      testRole,
		AccountID: uuid.NewString(),
		Username:  utils.RandomUsername(),
		Email:     tutils.RandomEmail(),
		Callback:  callback,
		t:         t,
	}
}

func (fp MockProvider) Name() string {
	return fp.faux.Name()
}

// PName returns a typed provider name for tests.
func (fp MockProvider) PName() provider.Name {
	return provider.Name(fp.Name())
}

// BeginAuth satisfies the goth.Provider
func (fp MockProvider) BeginAuth(state string) (goth.Session, error) {
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
	return &MockSession{
		Session: &faux.Session{
			ID:      fp.AccountID,
			Name:    fp.Username,
			Email:   fp.Email,
			AuthURL: authURL,
		},
		t: fp.t,
	}, nil
}

func (fp MockProvider) SetName(name string) {
	fp.faux.SetName(name)
}

func (fp MockProvider) UnmarshalSession(s string) (goth.Session, error) {
	sess, err := fp.faux.UnmarshalSession(s)
	if err != nil {
		return nil, err
	}
	return &MockSession{
		sess,
		fp.Role,
		fp.t,
	}, nil
}

func (fp MockProvider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*MockSession)
	return fp.faux.FetchUser(sess.Session)
}

func (fp MockProvider) Debug(b bool) {
	fp.faux.Debug(b)
}

func (fp MockProvider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return fp.faux.RefreshToken(refreshToken)
}

func (fp MockProvider) RefreshTokenAvailable() bool {
	return fp.faux.RefreshTokenAvailable()
}

// Authorize is used only for testing.
func (s *MockSession) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	tok := params.Get(key.Role)
	require.Equal(s.t, s.Role, tok)
	return s.Session.Authorize(provider, params)
}

// MockedProvider returns a mocked provider for tests.
func MockedProvider(t *testing.T, c *config.Config, callback string) (*config.Config, *MockProvider) {
	const (
		testClientKey = "provider-test-client-key"
		testSecret    = "provider-test-secret"
		testCallback  = "http://auth.exmaple.com/test/callback"
	)
	mp := NewMockProvider(t, callback)
	provider.AddExternal(mp.PName())
	t.Cleanup(func() {
		delete(provider.External, mp.PName())
	})
	if callback == "" {
		callback = testCallback
	}
	c.Authorization.Providers[mp.PName()] = config.Provider{
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
