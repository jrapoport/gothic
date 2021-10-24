package core

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_CreateUser(t *testing.T) {
	t.Parallel()
	em := tutils.RandomEmail()
	un := utils.RandomUsername()
	a := apiWithTempDB(t)
	ctx := rootContext(a)
	a.config.Signup.Username = true
	a.config.Signup.Default.Username = false
	_, err := a.CreateUser(nil, "", "", "", nil, false)
	assert.Error(t, err)
	ctx.SetAdminID(uuid.New())
	_, err = a.CreateUser(ctx, "", "", "", nil, false)
	assert.Error(t, err)
	_, err = a.CreateUser(ctx, "@", "", "", nil, false)
	assert.Error(t, err)
	_, err = a.CreateUser(ctx, em, "", "", nil, false)
	assert.Error(t, err)
	_, err = a.CreateUser(ctx, em, un, "", nil, false)
	assert.Error(t, err)
	adm := testUser(t, a)
	adm = promoteUser(t, a, adm)
	ctx.SetAdminID(adm.ID)
	_, err = a.CreateUser(ctx, em, un, "", nil, false)
	assert.Error(t, err)
	_, err = a.GrantBearerToken(ctx, adm)
	require.NoError(t, err)
	_, err = a.CreateUser(ctx, em, un, "", nil, false)
	assert.NoError(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleAdmin)
	require.NoError(t, err)
	_, err = a.CreateUser(ctx, em, un, "", nil, false)
	assert.Error(t, err)
	_, err = a.CreateUser(ctx, em, un, "", nil, true)
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleSuper)
	require.NoError(t, err)
	_, err = a.CreateUser(ctx, em, un, "", nil, true)
	assert.Error(t, err)
	ctx.SetProvider(a.Provider())
	ctx.SetUserID(adm.ID)
	em = tutils.RandomEmail()
	u, err := a.CreateUser(ctx, em, un, "", nil, true)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.Valid())
	assert.True(t, u.IsConfirmed())
	assert.True(t, u.IsAdmin())
}

func TestAPI_ChangeUserRole(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	ctx := rootContext(a)
	u := testUser(t, a)
	_, err := a.ChangeRole(nil, uuid.Nil, user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, uuid.New(), user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, uuid.New(), user.RoleAdmin)
	assert.Error(t, err)
	// promote
	ctx = rootContext(a)
	u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.NoError(t, err)
	assert.True(t, u.IsAdmin())
	// re-promote
	u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.NoError(t, err)
	assert.True(t, u.IsAdmin())
	// super-promote
	_, err = a.ChangeRole(ctx, u.ID, user.RoleSuper)
	assert.Error(t, err)
	// fail promote
	u.Role = user.RoleSuper
	err = a.conn.Save(u).Error
	require.NoError(t, err)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.Error(t, err)
	// demote
	u.Role = user.RoleAdmin
	err = a.conn.Save(u).Error
	require.NoError(t, err)
	u, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.NoError(t, err)
	assert.False(t, u.IsAdmin())
	banUser(t, a, u)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.NoError(t, err)
	// only super admin can promote to admin
	adm := testUser(t, a)
	adm.Role = user.RoleAdmin
	err = a.conn.Save(adm).Error
	require.NoError(t, err)
	ctx.SetAdminID(adm.ID)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.Error(t, err)
	// only super admin can demote to user
	u.Role = user.RoleAdmin
	err = a.conn.Save(u).Error
	require.NoError(t, err)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.Error(t, err)
}

func TestAPI_PromoteUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	ctx := rootContext(a)
	a.config.Signup.Username = true
	a.config.Signup.Default.Username = false
	_, err := a.PromoteUser(nil, uuid.Nil)
	assert.Error(t, err)
	ctx.SetAdminID(uuid.New())
	_, err = a.PromoteUser(ctx, uuid.Nil)
	assert.Error(t, err)
	_, err = a.PromoteUser(ctx, uuid.New())
	assert.Error(t, err)
	adm := testUser(t, a)
	adm = confirmUser(t, a, adm)
	ctx.SetAdminID(adm.ID)
	_, err = a.PromoteUser(ctx, uuid.New())
	assert.Error(t, err)
	_, err = a.GrantBearerToken(ctx, adm)
	require.NoError(t, err)
	_, err = a.PromoteUser(ctx, uuid.New())
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleAdmin)
	require.NoError(t, err)
	_, err = a.PromoteUser(ctx, uuid.New())
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleSuper)
	require.NoError(t, err)
	_, err = a.PromoteUser(ctx, uuid.New())
	assert.Error(t, err)
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	assert.False(t, u.IsAdmin())
	ctx = rootContext(a)
	pu, err := a.PromoteUser(ctx, u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, pu)
	assert.Equal(t, u.ID, pu.ID)
	assert.True(t, pu.Valid())
	assert.True(t, pu.IsAdmin())
	_, err = a.PromoteUser(ctx, u.ID)
	assert.NoError(t, err)
}

