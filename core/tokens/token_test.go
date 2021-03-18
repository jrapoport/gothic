package tokens

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
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

func tokenConn(t *testing.T) *store.Connection {
	conn, _ := tconn.TempConn(t)
	mg := migration.NewMigrationWithIndexes("1",
		token.AccessToken{}, token.AccessTokenIndexes)
	err := conn.RunMigration(mg)
	require.NoError(t, err)
	return conn
}

func TestUseToken(t *testing.T) {
	t.Parallel()
	const testToken = "1234567890asdfghjkl="
	conn := tokenConn(t)
	tk := token.NewAccessToken(testToken, token.SingleUse, token.NoExpiration)
	err := conn.Create(tk).Error
	assert.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Usable())
	err = UseToken(conn, tk)
	assert.NoError(t, err)
	assert.Equal(t, 1, tk.Used)
	assert.NotNil(t, tk.UsedAt)
	assert.True(t, tk.DeletedAt.Valid)
	assert.False(t, tk.Usable())
	err = UseToken(conn, tk)
	assert.Error(t, err)
}
