package tcore

import (
	"testing"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/require"
)

// TestUser creates a test user.
func TestUser(t *testing.T, a *core.API, pass string, admin bool) (*user.User, string) {
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	if pass == "" {
		pass = testPass
	}
	em := tutils.RandomEmail()
	ctx := context.Background()
	ctx.SetProvider(a.Provider())
	ctx.SetAdminID(user.SuperAdminID)
	u, err := a.Signup(ctx, em, "", pass, nil)
	require.NoError(t, err)
	require.NotNil(t, u)
	if admin {
		u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
		require.NoError(t, err)
		require.NotNil(t, u)
	}
	if !u.IsConfirmed() {
		return u, ""
	}
	bt, err := a.GrantBearerToken(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, bt)
	return u, bt.String()
}
