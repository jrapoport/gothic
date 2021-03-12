package core

import (
	"github.com/jrapoport/gothic/test/tconf"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_Signup(t *testing.T) {
	tests := []struct {
		name string
		auto bool
	}{
		{"Confirm", false},
		{"AutoConfirm", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testSignup(t, test.auto)
		})
	}
}

func testSignup(t *testing.T, autoconfirm bool) {
	a := createAPI(t)
	if !autoconfirm {
		a.config, _ = tconf.MockSMTP(t, a.config)
		err := a.OpenMail()
		require.NoError(t, err)
	}
	a.config.Signup.AutoConfirm = autoconfirm
	const (
		empty       = ""
		badUsername = "!"
		badPass     = "pass"
		testPass    = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
		userRx      = "^[a-zA-Z0-9_]{2,255}$"
		passRx      = "^[a-zA-Z0-9[:punct:]]{8,40}$"
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
		unRx     string
		pw       string
		pwRx     string
		data     types.Map
	}
	tests := []testCase{
		{email(), empty, empty, empty, empty, nil},
		{email(), empty, empty, badPass, empty, nil},
		{email(), empty, empty, testPass, empty, nil},
		{email(), empty, empty, testPass, passRx, nil},
		{email(), badUsername, empty, empty, empty, nil},
		{email(), badUsername, empty, badPass, empty, nil},
		{email(), badUsername, empty, testPass, empty, nil},
		{email(), badUsername, empty, testPass, passRx, nil},
		{email(), username(), empty, empty, empty, nil},
		{email(), username(), empty, badPass, empty, nil},
		{email(), username(), empty, testPass, empty, nil},
		{email(), username(), empty, testPass, passRx, nil},
		{email(), username(), userRx, empty, empty, nil},
		{email(), username(), userRx, badPass, empty, nil},
		{email(), username(), userRx, testPass, empty, nil},
		{email(), username(), userRx, testPass, passRx, nil},
	}
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
	ctx := testContext(a)
	for _, test := range tests {
		a.config.Validation.UsernameRegex = test.unRx
		a.config.Validation.PasswordRegex = test.pwRx
		checkUser := func(u *user.User) {
			require.NotNil(t, u)
			e, err := validate.Email(test.email)
			require.NoError(t, err)
			assert.Equal(t, e, u.Email)
			assert.Equal(t, test.username, u.Username)
			assert.Equal(t, autoconfirm, u.IsConfirmed())
			if !u.IsConfirmed() {
				var ct token.ConfirmToken
				err = a.conn.First(&ct, "user_id = ?", u.ID).Error
				assert.NoError(t, err)
			}
			if test.data == nil {
				return
			}
			assert.Equal(t, test.data, u.Data)
		}
		u, err := a.Signup(ctx, test.email, test.username, test.pw, test.data)
		assert.NoError(t, err)
		checkUser(u)
		u, err = a.GetUserWithEmail(test.email)
		require.NoError(t, err)
		checkUser(u)
	}
}

func TestAPI_Signup_Error(t *testing.T) {
	const (
		empty        = ""
		badEmail     = "bad"
		badUsername  = "!"
		badPass      = "pass"
		badRx        = " "
		testUsername = "foobar"
		testEmail    = "quack@example.com"
		testAddress  = "Foo Bar <quack@example.com>"
		testPass     = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
		userRx       = "^[a-zA-Z0-9_]{2,255}$"
		passRx       = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	)
	var testData = types.Map{
		key.Color:   "#ffaa11",
		"full_name": "Foo Bar",
		"avatar":    "http://example.com/user/image.png",
	}
	type testCase struct {
		email    string
		username string
		unRx     string
		pw       string
		pwRx     string
	}
	errors1 := []testCase{
		{empty, empty, empty, empty, empty},
		{empty, empty, empty, empty, passRx},
		{empty, empty, empty, badPass, empty},
		{empty, empty, empty, badPass, passRx},
		{empty, empty, empty, testPass, empty},
		{empty, empty, empty, testPass, passRx},
	}
	var tests []testCase
	for _, e := range []string{empty, badEmail} {
		for _, u := range []string{empty, badUsername, testUsername} {
			for _, rx := range []string{empty, userRx} {
				for _, test := range errors1 {
					test.email = e
					test.username = u
					test.unRx = rx
					tests = append(tests, test)
				}
			}
		}
	}
	errors2 := []testCase{
		{testEmail, empty, empty, empty, passRx},
		{testEmail, empty, empty, badPass, passRx},
		{testEmail, empty, userRx, empty, empty},
		{testEmail, empty, userRx, empty, passRx},
		{testEmail, empty, userRx, badPass, empty},
		{testEmail, empty, userRx, badPass, passRx},
		{testEmail, empty, userRx, testPass, empty},
		{testEmail, empty, userRx, testPass, passRx},
		{testEmail, badUsername, empty, empty, passRx},
		{testEmail, badUsername, empty, badPass, passRx},
		{testEmail, badUsername, userRx, empty, empty},
		{testEmail, badUsername, userRx, empty, passRx},
		{testEmail, badUsername, userRx, badPass, empty},
		{testEmail, badUsername, userRx, badPass, passRx},
		{testEmail, badUsername, userRx, testPass, empty},
		{testEmail, badUsername, userRx, testPass, passRx},
		{testEmail, testUsername, empty, empty, passRx},
		{testEmail, testUsername, empty, badPass, passRx},
		{testEmail, testUsername, userRx, empty, passRx},
		{testEmail, testUsername, userRx, badPass, passRx},
		{testEmail, testUsername, badRx, testPass, passRx},
		{testEmail, testUsername, badRx, testPass, passRx},
		{testEmail, testUsername, badRx, testPass, passRx},
		{testEmail, testUsername, badRx, testPass, passRx},
		{testEmail, testUsername, empty, testPass, badRx},
		{testEmail, testUsername, empty, testPass, badRx},
		{testEmail, testUsername, userRx, testPass, badRx},
		{testEmail, testUsername, userRx, testPass, badRx},
	}
	for _, test := range errors2 {
		tests = append(tests, test)
		test.email = testAddress
		tests = append(tests, test)
	}
	a := createAPI(t)
	a.config.Signup.Username = true
	ctx := testContext(a)
	for _, data := range []types.Map{nil, {}, testData} {
		for _, test := range tests {
			a.config.Validation.UsernameRegex = test.unRx
			a.config.Validation.PasswordRegex = test.pwRx
			_, err := a.Signup(ctx, test.email, test.username, test.pw, data)
			assert.Error(t, err)
		}
	}
	// invalid context
	_, err := a.Signup(nil, testEmail, testUsername, testPass, nil)
	assert.Error(t, err)
	// test email taken
	a.config.Validation.UsernameRegex = empty
	a.config.Validation.PasswordRegex = empty
	_, err = a.Signup(ctx, testEmail, testUsername, testPass, nil)
	assert.NoError(t, err)
	_, err = a.Signup(ctx, testEmail, testUsername, testPass, nil)
	assert.Error(t, err)
	ctx.SetProvider("bad-provider")
	_, err = a.Signup(ctx, testEmail, testUsername, testPass, nil)
	assert.Error(t, err)
}

