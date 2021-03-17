package users

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/providers"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/suite"
)

type CreateUserTestSuite struct {
	suite.Suite
	conn *store.Connection
	c    *config.Config
}

func TestCreateUser(t *testing.T) {
	ts := &CreateUserTestSuite{}
	suite.Run(t, ts)
}

func (ts *CreateUserTestSuite) SetupSuite() {
	ts.conn, ts.c = tconn.TempConn(ts.T())
	ts.c.UseInternal = true
	err := providers.LoadProviders(ts.c)
	ts.NoError(err)
}

func (ts *CreateUserTestSuite) TestCreateUser() {
	const (
		empty       = ""
		badUsername = "!"
		badPass     = "pass"
		testPass    = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	)
	var testData = types.Map{
		key.Color:   "#ffaa11",
		"full_name": "Foo Bar",
		"avatar":    "http://example.com/user/image.png",
	}
	username := func() string {
		return utils.RandomUsername()
	}
	email := func() string {
		return tutils.RandomEmail()
	}
	address := func() string {
		return tutils.RandomAddress()
	}
	type testCase struct {
		email    string
		username string
		pw       string
		data     types.Map
	}
	tests := []testCase{
		{email(), empty, empty, nil},
		{email(), empty, badPass, nil},
		{email(), empty, testPass, nil},
		{email(), empty, testPass, nil},
		{email(), badUsername, empty, nil},
		{email(), badUsername, badPass, nil},
		{email(), badUsername, testPass, nil},
		{email(), badUsername, testPass, nil},
		{email(), username(), empty, nil},
		{email(), username(), badPass, nil},
		{email(), username(), testPass, nil},
		{email(), username(), testPass, nil},
		{email(), username(), empty, nil},
		{email(), username(), badPass, nil},
		{email(), username(), testPass, nil},
		{email(), username(), testPass, nil},
	}
	p := ts.c.Provider()
	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]
		for _, data := range []types.Map{{}, testData} {
			test.email = email()
			test.data = data
			tests = append(tests, test)
		}
		for _, data := range []types.Map{nil, {}, testData} {
			test.email = address()
			test.data = data
			tests = append(tests, test)
		}
	}
	err := ts.conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			checkUser := func(u *user.User) {
				ts.NotNil(u)
				e, err := validate.Email(test.email)
				ts.Require().NoError(err)
				ts.Equal(e, u.Email)
				ts.Equal(test.username, u.Username)
				ts.False(u.IsConfirmed())
				ts.Equal(test.data, u.Data)
				ts.Equal(test.data, u.Metadata)

			}
			u, err := CreateUser(tx, p, test.email, test.username, test.pw, test.data, test.data)
			ts.NoError(err)
			checkUser(u)
			u, err = GetUserWithEmail(tx, test.email)
			ts.Require().NoError(err)
			checkUser(u)
		}
		return nil
	})
	ts.Require().NoError(err)
}

func (ts *CreateUserTestSuite) TestCreateUser_Error() {
	const (
		empty        = ""
		badEmail     = "bad"
		badUsername  = "!"
		badPass      = "pass"
		testUsername = "foobar"
		testEmail    = "quack@example.com"
		testPass     = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	)
	var testData = types.Map{
		key.Color:   "#ffaa11",
		"full_name": "Foo Bar",
		"avatar":    "http://example.com/user/image.png",
	}
	type testCase struct {
		email    string
		username string
		pw       string
	}
	errors1 := []testCase{
		{empty, empty, empty},
		{empty, empty, badPass},
		{empty, empty, testPass},
	}
	p := ts.c.Provider()
	var tests []testCase
	for _, e := range []string{empty, badEmail} {
		for _, u := range []string{empty, badUsername, testUsername} {
			for _, test := range errors1 {
				test.email = e
				test.username = u
				tests = append(tests, test)
			}
		}
	}
	for _, data := range []types.Map{nil, {}, testData} {
		for _, test := range tests {
			_, err := CreateUser(ts.conn, p, test.email, test.username, test.pw, data, data)
			ts.Error(err)
		}
	}
	// duplicate email
	_, err := CreateUser(ts.conn, p, testEmail, testUsername, testPass, nil, nil)
	ts.NoError(err)
	// no provider
	p = provider.Unknown
	_, err = CreateUser(ts.conn, p, testEmail, testUsername, testPass, nil, nil)
	ts.Error(err)
	// bad provider
	p = "bad-provider"
	_, err = CreateUser(ts.conn, p, testEmail, testUsername, testPass, nil, nil)
	ts.Error(err)
	// external provider
	p = provider.Google
	_, err = CreateUser(ts.conn, p, testEmail, testUsername, testPass, nil, nil)
	ts.Error(err)
}
