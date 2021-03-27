package user_test

import (
	"net/http"
	"testing"

	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	rest_user "github.com/jrapoport/gothic/hosts/rest/user"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPass  = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	Endpoint  = rest_user.User
	PassRoute = rest_user.User + rest_user.Password
)

func userResponse(t *testing.T, res string) *rest.UserResponse {
	r := new(rest.UserResponse)
	err := json.Unmarshal([]byte(res), r)
	assert.NoError(t, err)
	return r
}

func testUser(t *testing.T, srv *rest.Host) (*user.User, string) {
	em := tutils.RandomEmail()
	un := utils.RandomUsername()
	data := types.Map{
		"color": "orange",
		"happy": "salad",
		"pick":  42.0,
	}
	ctx := context.Background()
	ctx.SetProvider(srv.Provider())
	u, err := srv.Signup(ctx, em, un, testPass, data)
	require.NoError(t, err)
	require.NotNil(t, u)
	u, err = srv.Login(ctx, em, testPass)
	require.NoError(t, err)
	require.NotNil(t, u)
	bt, err := srv.GrantBearerToken(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, bt)
	return u, bt.String()
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		rest_user.RegisterServer,
	}, false)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	// unmasked
	u, tok := testUser(t, srv)
	test := &rest.UserResponse{
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: u.Username,
		Data:     u.Data,
	}
	res, err := thttp.DoAuthRequest(t, web, http.MethodGet, Endpoint, tok, nil, nil)
	assert.NoError(t, err)
	ur := userResponse(t, res)
	assert.Equal(t, test, ur)
	// masked
	srv.Config().MaskEmails = true
	test.Email = utils.MaskEmail(u.Email)
	res, err = thttp.DoAuthRequest(t, web, http.MethodGet, Endpoint, tok, nil, nil)
	assert.NoError(t, err)
	ur = userResponse(t, res)
	assert.Equal(t, test, ur)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		rest_user.RegisterServer,
	}, false)
	req := &rest_user.Request{
		Username: "peaches",
		Data: types.Map{
			"foo":   "bar",
			"tasty": "salad",
		},
	}
	// invalid username
	srv.Config().Signup.AutoConfirm = true
	u, tok := testUser(t, srv)
	srv.Config().Validation.UsernameRegex = "0"
	_, err := thttp.DoAuthRequest(t, web, http.MethodPut, Endpoint, tok, nil, req)
	assert.Error(t, err)
	// bad req
	res, err := thttp.DoAuthRequest(t, web, http.MethodPut, Endpoint, tok, nil, []byte("\n"))
	assert.Error(t, err)
	// unmasked
	srv.Config().MaskEmails = false
	srv.Config().Validation.UsernameRegex = ""
	res, err = thttp.DoAuthRequest(t, web, http.MethodPut, Endpoint, tok, nil, req)
	assert.NoError(t, err)
	ur := userResponse(t, res)
	err = mergo.Map(&u.Data, req.Data)
	assert.NoError(t, err)
	test := &rest.UserResponse{
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: req.Username,
		Data:     u.Data,
	}
	assert.Equal(t, test, ur)
	// masked
	srv.Config().MaskEmails = true
	req = &rest_user.Request{
		Username: "mario",
		Data: types.Map{
			"quack":   99.0,
			"peaches": "happy",
		},
	}
	res, err = thttp.DoAuthRequest(t, web, http.MethodPut, Endpoint, tok, nil, req)
	assert.NoError(t, err)
	ur = userResponse(t, res)
	err = mergo.Map(&u.Data, req.Data)
	assert.NoError(t, err)
	test = &rest.UserResponse{
		Role:     u.Role.String(),
		Email:    utils.MaskEmail(u.Email),
		Username: req.Username,
		Data:     u.Data,
	}
	assert.Equal(t, test, ur)
}

func TestChangePassword(t *testing.T) {
	t.Parallel()
	const newPassword = "gj8#xtg#yrabxpnno!p5f3t8na!hd3?4jq7majxs"
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		rest_user.RegisterServer,
	}, false)
	req := rest_user.Request{
		NewPassword: newPassword,
	}
	// invalid password
	u, tok := testUser(t, srv)
	srv.Config().Validation.PasswordRegex = "0"
	_, err := thttp.DoAuthRequest(t, web, http.MethodPut, PassRoute, tok, nil, req)
	assert.Error(t, err)
	// wrong old password
	srv.Config().Validation.PasswordRegex = ""
	_, err = thttp.DoAuthRequest(t, web, http.MethodPut, PassRoute, tok, nil, req)
	assert.Error(t, err)
	// success
	req.Password = testPass
	req.NewPassword = newPassword
	res, err := thttp.DoAuthRequest(t, web, http.MethodPut, PassRoute, tok, nil, req)
	assert.NoError(t, err)
	_, claims := tsrv.UnmarshalTokenResponse(t, srv.Config().JWT, res)
	assert.Equal(t, u.ID.String(), claims.Subject)
	u, err = srv.GetUser(u.ID)
	assert.NoError(t, err)
	err = u.Authenticate(newPassword)
	assert.NoError(t, err)
}

func TestErrors(t *testing.T) {
	t.Parallel()
	srv, web, _ := tsrv.RESTHost(t, []rest.RegisterServer{
		rest_user.RegisterServer,
	}, false)
	j := srv.Config().JWT
	srv.Config().Signup.AutoConfirm = true
	testRoute := func(t *testing.T, method, route string, checkReq bool) {
		// not authorized
		_, err := thttp.DoRequest(t, web, method, route, nil, nil)
		assert.Error(t, err)
		// no user id
		tok := thttp.BadToken(t, j)
		_, err = thttp.DoAuthRequest(t, web, method, route, tok, nil, nil)
		assert.Error(t, err)
		// user not found
		tok = thttp.UserToken(t, j, false, true)
		_, err = thttp.DoAuthRequest(t, web, method, route, tok, nil, &rest.Request{})
		assert.Error(t, err)
		if checkReq {
			// invalid req
			tok = thttp.UserToken(t, j, false, true)
			_, err = thttp.DoAuthRequest(t, web, method, route, tok, nil, []byte("\n"))
			assert.Error(t, err)
		}
		// user banned
		u, tok := tcore.TestUser(t, srv.API, "", false)
		_, err = srv.BanUser(context.Background(), u.ID)
		require.NoError(t, err)
		_, err = thttp.DoAuthRequest(t, web, method, route, tok, nil, nil)
		assert.Error(t, err)
	}
	tests := []struct {
		m        string
		r        string
		checkReq bool
	}{
		{http.MethodGet, Endpoint, false},
		{http.MethodPut, Endpoint, true},
		{http.MethodPut, PassRoute, true},
	}
	for _, test := range tests {
		testRoute(t, test.m, test.r, test.checkReq)
	}
}
