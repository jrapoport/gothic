package api

import (
	"net/http"
)

type ProviderSettings struct {
	Bitbucket bool `json:"bitbucket"`
	GitHub    bool `json:"github"`
	GitLab    bool `json:"gitlab"`
	Google    bool `json:"google"`
	Facebook  bool `json:"facebook"`
	Email     bool `json:"email"`
	SAML      bool `json:"saml"`
}

type ProviderLabels struct {
	SAML string `json:"saml,omitempty"`
}

type Settings struct {
	ExternalProviders ProviderSettings `json:"external"`
	ExternalLabels    ProviderLabels   `json:"external_labels"`
	DisableSignup     bool             `json:"disable_signup"`
	Autoconfirm       bool             `json:"autoconfirm"`
}

func (a *API) handleSettings(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, a.Settings())
}

func (a *API) Settings() *Settings {
	config := a.config
	external := a.config.External
	return &Settings{
		ExternalProviders: ProviderSettings{
			Bitbucket: external.Bitbucket.Enabled,
			GitHub:    external.Github.Enabled,
			GitLab:    external.Gitlab.Enabled,
			Google:    external.Google.Enabled,
			Facebook:  external.Facebook.Enabled,
			Email:     !external.Email.Disabled,
			SAML:      external.Saml.Enabled,
		},
		ExternalLabels: ProviderLabels{
			SAML: external.Saml.Name,
		},
		DisableSignup: config.DisableSignup,
		Autoconfirm:   config.Mailer.Autoconfirm,
	}
}
