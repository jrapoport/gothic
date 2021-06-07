package core

import (
	"sync"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/data-dog/go-sqlmock.v2"
)

const (
	testIP   = "127.0.0.1"
	testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
)

type listenerTestFunc func(mu *sync.RWMutex, evt events.Event, data *types.Map)

func createAPI(t *testing.T) *API {
	a := apiWithTempDB(t)
	a.config.Signup.Default.Username = false
	a.config.Signup.Default.Color = false
	a.config.Signup.Code = false
	a.config.Signup.Username = false
	a.config.Recaptcha.Key = ""
	return a
}

func mockAPI(t *testing.T) (*API, sqlmock.Sqlmock) {
	c := tconf.TempDB(t)
	a, err := NewAPI(c)
	require.NoError(t, err)
	require.NotNil(t, a)
	conn, mock := tconn.MockConn(t)
	a.conn = conn
	defer func() {
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}()
	return a, mock
}

// apiWithTempDB creates a new API with a temp db for tests.
func apiWithTempDB(t *testing.T) *API {
	return configuredAPI(t, tconf.TempDB(t))
}

// configuredAPI creates a new API for tests with config.
func configuredAPI(t *testing.T, c *config.Config) *API {
	c.Signup.AutoConfirm = true
	c.Mail.SpamProtection = false
	c.Mail.KeepAlive = false
	a, err := NewAPI(c)
	require.NoError(t, err)
	require.NotNil(t, a)
	t.Cleanup(func() {
		err = a.Shutdown()
		require.NoError(t, err)
	})
	return a
}

func unloadedAPI(t *testing.T) *API {
	c := tconf.TempDB(t)
	conn, err := store.Dial(c, nil)
	require.NoError(t, err)
	a := &API{
		config: c,
		conn:   conn,
		log:    c.Log(),
	}
	return a
}

func testContext(a *API) context.Context {
	ctx := context.Background()
	ctx.SetProvider(a.Provider())
	ctx.SetIPAddress(testIP)
	return ctx
}

func rootContext(a *API) context.Context {
	ctx := testContext(a)
	ctx.SetAdminID(user.SuperAdminID)
	return ctx
}

func testUser(t *testing.T, a *API) *user.User {
	p := a.Provider()
	email := tutils.RandomEmail()
	un := utils.RandomUsername()
	u, err := users.CreateUser(a.conn, p, email, un, testPass, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, u)
	return u
}

func confirmUser(t *testing.T, a *API, u *user.User) *user.User {
	now := time.Now()
	u.ConfirmedAt = &now
	u.Status = user.Active
	err := a.conn.Save(u).Error
	require.NoError(t, err)
	return u
}

func banUser(t *testing.T, a *API, u *user.User) *user.User {
	u.Status = user.Banned
	err := a.conn.Save(u).Error
	require.NoError(t, err)
	return u
}

func promoteUser(t *testing.T, a *API, u *user.User) *user.User {
	u = confirmUser(t, a, u)
	u.Role = user.RoleAdmin
	err := a.conn.Save(u).Error
	require.NoError(t, err)
	return u
}
