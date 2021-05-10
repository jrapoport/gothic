// +build release

package core

import "github.com/jrapoport/gothic/store"

func (a *API) debugSignup(_ *store.Connection, _, _ string) bool {
	return false
}
