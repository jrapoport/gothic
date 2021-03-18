package login

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	conn *store.Connection
	c    *config.Config
	jwt  config.JWT
}

func TestLogin(t *testing.T) {
	ts := &LoginTestSuite{}
	suite.Run(t, ts)
}

func (ts *LoginTestSuite) SetupSuite() {
	ts.conn, ts.c = tconn.TempConn(ts.T())
	ts.c.UseInternal = true
	ts.jwt = ts.c.JWT
}

func (ts *LoginTestSuite) testUser(p provider.Name, email, pw string) *user.User {
	u, err := users.CreateUser(ts.conn, p, email, "", pw, nil, nil)
	ts.Require().NoError(err)
	ts.Require().NotNil(u)
	now := time.Now()
	u.ConfirmedAt = &now
	u.Status = user.Active
	err = ts.conn.Save(u).Error
	ts.Require().NoError(err)
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
	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]
		test.email = address()
		tests = append(tests, test)
	}
	for _, test := range tests {
		p := ts.c.Provider()
		u := ts.testUser(p, test.email, test.pw)
		u2, err := UserLogin(ts.conn, p, test.email, test.pw)
		ts.NoError(err)
		ts.NotNil(u2)
		bt, err := tokens.GrantBearerToken(ts.conn, ts.jwt, u2)
		ts.NoError(err)
		ts.NotNil(bt)
		ts.Equal(u.ID, bt.IssuedTo())
	}
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
	u := ts.testUser(p, em, testPass)
	for _, test := range tests {
		_, err := UserLogin(ts.conn, p, test.email, test.pw)
		ts.Error(err)
	}
	// bad password
	_, err := UserLogin(ts.conn, p, em, badPass)
	ts.Error(err)
	// bad provider
	_, err = UserLogin(ts.conn, "bad-provider", em, testPass)
	ts.Error(err)
	// disabled provider
	_, err = UserLogin(ts.conn, provider.Google, em, testPass)
	ts.Error(err)
	// inactive user
	u.Status = user.Restricted
	err = ts.conn.Save(u).Error
	ts.NoError(err)
	_, err = UserLogin(ts.conn, p, em, testPass)
	ts.Error(err)
	// invalid user
	err = ts.conn.Delete(u).Error
	ts.NoError(err)
	_, err = UserLogin(ts.conn, p, em, testPass)
	ts.Error(err)
}

func (ts *LoginTestSuite) TestLogout() {
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	p := ts.c.Provider()
	em := tutils.RandomEmail()
	_ = ts.testUser(p, em, testPass)
	u, err := UserLogin(ts.conn, p, em, testPass)
	ts.NoError(err)
	ts.Require().NotNil(u)
	bt, err := tokens.GrantBearerToken(ts.conn, ts.jwt, u)
	ts.NoError(err)
	ts.Require().NotNil(bt)
	has, err := tokens.HasUsableRefreshToken(ts.conn, bt.UserID)
	ts.NoError(err)
	ts.True(has)
	err = UserLogout(ts.conn, bt.UserID)
	ts.NoError(err)
	has, err = tokens.HasUsableRefreshToken(ts.conn, bt.UserID)
	ts.NoError(err)
	ts.False(has)
	// system user
	err = UserLogout(ts.conn, uuid.Nil)
	ts.NoError(err)
	// unknown user
	err = UserLogout(ts.conn, uuid.New())
	ts.NoError(err)
}
