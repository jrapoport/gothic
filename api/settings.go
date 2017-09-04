package api

import "net/http"

type ExternalProviderSettings struct {
	BitBucket bool `json:"bitbucket"`
	GitHub    bool `json:"github"`
	GitLab    bool `json:"gitlab"`
	Google    bool `json:"google"`
}

type Settings struct {
	ExternalProviders ExternalProviderSettings `json:"external"`
	SignupEnabled     bool                     `json:"signup_enabled"`
	Autoconfirm       bool                     `json:"autoconfirm"`
}

func (a *API) Settings(w http.ResponseWriter, r *http.Request) error {
	config := a.getConfig(r.Context())

	return sendJSON(w, http.StatusOK, &Settings{
		ExternalProviders: ExternalProviderSettings{
			BitBucket: config.External.Bitbucket.Enabled,
			GitHub:    config.External.Github.Enabled,
			GitLab:    config.External.Gitlab.Enabled,
			Google:    config.External.Google.Enabled,
		},
		SignupEnabled: config.SignupEnabled,
		Autoconfirm:   config.Mailer.Autoconfirm,
	})
}
