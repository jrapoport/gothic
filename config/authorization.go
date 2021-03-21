package config

import (
	"net/url"
	"strings"

	"github.com/jrapoport/gothic/models/types/provider"
)

// Providers is a map of configured external providers.
type Providers map[provider.Name]Provider

// Authorization holds external provider configs
type Authorization struct {
	UseInternal bool      `json:"provider_internal" yaml:"provider_internal" mapstructure:"provider_internal"`
	RedirectURL string    `json:"provider_redirect_url" yaml:"provider_redirect_url" mapstructure:"provider_redirect_url"`
	Providers   Providers `json:"provider" yaml:"provider" mapstructure:"provider"`
	internal    provider.Name
}

// Provider name of the internal provider (if enabled).
func (a *Authorization) Provider() provider.Name {
	if !a.UseInternal {
		return provider.Unknown
	}
	return a.internal
}

func (a *Authorization) useInternalProvider(name string) {
	a.internal = provider.NormalizeName(name)
}

const hostToken = ":host"

func (a *Authorization) normalize(srv Service, host string) error {
	if a.UseInternal {
		a.useInternalProvider(srv.Name)
	}
	if a.RedirectURL != "" {
		_, err := url.Parse(a.RedirectURL)
		if err != nil {
			return err
		}
	}
	const urlFormat = "http://" + hostToken + "/account/auth/callback"
	for k, p := range a.Providers {
		if p.ClientKey == "" {
			delete(a.Providers, k)
			continue
		}
		callback := p.CallbackURL
		if callback == "" {
			callback = urlFormat
		}
		callback = FormatCallback(callback, host)
		_, err := url.Parse(callback)
		if err != nil {
			return err
		}
		p.CallbackURL = callback
		a.Providers[k] = p
	}
	return nil
}

// Provider holds external provider config.
type Provider struct {
	ClientKey   string   `json:"client_key" yaml:"client_key" mapstructure:"client_key"`
	Secret      string   `json:"secret"`
	CallbackURL string   `json:"callback_url" yaml:"callback_url" mapstructure:"callback_url"`
	Scopes      []string `json:"scopes"`
}

/*
// SAML is for SAML support.
type SAML struct {
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	MetadataURL string `json:"metadata_url" yaml:"metadata_url" mapstructure:"metadata_url"`
	APIBase     string `json:"api_base" yaml:"api_base" mapstructure:"api_base"`
	Cert 		string `json:"cert"`
	Key  		string `json:"key"`
}
*/

// FormatCallback formats the link URL replacing the ':host' with the RESTAddress().
func FormatCallback(callback, host string) string {
	return strings.ReplaceAll(callback, hostToken, host)
}
