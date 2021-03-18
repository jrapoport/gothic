package settings

import (
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/health"
	"github.com/jrapoport/gothic/models/types/provider"
)

// Settings is the api settings response
type Settings struct {
	health.Health
	Signup `json:"signup"`
	Mail   `json:"mail"`
}

// Signup settings
type Signup struct {
	Disabled    bool     `json:"disabled,omitempty"`
	AutoConfirm bool     `json:"autoconfirm,omitempty"`
	Provider    Provider `json:"provider,omitempty"`
}

// Provider settings
type Provider = struct {
	Internal string                 `json:"internal,omitempty"`
	External map[provider.Name]bool `json:"external,omitempty"`
}

// Mail settings
type Mail struct {
	Disabled       bool   `json:"disabled,omitempty"`
	Host           string `json:"host,omitempty"`
	Port           int    `json:"port,omitempty"`
	Authentication string `json:"authentication,omitempty"`
	Encryption     string `json:"encryption,omitempty"`
}

// Current returns the currently configured settings.
func Current(c *config.Config) Settings {
	m := Mail{
		Disabled: true,
	}
	if c.Mail.Host != "" {
		m = Mail{
			Host:           c.Mail.Host,
			Port:           c.Mail.Port,
			Authentication: c.Mail.Authentication,
			Encryption:     c.Mail.Encryption,
		}
	}
	p := Provider{
		Internal: c.Provider().String(),
		External: map[provider.Name]bool{},
	}
	for name := range c.Providers {
		p.External[name] = true
	}
	return Settings{
		health.Check(c),
		Signup{
			Disabled:    c.Signup.Disabled,
			AutoConfirm: c.Signup.AutoConfirm,
			Provider:    p,
		},
		m,
	}
}
