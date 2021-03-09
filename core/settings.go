package core

import "github.com/jrapoport/gothic/core/settings"

// Settings returns the current api settings.
func (a *API) Settings() settings.Settings {
	return settings.Current(a.config)
}