func TestAPI_Signup_Disabled(t *testing.T) {
	a := createAPI(t)
	a.config.Signup.Disabled = true
	_, err := a.Signup(context.Background(), "", "", "", nil)
	assert.Error(t, err)
}

func TestAPI_Signup_ReCaptcha(t *testing.T) {
	a := createAPI(t)
	ctx := context.Background()
	_, err := a.Signup(ctx, "", "", "", nil)
	assert.Error(t, err)
	a.config.Recaptcha.Key = validate.ReCaptchaDebugKey
	ctx.SetReCaptcha("bad")
	_, err = a.Signup(ctx, "", "", "", nil)
	assert.Error(t, err)
	ctx = context.Background()
	ctx.SetIPAddress(testIP)
	_, err = a.Signup(ctx, "", "", "", nil)
	assert.Error(t, err)
	ctx = context.Background()
	ctx.SetIPAddress(testIP)
	ctx.SetReCaptcha(validate.ReCaptchaDebugKey)
	_, err = a.Signup(ctx, "", "", "", nil)
	assert.Error(t, err)
	ctx.SetReCaptcha(validate.ReCaptchaDebugToken)
	_, err = a.Signup(ctx, "", "", "", nil)
	assert.Error(t, err)
	a.config.Validation.UsernameRegex = ""
	em := tutils.RandomEmail()
	_, err = a.Signup(ctx, em, "", testPass, nil)
	assert.NoError(t, err)
}

func TestAPI_Signup_SignupCode(t *testing.T) {
	const (
		testUsername = "foobar"
		testPass     = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	)
	email := func() string {
		return tutils.RandomEmail()
	}
	a := createAPI(t)
	ctx := testContext(a)
	sc, err := codes.CreateCode(a.conn, uuid.Nil, code.PIN, 1, true)
	require.NoError(t, err)
	require.NotNil(t, sc)
	// Signup code off
	a.config.Signup.Code = false
	// no code
	_, err = a.Signup(ctx, email(), testUsername, testPass, nil)
	assert.NoError(t, err)
	// ignore code
	ctx.SetCode("ignore")
	_, err = a.Signup(ctx, email(), testUsername, testPass, nil)
	assert.NoError(t, err)
	// Signup code on
	a.config.Signup.Code = true
	// no code
	ctx = testContext(a)
	_, err = a.Signup(ctx, email(), testUsername, testPass, nil)
	assert.Error(t, err)
	// bad code
	ctx.SetCode("bad")
	_, err = a.Signup(ctx, email(), testUsername, testPass, nil)
	assert.Error(t, err)
	// good code
	ctx.SetCode(sc.Code())
	_, err = a.Signup(ctx, email(), testUsername, testPass, nil)
	assert.NoError(t, err)
	// check code was used
	_, err = codes.GetUsableCode(a.conn, sc.Code())
	assert.Error(t, err)
	// can't re-use code
	_, err = a.Signup(ctx, email(), testUsername, testPass, nil)
	assert.Error(t, err)
}

