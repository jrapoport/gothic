package users

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	const (
		name1 = "peaches"
		name2 = "foobar"
	)
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	username := u.Username
	_, err := Update(conn, u, nil, types.Map{
		key.FirstName: name1,
	})
	assert.Error(t, err)
	err = ConfirmUser(conn, u, time.Now())
	require.NoError(t, err)
	ok, err := Update(conn, u, nil, types.Map{
		key.FirstName: name1,
	})
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, username, u.Username)
	assert.Equal(t, name1, u.Data[key.FirstName])
	username = name1
	ok, err = Update(conn, u, &username, nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, username, u.Username)
	username = name2
	ok, err = Update(conn, u, &username, types.Map{
		key.FirstName: name2,
		key.LastName:  name1,
	})
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, username, u.Username)
	assert.Equal(t, name2, u.Data[key.FirstName])
	assert.Equal(t, name1, u.Data[key.LastName])
	ok, err = Update(conn, u, nil, nil)
	assert.NoError(t, err)
	assert.False(t, ok)
	banUser(t, conn, u)
	ok, err = Update(conn, u, &username, nil)
	assert.Error(t, err)
	assert.False(t, ok)
	ok, err = Update(conn, nil, &username, nil)
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestConfirmUser(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	now := time.Now()
	err := ConfirmUser(conn, u, now)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	err = ConfirmUser(conn, u, time.Now())
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	assert.Equal(t, now.String(), u.ConfirmedAt.String())
	err = ConfirmUser(conn, nil, now)
	assert.Error(t, err)
	banUser(t, conn, u)
	err = ConfirmUser(conn, u, now)
	assert.Error(t, err)
	u = testUser(t, conn, c.Provider())
	err = ConfirmUser(conn, u, time.Time{})
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
}

func TestConfirmIfNeeded(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(conn, u.ID, 0)
	assert.NoError(t, err)
	// nil user
	_, err = ConfirmIfNeeded(conn, ct, nil)
	assert.Error(t, err)
	// invalid user
	_, err = ConfirmIfNeeded(conn, ct, new(user.User))
	assert.Error(t, err)
	// success
	confirmed, err := ConfirmIfNeeded(conn, ct, u)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	assert.Equal(t, confirmed, u.IsConfirmed())
	has, err := ct.HasToken(conn)
	assert.NoError(t, err)
	assert.False(t, has)
	// cannot reuse token
	_, err = ConfirmIfNeeded(conn, ct, u)
	assert.Error(t, err)
	// already confirmed
	ct, err = tokens.GrantConfirmToken(conn, u.ID, 0)
	assert.NoError(t, err)
	confirmed, err = ConfirmIfNeeded(conn, ct, u)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	assert.NotEqual(t, confirmed, u.IsConfirmed())
	ct.Token = "bad"
	_, err = ConfirmIfNeeded(conn, ct, u)
	assert.Error(t, err)
}

func TestChangeRole(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	err := ConfirmUser(conn, u, time.Now())
	require.NoError(t, err)
	err = ChangeRole(conn, u, user.RoleAdmin)
	assert.NoError(t, err)
	err = ChangeRole(conn, u, user.RoleAdmin)
	assert.NoError(t, err)
	err = ChangeRole(conn, u, user.InvalidRole)
	assert.Error(t, err)
	err = ChangeRole(conn, nil, user.RoleAdmin)
	assert.Error(t, err)
	banUser(t, conn, u)
	err = ChangeRole(conn, u, user.RoleAdmin)
	assert.Error(t, err)
}

func TestChangeEmail(t *testing.T) {
	t.Parallel()
	var newEmail = tutils.RandomEmail()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	assert.NotEqual(t, newEmail, u.Email)
	err := ChangeEmail(conn, u, newEmail)
	assert.Error(t, err)
	err = ConfirmUser(conn, u, time.Now())
	require.NoError(t, err)
	err = ChangeEmail(conn, u, newEmail)
	assert.NoError(t, err)
	assert.Equal(t, newEmail, u.Email)
	err = ChangeEmail(conn, u, "@")
	assert.Error(t, err)
	err = ChangeEmail(conn, u, newEmail)
	assert.NoError(t, err)
	banUser(t, conn, u)
	err = ChangeEmail(conn, u, newEmail)
	assert.Error(t, err)
	err = ChangeEmail(conn, nil, newEmail)
	assert.Error(t, err)
}

func TestChangePassword(t *testing.T) {
	t.Parallel()
	var newPassword = utils.SecureToken()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	err := ChangePassword(conn, u, newPassword)
	assert.NoError(t, err)
	err = u.Authenticate(newPassword)
	assert.NoError(t, err)
	u.Provider = provider.Google
	err = ChangePassword(conn, u, newPassword)
	assert.Error(t, err)
	banUser(t, conn, u)
	err = ChangePassword(conn, u, newPassword)
	assert.Error(t, err)
	err = ChangePassword(conn, nil, newPassword)
	assert.Error(t, err)
}

func TestLockUser(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	require.False(t, u.IsLocked())
	err := LockUser(conn, nil)
	assert.NoError(t, err)
	err = LockUser(conn, u)
	assert.NoError(t, err)
	require.True(t, u.IsLocked())
}

func TestBanUser(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	require.False(t, u.IsBanned())
	err := BanUser(conn, nil)
	assert.NoError(t, err)
	err = LockUser(conn, u)
	assert.NoError(t, err)
	require.True(t, u.IsLocked())
	assert.False(t, u.IsBanned())
	err = BanUser(conn, u)
	assert.NoError(t, err)
	err = LockUser(conn, u)
	require.NoError(t, err)
	assert.True(t, u.IsLocked())
	assert.True(t, u.IsBanned())
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	require.False(t, u.DeletedAt.Valid)
	// success
	err := DeleteUser(conn, u, false)
	assert.NoError(t, err)
	require.True(t, u.DeletedAt.Valid)
	_, err = GetUser(conn, u.ID)
	assert.Error(t, err)
	// ignored
	err = DeleteUser(conn, user.NewSystemUser(), false)
	assert.NoError(t, err)
	// ignored
	err = DeleteUser(conn, nil, false)
	assert.NoError(t, err)
	// ignore banned
	u = testUser(t, conn, c.Provider())
	require.False(t, u.DeletedAt.Valid)
	err = BanUser(conn, u)
	require.NoError(t, err)
	err = DeleteUser(conn, u, false)
	assert.NoError(t, err)
	require.False(t, u.DeletedAt.Valid)
	u, err = GetUser(conn, u.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsBanned())
	taken, err := IsEmailTaken(conn, u.Email)
	require.NoError(t, err)
	assert.True(t, taken)
	// hard delete
	u = testUser(t, conn, c.Provider())
	require.False(t, u.DeletedAt.Valid)
	err = DeleteUser(conn, u, true)
	assert.NoError(t, err)
	var count int64
	err = conn.Unscoped().
		Model(user.User{}).
		Where("id = ?", u.ID).
		Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
