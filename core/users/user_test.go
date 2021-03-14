package users

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUser(t *testing.T, conn *store.Connection, p provider.Name) *user.User {
	email := tutils.RandomEmail()
	username := utils.RandomUsername()
	data := types.Map{
		"helo":   "world",
		"foobar": 13.37,
	}
	meta := types.Map{
		key.IPAddress: "127.0.0.1",
	}
	u, err := createUser(conn, p, email, username, "password", data, meta)
	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, username, u.Username)
	assert.Equal(t, data, u.Data)
	assert.Equal(t, meta, u.Metadata)
	return u
}

func banUser(t *testing.T, conn *store.Connection, u *user.User) {
	u.Status = user.Banned
	err := conn.Save(u).Error
	require.NoError(t, err)
}

func TestGetUserWithID(t *testing.T) {
	conn, c := tconn.TempConn(t)
	u1 := testUser(t, conn, c.Provider())
	u2, err := GetUser(conn, u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.ID, u2.ID)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = GetUser(conn, user.SystemID)
	assert.Error(t, err)
	_, err = GetUser(conn, uuid.New())
	assert.Error(t, err)
}

func TestGetActiveUserWithID(t *testing.T) {
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	_, err := GetActiveUser(conn, u.ID)
	assert.Error(t, err)
	err = ConfirmUser(conn, u, time.Now())
	require.NoError(t, err)
	tu, err := GetActiveUser(conn, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, tu.ID)
	banUser(t, conn, u)
	_, err = GetActiveUser(conn, u.ID)
	assert.Error(t, err)
	err = conn.Delete(u).Error
	require.NoError(t, err)
	_, err = GetActiveUser(conn, u.ID)
	assert.Error(t, err)
}

func TestGetAuthenticatedUserWithID(t *testing.T) {
	conn, c := tconn.TempConn(t)
	test := testUser(t, conn, c.Provider())
	_, err := GetAuthenticatedUser(conn, uuid.Nil)
	assert.Error(t, err)
	_, err = GetAuthenticatedUser(conn, test.ID)
	assert.Error(t, err)
	_, err = tokens.GrantBearerToken(conn, c.JWT, test)
	require.NoError(t, err)
	u, err := GetAuthenticatedUser(conn, test.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, test.ID, u.ID)
	err = ConfirmUser(conn, test, time.Now())
	require.NoError(t, err)
	u, err = GetAuthenticatedUser(conn, test.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, test.ID, u.ID)
	u.Status = user.Locked
	err = conn.Save(u).Error
	require.NoError(t, err)
	_, err = GetAuthenticatedUser(conn, test.ID)
	assert.Error(t, err)
}

func TestIsEmailTaken(t *testing.T) {
	var email = tutils.RandomEmail()
	var uname = utils.RandomUsername()
	testIsTakenFunc(t, email, uname, email, IsEmailTaken)
}

func TestIsUsernameTaken(t *testing.T) {
	var email = tutils.RandomEmail()
	var uname = utils.RandomUsername()
	testIsTakenFunc(t, email, uname, uname, IsUsernameTaken)
}

type isTakenFunc func(tx *store.Connection, val string) (bool, error)

func testIsTakenFunc(t *testing.T, email, uname, taken string, takenFn isTakenFunc) {
	var pw = []byte("")
	conn, c := tconn.TempConn(t)
	found, err := takenFn(conn, taken)
	assert.NoError(t, err)
	assert.False(t, found)
	p := c.Provider()
	u := user.NewUser(p, user.RoleUser, email, uname, pw, nil, nil)
	err = conn.Create(u).Error
	assert.NoError(t, err)
	found, err = takenFn(conn, taken)
	assert.NoError(t, err)
	assert.True(t, found)
	found, err = takenFn(conn, "")
	assert.Error(t, err)
	assert.False(t, found)
}

func TestGetUserWithEmail(t *testing.T) {
	conn, c := tconn.TempConn(t)
	u1 := testUser(t, conn, c.Provider())
	u2, err := GetUserWithEmail(conn, u1.EmailAddress().String())
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.ID, u2.ID)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = GetUserWithEmail(conn, "invalid email address")
	assert.Error(t, err)
	_, err = GetUserWithEmail(conn, "does-not-exist@exmaple.com")
	assert.Error(t, err)
}

func TestHasUserWithEmail(t *testing.T) {
	conn, c := tconn.TempConn(t)
	u1 := testUser(t, conn, c.Provider())
	u2, err := HasUserWithEmail(conn, u1.EmailAddress().String())
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.ID, u2.ID)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = HasUserWithEmail(conn, "invalid email address")
	assert.Error(t, err)
	u2, err = HasUserWithEmail(conn, "does-not-exist@exmaple.com")
	assert.NoError(t, err)
	assert.Nil(t, u2)
}

func TestRandomUsername(t *testing.T) {
	conn, c := tconn.TempConn(t)
	uname, err := RandomUsername(conn, false)
	assert.NoError(t, err)
	used := []string{uname}
	assert.NotEmpty(t, uname)
	// fill users
	for i := 0; i < 100; i++ {
		em := tutils.RandomEmail()
		un := utils.RandomUsername()
		used = append(used, un)
		p := c.Provider()
		u := user.NewUser(p, user.RoleUser, em, un, []byte(""), nil, nil)
		err = conn.Create(u).Error
		assert.NoError(t, err)
	}
	uname, err = RandomUsername(conn, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, uname)
	assert.NotContains(t, used, uname)
}
