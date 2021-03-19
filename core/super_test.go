package core

import (
	"testing"

	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/stretchr/testify/assert"
)

func TestAPI_CreateSuperAdmin(t *testing.T) {
	u := user.NewSuperAdmin("")
	a := unloadedAPI(t)
	has, err := a.conn.Has(&u, key.ID+" = ?", user.SuperAdminID)
	assert.NoError(t, err)
	assert.False(t, has)
	err = a.CreateSuperAdmin()
	assert.NoError(t, err)
	has, err = a.conn.Has(&u, key.ID+" = ?", user.SuperAdminID)
	assert.NoError(t, err)
	assert.True(t, has)
	pw := a.config.RootPassword
	a.config.RootPassword = "bad"
	err = a.CreateSuperAdmin()
	assert.Error(t, err)
	a.config.RootPassword = pw
	a.config.Name = "bad"
	err = a.CreateSuperAdmin()
	assert.Error(t, err)
}

func TestAPI_GetSuperAdmin(t *testing.T) {
	a := unloadedAPI(t)
	_, err := a.GetSuperAdmin("")
	assert.Error(t, err)
	err = a.CreateSuperAdmin()
	assert.NoError(t, err)
	su, err := a.GetSuperAdmin("")
	assert.NoError(t, err)
	assert.Equal(t, user.SuperAdminID, su.ID)
	_, err = a.GetSuperAdmin("bad")
	assert.Error(t, err)
	a.config.RootPassword = "bad"
	_, err = a.GetSuperAdmin("")
	assert.Error(t, err)
	a.config.Name = "bad"
	_, err = a.GetSuperAdmin("")
	assert.Error(t, err)
}