func TestAPI_Signup_Username(t *testing.T) {
	const (
		noUsername   = ""
		badUsername  = "!"
		testUsername = "foobar"
		testPassword = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
		userRx       = "^[a-zA-Z0-9_]{2,255}$"
	)
	email := func() string {
		return tutils.RandomEmail()
	}
	a := createAPI(t)
	ctx := testContext(a)
	type testCase struct {
		username string
		req      bool
		def      bool
		rx       string
		Err      assert.ErrorAssertionFunc
		Empty    assert.ValueAssertionFunc
	}
	tests := []testCase{
		// not required
		{noUsername, false, false, "", assert.NoError, assert.Empty},
		// required
		{noUsername, true, false, "", assert.Error, assert.Empty},
		// bad username, not required
		{badUsername, false, false, userRx, assert.Error, assert.Empty},
		// bad username, required
		{badUsername, true, false, userRx, assert.Error, assert.Empty},
		// good username, not required
		{testUsername, false, false, userRx, assert.NoError, assert.NotEmpty},
		// good username, required
		{testUsername, true, false, userRx, assert.NoError, assert.NotEmpty},
		// default, not required
		{noUsername, false, true, "", assert.NoError, assert.NotEmpty},
		// default, required
		{noUsername, true, true, "", assert.NoError, assert.NotEmpty},
		// default, not required
		{noUsername, false, true, userRx, assert.NoError, assert.NotEmpty},
		// default, required
		{noUsername, true, true, userRx, assert.NoError, assert.NotEmpty},
	}
	for _, test := range tests {
		a.config.Signup.Username = test.req
		a.config.Signup.Default.Username = test.def
		a.config.Validation.UsernameRegex = test.rx
		u, err := a.Signup(ctx, email(), test.username, testPassword, nil)
		test.Err(t, err)
		if u != nil {
			test.Empty(t, u.Username)
		}
	}
}

func TestAPI_Signup_DefaultColor(t *testing.T) {
	const testPassword = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	email := func() string {
		return tutils.RandomEmail()
	}
	a := createAPI(t)
	ctx := testContext(a)
	a.config.Signup.Username = false
	a.config.Signup.Default.Color = false
	u, err := a.Signup(ctx, email(), "", testPassword, nil)
	require.NoError(t, err)
	require.NotNil(t, u)
	u, err = a.Signup(ctx, email(), "", testPassword, types.Map{})
	require.NoError(t, err)
	require.NotNil(t, u)
	require.NotNil(t, u.Data)
	_, ok := u.Data[key.Color]
	assert.False(t, ok)
	a.config.Signup.Default.Color = true
	u, err = a.Signup(ctx, email(), "", testPassword, nil)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.NotEmpty(t, u.Data)
	clr1, ok := u.Data[key.Color].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, clr1)
	u, err = a.Signup(ctx, email(), "", testPassword, u.Data)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.NotEmpty(t, u.Data)
	clr2, ok := u.Data[key.Color].(string)
	assert.True(t, ok)
	assert.Equal(t, clr1, clr2)
}

func TestAPI_Signup_Event(t *testing.T) {
	a := createAPI(t)
	testListen := func(mu *sync.RWMutex, lis events.Event, data *types.Map) {
		c := a.Listen(lis)
		go func() {
			for {
				msg, ok := <-c
				if !ok {
					break
				}
				mu.Lock()
				evt, ok := msg[key.Event].(events.Event)
				if ok && evt == events.Signup {
					*data = msg
				}
				mu.Unlock()
			}
		}()
	}
	testAddListener := func(mu *sync.RWMutex, lis events.Event, data *types.Map) {
		a.AddListener(lis, func(evt events.Event, msg types.Map) {
			mu.Lock()
			defer mu.Unlock()
			if evt == events.Signup {
				*data = msg
			}
		})
	}
	tests := []struct {
		name   string
		lis    events.Event
		listen listenerTestFunc
	}{
		{"Listen_Signup", events.Signup, testListen},
		{"Listen_All", events.All, testListen},
		{"AddListener_Signup", events.Signup, testAddListener},
		{"AddListener_All", events.All, testAddListener},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testSignupEvent(t, a, test.lis, test.listen)
		})
	}
}

func testSignupEvent(t *testing.T, a *API, lis events.Event, l listenerTestFunc) {
	const (
		testUsername = "foobar"
		testPass     = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
		testIP       = "127.0.0.1"
	)
	email := tutils.RandomEmail()
	var data types.Map
	var mu sync.RWMutex
	l(&mu, lis, &data)
	p := a.Provider()
	ctx := testContext(a)
	u, err := a.Signup(ctx, email, testUsername, testPass, nil)
	assert.NoError(t, err)
	assert.NotNil(t, t, u)
	assert.Equal(t, p, u.Provider)
	assert.Equal(t, email, u.Email)
	assert.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return data != nil
	}, 10*time.Second, 100*time.Millisecond)
	assert.Equal(t, events.Signup, data[key.Event].(events.Event))
	assert.Equal(t, testIP, data[key.IPAddress].(string))
	assert.Equal(t, p, data[key.Provider].(provider.Name))
	assert.Equal(t, email, data[key.Email].(string))
}
