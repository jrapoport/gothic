package core

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/stretchr/testify/assert"
)

func loginAPI(t *testing.T) *API {
	a := apiWithTempDB(t)
	a.config.Signup.Default.Username = false
	a.config.Signup.Default.Color = false
	a.config.Signup.Code = false
	a.config.Recaptcha.Key = ""
	a.mail.UseSpamProtection(false)
	return a
}

func TestAPI_Login(t *testing.T) {
	t.Parallel()
	const (
		empty          = ""
		badEmail       = "bad"
		unknownEmail   = "quack@example.com"
		unknownAddress = "Foo Bar <quack@example.com>"
		badPass        = "pass"
		testPass       = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	)
	a := loginAPI(t)
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	tests := []struct {
		email string
		pw    string
		Err   assert.ErrorAssertionFunc
	}{
		{empty, empty, assert.Error},
		{badEmail, empty, assert.Error},
		{unknownEmail, empty, assert.Error},
		{unknownAddress, empty, assert.Error},
		{empty, badPass, assert.Error},
		{badEmail, badPass, assert.Error},
		{unknownEmail, badPass, assert.Error},
		{unknownAddress, badPass, assert.Error},
		{empty, testPass, assert.Error},
		{badEmail, testPass, assert.Error},
		{unknownEmail, testPass, assert.Error},
		{unknownAddress, testPass, assert.Error},
		{u.Email, testPass, assert.NoError},
		{u.EmailAddress().String(), testPass, assert.NoError},
	}
	ctx := testContext(a)
	for _, test := range tests {
		_, err := a.Login(ctx, test.email, test.pw)
		test.Err(t, err)
	}
}

func TestAPI_Login_Disabled(t *testing.T) {
	t.Parallel()
	a := loginAPI(t)
	ctx := context.Background()
	_, err := a.Login(ctx, "", "")
	assert.Error(t, err)
}

func TestAPI_Login_ReCaptcha(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	a := loginAPI(t)
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	a.config.Security.Recaptcha.Login = true
	ctx.SetProvider(a.Provider())
	_, err := a.Login(ctx, u.Email, "")
	assert.Error(t, err)
	a.config.Recaptcha.Key = validate.ReCaptchaDebugKey
	ctx.SetReCaptcha("bad")
	_, err = a.Login(ctx, u.Email, "")
	assert.Error(t, err)
	ctx = testContext(a)
	_, err = a.Login(ctx, u.Email, "")
	assert.Error(t, err)
	ctx.SetReCaptcha("bad")
	_, err = a.Login(ctx, u.Email, "")
	assert.Error(t, err)
	ctx.SetReCaptcha(validate.ReCaptchaDebugToken)
	_, err = a.Login(ctx, u.Email, "")
	assert.Error(t, err)
	_, err = a.Login(ctx, u.Email, testPass)
	assert.NoError(t, err)
}

func TestAPI_Login_Event(t *testing.T) {
	t.Parallel()
	a := loginAPI(t)
	testListen := func(mu *sync.RWMutex, lis events.Event, data *types.Map) {
		c := a.Listen(lis)
		go func() {
			for {
				msg, open := <-c
				if !open {
					break
				}
				evt, ok := msg[key.Event].(events.Event)
				if !ok {
					continue
				}
				mu.Lock()
				switch evt {
				case events.Login:
					*data = msg
				case events.Logout:
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
			if evt == events.Login {
				*data = msg
			}
			if evt == events.Logout {
				*data = msg
			}
		})
	}
	evts := []events.Event{events.Login, events.Logout}
	all := []events.Event{events.All}
	tests := []struct {
		name   string
		lis    []events.Event
		listen listenerTestFunc
	}{
		{"Listen_Login", evts, testListen},
		{"Listen_All", all, testListen},
		{"AddListener_Login", evts, testAddListener},
		{"AddListener_All", all, testAddListener},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testLoginEvent(t, a, test.lis, test.listen)
		})
	}
}

func testLoginEvent(t *testing.T, a *API, lis []events.Event, l listenerTestFunc) {
	var data types.Map
	var mu sync.RWMutex
	for _, li := range lis {
		l(&mu, li, &data)
	}
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	ctx := testContext(a)
	u2, err := a.Login(ctx, u.Email, testPass)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u.ID, u2.ID)
	assert.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return data != nil
	}, 100*time.Millisecond, 10*time.Millisecond)
	assert.Equal(t, events.Login, data[key.Event].(events.Event))
	assert.Equal(t, testIP, data[key.IPAddress].(string))
	assert.Equal(t, u.ID, data[key.UserID].(uuid.UUID))
	bt, err := a.GrantBearerToken(ctx, u2)
	assert.NoError(t, err)
	assert.NotNil(t, bt)
	data = nil
	err = a.Logout(ctx, bt.UserID)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return data != nil
	}, 100*time.Millisecond, 10*time.Millisecond)
	assert.Equal(t, events.Logout, data[key.Event].(events.Event))
	assert.Equal(t, testIP, data[key.IPAddress].(string))
	assert.Equal(t, u.ID, data[key.UserID].(uuid.UUID))
}
