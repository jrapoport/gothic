package config

import (
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// ENVPrefix is the prefix used for env vars.
const ENVPrefix = "GOTHIC"

const (
	// Auth0DomainEnv is the auth0 domain env var.
	Auth0DomainEnv = ENVPrefix + "_AUTH0_DOMAIN"

	// AzureADTenantEnv is the azure ad tenant env var.
	AzureADTenantEnv = ENVPrefix + "_AZURE_AD_TENANT"

	// CloudFoundryURLEnv is the cloud foundry env var.
	CloudFoundryURLEnv = ENVPrefix + "_CLOUDFOUNDRY_URL"

	// NextCloudURLEnv is the next cloud url env var.
	NextCloudURLEnv = ENVPrefix + "_NEXTCLOUD_URL"

	// OktaURLEnv is the okta url env var.
	OktaURLEnv = ENVPrefix + "_OKTA_URL"

	// OpenIDConnectURLEnv is the open id connect discovery url env var.
	OpenIDConnectURLEnv = ENVPrefix + "_OPENID_CONNECT_URL"

	// TwitterAuthorizeEnv will use twitter authorization instead of authentication.
	TwitterAuthorizeEnv = ENVPrefix + "_TWITTER_AUTHORIZE"
)

// IsDebug wil be true if the build is a debug build, false if release.
func (Config) IsDebug() bool {
	return debug
}

// Env returns the current env (debug or release).
func (c Config) Env() string {
	if c.IsDebug() {
		return "debug"
	}
	return "prod"
}

// loadEnv does some viper magic to make env vars work as expected.
// see comment below and https://github.com/spf13/viper/issues/188
func loadEnv() (*viper.Viper, error) {
	v := viper.New()
	// set default values in viper.
	// Viper needs to know if a key exists in order to override it.
	// https://github.com/spf13/viper/issues/188
	// we can safely ignore this error because we control the struct
	b, _ := yaml.Marshal(configDefaults)
	defaults := map[string]interface{}{}
	// we can safely ignore this error because we control the struct
	_ = yaml.Unmarshal(b, &defaults)
	v.SetTypeByDefaultValue(true)
	for key, val := range defaults {
		v.SetDefault(key, val)
	}
	// tell viper to overwrite env variables
	v.AutomaticEnv()
	v.SetEnvPrefix(ENVPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)
	return v, nil
}
