package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountServer_Logout(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	u := testUser(t, srv)
	// not authorized
	ctx := context.Background()
	_, err := srv.Logout(ctx, nil)
	assert.Error(t, err)
	// bad token
	ctx = context.Background()
	ctx = testAuthCtx(t, srv, u)
	_, err = srv.Logout(ctx, nil)
	assert.NoError(t, err)
}
