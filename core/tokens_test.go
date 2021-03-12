package core

import (
	"testing"

	"github.com/jrapoport/gothic/models/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_GrantBearerToken(t *testing.T) {
	a := apiWithTempDB(t)
	u := testUser(t, a)
	// nil
	_, err := a.GrantBearerToken(nil, nil)
	assert.Error(t, err)
	// invalid user
	_, err = a.GrantBearerToken(nil, new(user.User))
	assert.Error(t, err)
	// not confirmed
	_, err = a.GrantBearerToken(nil, u)
	assert.NoError(t, err)
	u = confirmUser(t, a, u)
	bt, err := a.GrantBearerToken(nil, u)
	assert.NoError(t, err)
	require.NotNil(t, bt)
	assert.NotEmpty(t, bt.AccessToken.String())
	assert.Equal(t, u.ID, bt.UserID)
	require.NotNil(t, bt.RefreshToken)
	assert.NotEmpty(t, bt.RefreshToken.String())
	assert.Equal(t, u.ID, bt.RefreshToken.UserID)
}

func TestAPI_RefreshBearerToken(t *testing.T) {
	a := apiWithTempDB(t)
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	bt1, err := a.GrantBearerToken(nil, u)
	require.NoError(t, err)
	rt := bt1.RefreshToken
	require.NotNil(t, rt)
	bt2, err := a.RefreshBearerToken(nil, rt.Token)
	assert.NoError(t, err)
	assert.NotNil(t, bt2)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.NoError(t, err)
	banUser(t, a, u)
	rt = bt2.RefreshToken
	_, err = a.RefreshBearerToken(nil, rt.Token)
	assert.Error(t, err)
	err = a.conn.Delete(u).Error
	assert.NoError(t, err)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	err = a.conn.Delete(rt).Error
	assert.NoError(t, err)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
}
