package conf

import "errors"

// ExternalConfig holds SSO OAuth configs
type ExternalConfig struct {
	Bitbucket   OAuthProvider `json:"bitbucket"`
	Github      OAuthProvider `json:"github"`
	Gitlab      OAuthProvider `json:"gitlab"`
	Google      OAuthProvider `json:"google"`
	Facebook    OAuthProvider `json:"facebook"`
	Email       EmailProvider `json:"email"`
	Saml        SamlProvider  `json:"saml"`
	RedirectURL string        `json:"redirect_url"`
}

type EmailProvider struct {
	Disabled bool `json:"disabled"`
}

type SamlProvider struct {
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	MetadataURL string `json:"metadata_url" split_words:"true"`
	APIBase     string `json:"api_base" split_words:"true"`
	SigningCert string `json:"signing_cert" split_words:"true"`
	SigningKey  string `json:"signing_key" split_words:"true"`
}

// OAuthProvider holds all config related to external account providers.
type OAuthProvider struct {
	ClientID    string `json:"client_id" split_words:"true"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri" split_words:"true"`
	URL         string `json:"url"`
	Enabled     bool   `json:"enabled"`
}

func (o *OAuthProvider) Validate() error {
	if !o.Enabled {
		return errors.New("provider is not enabled")
	}
	if o.ClientID == "" {
		return errors.New("missing Oauth client ID")
	}
	if o.Secret == "" {
		return errors.New("missing Oauth secret")
	}
	if o.RedirectURI == "" {
		return errors.New("missing redirect URI")
	}
	return nil
}
