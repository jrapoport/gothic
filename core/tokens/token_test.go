package tokens

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/require"
)

func testUser(t *testing.T, conn *store.Connection, c *config.Config) *user.User {
	p := c.Provider()
	em := tutils.RandomEmail()
	r := user.RoleUser
	u := user.NewUser(p, r, em, "", []byte(""), nil, nil)
	now := time.Now()
	u.ConfirmedAt = &now
	u.Status = user.Active
	err := conn.Create(u).Error
	require.NoError(t, err)
	return u
}