func TestAPI_UpdateUserMetadata(t *testing.T) {
	const (
		testKey   = "test-key"
		testValue = "test-value"
	)
	var testMeta = types.Map{
		testKey:       testValue,
		key.IPAddress: "foo",
	}
	t.Parallel()
	a := apiWithTempDB(t)
	ctx := rootContext(a)
	a.config.Signup.Username = true
	a.config.Signup.Default.Username = false
	_, err := a.UpdateUserMetadata(nil, uuid.Nil, nil)
	assert.Error(t, err)
	ctx.SetAdminID(uuid.New())
	_, err = a.UpdateUserMetadata(ctx, uuid.Nil, nil)
	assert.Error(t, err)
	_, err = a.UpdateUserMetadata(ctx, uuid.New(), nil)
	assert.Error(t, err)
	_, err = a.UpdateUserMetadata(ctx, uuid.New(), testMeta)
	assert.Error(t, err)
	adm := testUser(t, a)
	adm = confirmUser(t, a, adm)
	ctx.SetAdminID(adm.ID)
	_, err = a.UpdateUserMetadata(ctx, uuid.New(), testMeta)
	assert.Error(t, err)
	_, err = a.GrantBearerToken(ctx, adm)
	require.NoError(t, err)
	_, err = a.UpdateUserMetadata(ctx, uuid.New(), testMeta)
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleAdmin)
	require.NoError(t, err)
	_, err = a.UpdateUserMetadata(ctx, uuid.New(), testMeta)
	assert.Error(t, err)
	u := testUser(t, a)
	assert.False(t, u.IsAdmin())
	assert.Equal(t, testIP, u.Metadata[key.IPAddress])
	u, err = a.UpdateUserMetadata(ctx, u.ID, testMeta)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, testValue, u.Metadata[testKey])
	assert.Equal(t, testIP, u.Metadata[key.IPAddress])
	err = users.ChangeRole(a.conn, u, user.RoleAdmin)
	require.NoError(t, err)
	_, err = a.UpdateUserMetadata(ctx, u.ID, testMeta)
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleSuper)
	require.NoError(t, err)
	u, err = a.UpdateUserMetadata(ctx, u.ID, testMeta)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, testValue, u.Metadata[testKey])
	assert.Equal(t, testIP, u.Metadata[key.IPAddress])
	err = a.conn.DB.Migrator().DropColumn(new(user.User), "Metadata")
	require.NoError(t, err)
	_, err = a.UpdateUserMetadata(ctx, u.ID, testMeta)
	assert.Error(t, err)

}

func TestAPI_DeleteUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	ctx := rootContext(a)
	a.config.Signup.Username = true
	a.config.Signup.Default.Username = false
	err := a.DeleteUser(nil, uuid.Nil, false)
	assert.Error(t, err)
	ctx.SetAdminID(uuid.New())
	err = a.DeleteUser(ctx, uuid.Nil, false)
	assert.Error(t, err)
	err = a.DeleteUser(ctx, uuid.New(), false)
	assert.Error(t, err)
	adm := testUser(t, a)
	adm = confirmUser(t, a, adm)
	ctx.SetAdminID(adm.ID)
	err = a.DeleteUser(ctx, uuid.New(), false)
	assert.Error(t, err)
	_, err = a.GrantBearerToken(ctx, adm)
	require.NoError(t, err)
	err = a.DeleteUser(ctx, uuid.New(), false)
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleAdmin)
	require.NoError(t, err)
	err = a.DeleteUser(ctx, uuid.New(), false)
	assert.Error(t, err)
	err = users.ChangeRole(a.conn, adm, user.RoleSuper)
	require.NoError(t, err)
	err = a.DeleteUser(ctx, uuid.New(), false)
	assert.Error(t, err)
	u := testUser(t, a)
	assert.False(t, u.IsAdmin())
	err = a.DeleteUser(ctx, u.ID, false)
	assert.NoError(t, err)
	err = a.DeleteUser(ctx, u.ID, false)
	assert.Error(t, err)
	u, err = a.GetUser(u.ID)
	assert.Error(t, err)
	assert.Nil(t, u)
	u = testUser(t, a)
	assert.False(t, u.IsAdmin())
	err = a.DeleteUser(ctx, u.ID, true)
	assert.NoError(t, err)
	// make an admin
	adm = testUser(t, a)
	adm = confirmUser(t, a, adm)
	adm.Role = user.RoleAdmin
	err = a.conn.Save(adm).Error
	require.NoError(t, err)
	ctx.SetAdminID(adm.ID)
	err = a.DeleteUser(ctx, adm.ID, true)
	assert.Error(t, err)
}
