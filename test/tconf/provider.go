package tconf

import (
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

// MockProvider is a provider mock.
type MockProvider struct {
	faux.Provider
	AccountID string
	Email     string
	Username  string
	Callback  string
}

// BeginAuth satisfies the goth.Provider
func (fp MockProvider) BeginAuth(state string) (goth.Session, error) {
	s, err := fp.Provider.BeginAuth(state)
	if err != nil {
		return nil, err
	}
	authURL, err := s.GetAuthURL()
	if err != nil {
		return nil, err
	}
	if fp.Callback != "" {
		au, err := url.Parse(authURL)
		if err != nil {
			return nil, err
		}
		authURL = fp.Callback + "?" + au.RawQuery
	}
	return &faux.Session{
		ID:      fp.AccountID,
		Name:    fp.Username,
		Email:   fp.Email,
		AuthURL: authURL,
	}, nil
}

// PName returns a typed provider name for tests.
func (fp MockProvider) PName() provider.Name {
	return provider.Name(fp.Provider.Name())
}

// MockedProvider returns a mocked provider for tests.
func MockedProvider(t *testing.T, c *config.Config, callback string) (*config.Config, *MockProvider) {
	const (
		testClientKey = "provider-test-client-key"
		testSecret    = "provider-test-secret"
		testCallback  = "http://auth.exmaple.com/test/callback"
	)
	fp := &MockProvider{
		AccountID: uuid.NewString(),
		Username:  utils.RandomUsername(),
		Email:     tutils.RandomEmail(),
		Callback:  callback,
	}
	provider.External[fp.PName()] = struct{}{}
	t.Cleanup(func() {
		delete(provider.External, fp.PName())
	})
	if callback == "" {
		callback = testCallback
	}
	c.Authorization.Providers[fp.PName()] = config.Provider{
		ClientKey:   testClientKey,
		Secret:      testSecret,
		CallbackURL: callback,
	}
	return c, fp
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
