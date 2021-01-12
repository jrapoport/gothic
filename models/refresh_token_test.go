package models

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/storage/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RefreshTokenTestSuite struct {
	suite.Suite
	db *storage.Connection
}

func (ts *RefreshTokenTestSuite) SetupTest() {
	err := storage.TruncateAll(ts.db)
	assert.NoError(ts.T(), err)
}

func TestRefreshToken(t *testing.T) {
	globalConfig, err := conf.LoadGlobal(modelsTestConfig)
	require.NoError(t, err)

	conn, err := test.SetupDBConnection(t, globalConfig)
	require.NoError(t, err)

	ts := &RefreshTokenTestSuite{
		db: conn,
	}

	suite.Run(t, ts)
}

func (ts *RefreshTokenTestSuite) TestGrantAuthenticatedUser() {
	u := ts.createUser()
	r, err := GrantAuthenticatedUser(ts.db, u)
	require.NoError(ts.T(), err)

	require.NotEmpty(ts.T(), r.Token)
	require.Equal(ts.T(), u.ID, r.UserID)
}

func (ts *RefreshTokenTestSuite) TestGrantRefreshTokenSwap() {
	u := ts.createUser()
	r, err := GrantAuthenticatedUser(ts.db, u)
	require.NoError(ts.T(), err)

	s, err := GrantRefreshTokenSwap(ts.db, u, r)
	require.NoError(ts.T(), err)

	_, nr, err := FindUserWithRefreshToken(ts.db, r.Token)
	require.NoError(ts.T(), err)

	require.Equal(ts.T(), r.ID, nr.ID)
	require.True(ts.T(), nr.Revoked, "expected old token to be revoked")

	require.NotEqual(ts.T(), r.ID, s.ID)
	require.Equal(ts.T(), u.ID, s.UserID)
}

func (ts *RefreshTokenTestSuite) TestLogout() {
	u := ts.createUser()
	r, err := GrantAuthenticatedUser(ts.db, u)
	require.NoError(ts.T(), err)

	require.NoError(ts.T(), Logout(ts.db, u.ID))
	u, r, err = FindUserWithRefreshToken(ts.db, r.Token)
	require.Errorf(ts.T(), err, "expected error when there are no refresh tokens to authenticate. user: %v token: %v", u, r)

	require.True(ts.T(), IsNotFoundError(err), "expected NotFoundError")
}

func (ts *RefreshTokenTestSuite) createUser() *User {
	return ts.createUserWithEmail("testusergothic.com")
}

func (ts *RefreshTokenTestSuite) createUserWithEmail(email string) *User {
	user, err := NewUser(email, "secret", "test", nil)
	require.NoError(ts.T(), err)

	err = ts.db.Create(user).Error
	require.NoError(ts.T(), err)

	return user
}
