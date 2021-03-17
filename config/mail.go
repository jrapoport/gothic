package config

import (
	"net/url"
	"strings"
	"time"
)

// Mail config
type Mail struct {
	// SMTP is the smtp configuration.
	SMTP `yaml:",inline" mapstructure:",squash"`
	// Name is the name to use with outbound mail (default: ServiceName).
	Name string `json:"name"`
	// Link is the url to use with outbound mail (default: SiteURL).
	Link string `json:"link"`
	// Logo is the logo to use with outbound mail (default: SiteLogo).
	Logo string `json:"logo"`
	// From is the originating address to use with outbound email.
	From string `json:"from"`
	// Theme is the theme to use when formatting emails.
	// Options are "default", "flat", "" = default
	Theme         string `json:"theme"`
	MailTemplates `yaml:",inline" mapstructure:",squash"`
}

// SMTP config
type SMTP struct {
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty" mapstructure:",omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	// Authentication options are "plain", "login", "cram-md5", "" = plain
	Authentication string `json:"authentication"`
	// Encryption options are "none", "ssl", "tls", "" = none
	Encryption string        `json:"encryption"`
	KeepAlive  bool          `json:"keepalive"`
	Expiration time.Duration `json:"expiration"`
	SendLimit  time.Duration `json:"send_limit" yaml:"send_limit" mapstructure:"send_limit"`
	// SpamProtection enables smtp email account validation
	SpamProtection bool `json:"spam_protection" yaml:"spam_protection" mapstructure:"spam_protection"`
}

func (m *Mail) normalize(srv Service) error {
	if m.Name == "" {
		n := strings.ToLower(srv.Name)
		n = strings.Title(n)
		m.Name = n
	}
	if m.Link == "" {
		m.Link = srv.SiteURL
	}
	if m.Link != "" {
		u, err := url.Parse(m.Link)
		if err != nil {
			return err
		}
		if m.From == "" {
			m.From = mailDefaults.From
		}
		m.From = strings.Replace(m.From, ":name", m.Name, 1)
		m.From = strings.Replace(m.From, ":link_hostname", u.Hostname(), 1)
	}
	return m.MailTemplates.normalize(srv)
}

// CheckSendLimit returns ErrRateLimitExceeded if lest exceeds the send limit
func (m Mail) CheckSendLimit(last *time.Time) error {
	if last == nil {
		return nil
	}
	limit := m.SendLimit
	// FUTURE happened before now
	if last.Add(limit).After(time.Now().UTC()) {
		return ErrRateLimitExceeded
	}
	return nil
}
