package core

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_GetUser(t *testing.T) {
	a := apiWithTempDB(t)
	u1 := testUser(t, a)
	u2, err := a.GetUser(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Username, u2.Username)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = a.GetUser(user.SystemID)
	assert.Error(t, err)
	_, err = a.GetUser(uuid.New())
	assert.Error(t, err)
}

func TestAPI_GetAuthenticatedUser(t *testing.T) {
	a := apiWithTempDB(t)
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	_, err := a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	bt, err := tokens.GrantBearerToken(a.conn, a.config.JWT, u)
	require.NoError(t, err)
	au, err := a.GetAuthenticatedUser(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
	err = a.conn.Delete(u).Error
	assert.NoError(t, err)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	err = a.conn.Delete(bt.RefreshToken).Error
	assert.NoError(t, err)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
}

func TestAPI_GetUserWithEmail(t *testing.T) {
	a := apiWithTempDB(t)
	u1 := testUser(t, a)
	u2, err := a.GetUserWithEmail(u1.EmailAddress().String())
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Username, u2.Username)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = a.GetUserWithEmail("does-not-exist@example.com")
	assert.Error(t, err)
}

func TestAPI_ChangePassword(t *testing.T) {
	var newPass = utils.SecureToken()
	a := apiWithTempDB(t)
	ctx := testContext(a)
	a.config.Validation.PasswordRegex = ""
	u := testUser(t, a)
	_, err := a.ChangePassword(nil, uuid.Nil, testPass, newPass)
	assert.Error(t, err)
	_, err = a.ChangePassword(ctx, u.ID, "", newPass)
	assert.Error(t, err)
	_, err = a.ChangePassword(ctx, u.ID, testPass, newPass)
	assert.NoError(t, err)
	banUser(t, a, u)
	_, err = a.ChangePassword(ctx, u.ID, newPass, testPass)
	assert.Error(t, err)
	a.config.Validation.PasswordRegex = "!"
	_, err = a.ChangePassword(ctx, u.ID, "", "")
	assert.Error(t, err)
}

func TestAPI_CmdChangeUserRole(t *testing.T) {
	a := apiWithTempDB(t)
	ctx := testContext(a)
	u := testUser(t, a)
	_, err := a.ChangeRole(nil, uuid.Nil, user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, uuid.New(), user.RoleAdmin)
	assert.Error(t, err)
	// promote
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
	// sneak promote
	u.Role = user.RoleSuper
	err = a.conn.Save(u).Error
	assert.NoError(t, err)
	u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.NoError(t, err)
	assert.True(t, u.IsAdmin())
	// demote
	u, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.NoError(t, err)
	assert.False(t, u.IsAdmin())
	banUser(t, a, u)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.Error(t, err)
}

func TestAPI_ConfirmUser(t *testing.T) {
	a := apiWithTempDB(t)
	u := testUser(t, a)
	ctx := testContext(a)
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	// bad token
	u, err = a.ConfirmUser(nil, "")
	assert.Error(t, err)
	// good token
	u, err = a.ConfirmUser(ctx, ct.String())
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	ct, err = tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	_, err = a.ConfirmUser(ctx, ct.String())
	assert.NoError(t, err)

}

func TestAPI_ConfirmPassword(t *testing.T) {
	const (
		empty   = ""
		badPass = "pass"
		passRx  = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	)
	a := apiWithTempDB(t)
	u := testUser(t, a)
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	// password validation
	a.config.Validation.PasswordRegex = passRx
	ctx := testContext(a)
	// bad token
	u, err = a.ConfirmPasswordChange(nil, "", testPass)
	assert.Error(t, err)
	_, err = a.ConfirmPasswordChange(nil, ct.String(), empty)
	assert.Error(t, err)
	_, err = a.ConfirmPasswordChange(ctx, ct.String(), badPass)
	assert.Error(t, err)
	u, err = a.ConfirmPasswordChange(ctx, ct.String(), testPass)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	err = u.Authenticate(testPass)
	assert.NoError(t, err)
	// no password validation
	a.config.Validation.PasswordRegex = empty
	ct, err = tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	u, err = a.ConfirmPasswordChange(ctx, ct.String(), empty)
	assert.NoError(t, err)
	err = u.Authenticate(empty)
	assert.NoError(t, err)
	// can't reuse token
	_, err = a.ConfirmPasswordChange(ctx, ct.String(), empty)
	assert.Error(t, err)
}

func TestAPI_ConfirmEmail(t *testing.T) {
	const (
		empty    = ""
		badEmail = "@"
	)
	var testEmail = tutils.RandomEmail()
	a := apiWithTempDB(t)
	ctx := testContext(a)
	u := testUser(t, a)
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	// bad token
	_, err = a.ConfirmChangeEmail(nil, "", testEmail)
	assert.Error(t, err)
	_, err = a.ConfirmChangeEmail(ctx, ct.String(), empty)
	assert.Error(t, err)
	_, err = a.ConfirmChangeEmail(ctx, ct.String(), badEmail)
	assert.Error(t, err)
	u, err = a.ConfirmChangeEmail(ctx, ct.String(), testEmail)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	// can't reuse token
	_, err = a.ConfirmChangeEmail(ctx, ct.String(), testEmail)
	assert.Error(t, err)
}

func TestAPI_UpdateUser(t *testing.T) {
	var testName = "peaches"
	data := types.Map{
		"foo":   "bar",
		"tasty": "salad",
	}
	a := apiWithTempDB(t)
	u := testUser(t, a)
	confirmUser(t, a, u)
	ctx := testContext(a)
	// system user
	_, err := a.UpdateUser(nil, uuid.Nil, nil, nil)
	assert.Error(t, err)
	// user not found
	_, err = a.UpdateUser(ctx, uuid.New(), nil, nil)
	assert.Error(t, err)
	a.config.Validation.UsernameRegex = "0"
	_, err = a.UpdateUser(ctx, u.ID, &testName, nil)
	assert.Error(t, err)
	a.config.Validation.UsernameRegex = ""
	u, err = a.UpdateUser(ctx, u.ID, &testName, nil)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, testName, u.Username)
	var fooName = "foo"
	u, err = a.UpdateUser(ctx, u.ID, &fooName, data)
	assert.NoError(t, err)
	assert.Equal(t, fooName, u.Username)
	assert.EqualValues(t, data, u.Data)
}

func TestAPI_BanUser(t *testing.T) {
	a := apiWithTempDB(t)
	u := testUser(t, a)
	assert.False(t, u.IsBanned())
	// no user id
	_, err := a.BanUser(nil, uuid.Nil)
	assert.Error(t, err)
	// bad user id
	_, err = a.BanUser(nil, uuid.New())
	assert.Error(t, err)
	// ban user
	u, err = a.BanUser(nil, u.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsBanned())
	// "re" ban user
	u, err = a.BanUser(nil, u.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsBanned())
}

func TestAPI_DeleteUser(t *testing.T) {
	a := apiWithTempDB(t)
	u := testUser(t, a)
	assert.True(t, u.Valid())
	assert.False(t, u.DeletedAt.Valid)
	// no user id
	err := a.DeleteUser(nil, uuid.Nil)
	assert.Error(t, err)
	// bad user id
	err = a.DeleteUser(nil, uuid.New())
	assert.Error(t, err)
	// delete user
	err = a.DeleteUser(nil, u.ID)
	assert.NoError(t, err)
	_, err = a.GetUser(u.ID)
	assert.Error(t, err)
}
