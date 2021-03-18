package config

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/jrapoport/gothic/store/types/provider"
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
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	name = strings.ToLower(name)
	name = reg.ReplaceAllString(name, "")
	a.internal = provider.Name(name)
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

// FormatCallback formats the link URL replacing the ':host' with the RESTAddress().
func FormatCallback(callback, host string) string {
	return strings.ReplaceAll(callback, hostToken, host)
}
