package login

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	c   *config.Config
	jwt config.JWT
}

func TestLogin(t *testing.T) {
	ts := &LoginTestSuite{}
	suite.Run(t, ts)
}

func (ts *LoginTestSuite) SetupSuite() {
	ts.c = tconf.TempDB(ts.T())
	ts.c.UseInternal = true
	ts.jwt = ts.c.JWT
}

func testUser(t *testing.T, tx *store.Connection, p provider.Name, email, pw string) *user.User {
	u, err := users.CreateUser(tx, p, email, "", pw, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, u)
	now := time.Now()
	u.ConfirmedAt = &now
	u.Status = user.Active
	err = tx.Save(u).Error
	require.NoError(t, err)
	return u
}

func (ts *LoginTestSuite) TestLogin() {
	const (
		empty    = ""
		badPass  = "pass"
		testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	)
	email := func() string {
		return tutils.RandomEmail()
	}
	address := func() string {
		return tutils.RandomAddress()
	}
	type testCase struct {
		email string
		pw    string
	}
	tests := []testCase{
		{email(), empty},
		{email(), badPass},
		{email(), testPass},
	}
	ex := tests[len(tests)-1]
	ex.email = address()
	tests = append(tests, ex)
	conn := tconn.Conn(ts.T(), ts.c)
	err := conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			p := ts.c.Provider()
			u := testUser(ts.T(), tx, p, test.email, test.pw)
			u2, err := UserLogin(tx, p, test.email, test.pw)
			ts.NoError(err)
			ts.NotNil(u2)
			bt, err := tokens.GrantBearerToken(tx, ts.jwt, u2)
			ts.NoError(err)
			ts.NotNil(bt)
			ts.Equal(u.ID, bt.IssuedTo())
		}
		return nil
	})
	ts.NoError(err)
}

func (ts *LoginTestSuite) TestLogin_Error() {
	const (
		empty        = ""
		badEmail     = "bad"
		unknownEmail = "unknown@example.com"
		badPass      = "pass"
		testPass     = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	)
	type testCase struct {
		email string
		pw    string
	}
	errors1 := []testCase{
		{empty, empty},
		{empty, testPass},
	}
	var tests []testCase
	for _, e := range []string{empty, badEmail, unknownEmail} {
		for _, test := range errors1 {
			test.email = e
			tests = append(tests, test)
		}
	}
	p := ts.c.Provider()
	em := tutils.RandomEmail()
	conn := tconn.Conn(ts.T(), ts.c)
	u := testUser(ts.T(), conn, p, em, testPass)
	for _, test := range tests {
		_, err := UserLogin(conn, p, test.email, test.pw)
		ts.Error(err)
	}
	// bad password
	_, err := UserLogin(conn, p, em, badPass)
	ts.Error(err)
	// bad provider
	_, err = UserLogin(conn, "bad-provider", em, testPass)
	ts.Error(err)
	// disabled provider
	_, err = UserLogin(conn, provider.Google, em, testPass)
	ts.Error(err)
	// inactive user
	u.Status = user.Restricted
	err = conn.Save(u).Error
	ts.NoError(err)
	_, err = UserLogin(conn, p, em, testPass)
	ts.Error(err)
	// invalid user
	err = conn.Delete(u).Error
	ts.NoError(err)
	_, err = UserLogin(conn, p, em, testPass)
	ts.Error(err)
}

func (ts *LoginTestSuite) TestLogout() {
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	p := ts.c.Provider()
	em := tutils.RandomEmail()
	conn := tconn.Conn(ts.T(), ts.c)
	_ = testUser(ts.T(), conn, p, em, testPass)
	u, err := UserLogin(conn, p, em, testPass)
	ts.NoError(err)
	ts.Require().NotNil(u)
	bt, err := tokens.GrantBearerToken(conn, ts.jwt, u)
	ts.NoError(err)
	ts.Require().NotNil(bt)
	has, err := tokens.HasUsableRefreshToken(conn, bt.UserID)
	ts.NoError(err)
	ts.True(has)
	err = UserLogout(conn, bt.UserID)
	ts.NoError(err)
	has, err = tokens.HasUsableRefreshToken(conn, bt.UserID)
	ts.NoError(err)
	ts.False(has)
	// system user
	err = UserLogout(conn, uuid.Nil)
	ts.NoError(err)
	// unknown user
	err = UserLogout(conn, uuid.New())
	ts.NoError(err)
}
