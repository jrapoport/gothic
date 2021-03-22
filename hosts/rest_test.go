package hosts

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account"
	"github.com/jrapoport/gothic/hosts/rest/account/login"
	"github.com/jrapoport/gothic/hosts/rest/user"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRESTHost(t *testing.T) {
	t.Parallel()
	c := tconf.TempDB(t)
	c.Signup.AutoConfirm = true
	c.Security.MaskEmails = false
	a, err := core.NewAPI(c)
	require.NoError(t, err)
	h := NewRESTHost(a, "127.0.0.1:0")
	require.NotNil(t, h)
	err = h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.Online()
	}, 1*time.Second, 10*time.Millisecond)
	// create a test user
	const pass = "1234567890asdfghjkl"
	test, _ := tcore.TestUser(t, a, pass, false)
	// unauthenticated call
	loginURI := func() string {
		require.NotEmpty(t, h.Address())
		return "http://" + h.Address() + account.Account + login.Login
	}
	b, err := json.Marshal(&login.Request{
		Email:    test.Email,
		Password: pass,
	})
	require.NoError(t, err)
	res, err := http.Post(loginURI(), rest.JSONContent, bytes.NewBuffer(b))
	require.NoError(t, err)
	b, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	err = res.Body.Close()
	require.NoError(t, err)
	ur, claims := tsrv.UnmarshalUserResponse(t, c.JWT, string(b))
	assert.EqualValues(t, tokens.Bearer, ur.Token.Type)
	assert.Equal(t, test.ID.String(), claims.Subject)
	assert.Equal(t, test.Email, ur.Email)
	// authenticated call (error)
	getUserURI := func() string {
		require.NotEmpty(t, h.Address())
		return "http://" + h.Address() + user.Endpoint
	}
	req := thttp.Request(t, http.MethodGet, getUserURI(), ur.Token.Access, nil, nil)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	b, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	err = res.Body.Close()
	require.NoError(t, err)
	ur, _ = tsrv.UnmarshalUserResponse(t, c.JWT, string(b))
	assert.Equal(t, test.Email, ur.Email)
	err = h.Shutdown()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return !h.Online()
	}, 1*time.Second, 10*time.Millisecond)
}
