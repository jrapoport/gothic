package tokens

import (
	"testing"

	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantBearerToken(t *testing.T) {
	conn, c := tconn.TempConn(t)
	// system user
	_, err := GrantBearerToken(conn, c.JWT, nil)
	assert.Error(t, err)
	_, err = GrantBearerToken(conn, c.JWT, new(user.User))
	assert.Error(t, err)
	u := testUser(t, conn, c)
	bt, err := GrantBearerToken(conn, c.JWT, u)
	assert.NoError(t, err)
	require.NotNil(t, bt)
	assert.True(t, bt.Usable())
	assert.NotEmpty(t, bt.AccessToken.String())
	assert.Equal(t, u.ID, bt.UserID)
	require.NotNil(t, bt.RefreshToken)
	assert.NotEmpty(t, bt.RefreshToken.String())
	assert.Equal(t, u.ID, bt.RefreshToken.UserID)
}

func TestRefreshBearerToken(t *testing.T) {
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c)
	bt, err := GrantBearerToken(conn, c.JWT, u)
	assert.NoError(t, err)
	require.NotNil(t, bt)
	assert.True(t, bt.Usable())
	bt, err = RefreshBearerToken(conn, c.JWT, u, bt.RefreshToken.String())
	assert.NoError(t, err)
	require.NotNil(t, bt)
	assert.True(t, bt.Usable())
	assert.NotEmpty(t, bt.AccessToken.String())
	assert.Equal(t, u.ID, bt.UserID)
	require.NotNil(t, bt.RefreshToken)
	assert.NotEmpty(t, bt.RefreshToken.String())
	assert.Equal(t, u.ID, bt.RefreshToken.UserID)
	// mismatched user
	u = testUser(t, conn, c)
	_, err = RefreshBearerToken(conn, c.JWT, u, bt.RefreshToken.String())
	assert.Error(t, err)
}
