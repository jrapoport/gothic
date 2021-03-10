package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"

func userResponse(t *testing.T, res string) *rest.UserResponse {
	r := new(rest.UserResponse)
	err := json.Unmarshal([]byte(res), r)
	assert.NoError(t, err)
	return r
}

func testUser(t *testing.T, srv *userServer) (*user.User, string) {
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

func TestUserServer_GetUser(t *testing.T) {
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	getUser := func(tok string) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodGet, Endpoint, tok, nil, nil)
		r, err := rest.ParseClaims(r, srv.Config().JWT, tok)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		srv.GetUser(w, r)
		return w
	}
	// unmasked
	u, tok := testUser(t, srv)
	test := &rest.UserResponse{
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: u.Username,
		Data:     u.Data,
	}
	res := getUser(tok)
	assert.Equal(t, http.StatusOK, res.Code)
	ur := userResponse(t, res.Body.String())
	assert.Equal(t, test, ur)
	// masked
	srv.Config().MaskEmails = true
	test.Email = utils.MaskEmail(u.Email)
	res = getUser(tok)
	assert.Equal(t, http.StatusOK, res.Code)
	ur = userResponse(t, res.Body.String())
	assert.Equal(t, test, ur)
}

func TestUserServer_UpdateUser(t *testing.T) {
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	updateUser := func(tok string, body interface{}) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodPut, Endpoint, tok, nil, body)
		r, err := rest.ParseClaims(r, srv.Config().JWT, tok)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		srv.UpdateUser(w, r)
		return w
	}
	req := &Request{
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
	res := updateUser(tok, req)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// bad req
	res = updateUser(tok, []byte("\n"))
	assert.NotEqual(t, http.StatusOK, res.Code)
	// unmasked
	srv.Config().MaskEmails = false
	srv.Config().Validation.UsernameRegex = ""
	res = updateUser(tok, req)
	assert.Equal(t, http.StatusOK, res.Code)
	ur := userResponse(t, res.Body.String())
	err := mergo.Map(&u.Data, req.Data)
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
	req = &Request{
		Username: "mario",
		Data: types.Map{
			"quack":   99.0,
			"peaches": "happy",
		},
	}
	res = updateUser(tok, req)
	assert.Equal(t, http.StatusOK, res.Code)
	ur = userResponse(t, res.Body.String())
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

func TestUserServer_ChangePassword(t *testing.T) {
	const newPassword = "gj8#xtg#yrabxpnno!p5f3t8na!hd3?4jq7majxs"
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	changePassword := func(tok string, body interface{}) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodPut, Password, tok, nil, body)
		r, err := rest.ParseClaims(r, srv.Config().JWT, tok)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		srv.ChangePassword(w, r)
		return w
	}
	req := Request{
		Password: newPassword,
	}
	// invalid password
	u, tok := testUser(t, srv)
	srv.Config().Validation.PasswordRegex = "0"
	res := changePassword(tok, req)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// wrong old password
	srv.Config().Validation.PasswordRegex = ""
	res = changePassword(tok, req)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// success
	req.OldPassword = testPass
	res = changePassword(tok, req)
	assert.Equal(t, http.StatusOK, res.Code)
	_, claims := tsrv.UnmarshalTokenResponse(t, srv.Config().JWT, res.Body.String())
	assert.Equal(t, u.ID.String(), claims.Subject)
	u, err := srv.API.GetUser(u.ID)
	assert.NoError(t, err)
	err = u.Authenticate(newPassword)
	assert.NoError(t, err)
}

func TestRequestErrors(t *testing.T) {
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	j := srv.Config().JWT
	srv.Config().Signup.AutoConfirm = true
	testRoute := func(t *testing.T, h http.HandlerFunc, method, route string, checkReq bool) {
		// not authorized
		r := thttp.Request(t, method, route, "", nil, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.NotEqual(t, http.StatusOK, w.Code)
		// no user id
		tok := thttp.BadToken(t, j)
		r = thttp.Request(t, method, route, tok, nil, nil)
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.NotEqual(t, http.StatusOK, w.Code)
		// user not found
		tok = thttp.UserToken(t, j, false, true)
		r = thttp.Request(t, method, route, tok, nil, &rest.Request{})
		r, err := rest.ParseClaims(r, j, tok)
		require.NoError(t, err)
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.NotEqual(t, http.StatusOK, w.Code)
		if checkReq {
			// invalid req
			tok = thttp.UserToken(t, j, false, true)
			r = thttp.Request(t, method, route, tok, nil, []byte("\n"))
			r, err = rest.ParseClaims(r, j, tok)
			require.NoError(t, err)
			w = httptest.NewRecorder()
			h.ServeHTTP(w, r)
			assert.NotEqual(t, http.StatusOK, w.Code)
		}
		// user banned
		u, tok := tcore.TestUser(t, srv.API, "", false)
		_, err = srv.BanUser(context.Background(), u.ID)
		require.NoError(t, err)
		r, err = rest.ParseClaims(r, j, tok)
		require.NoError(t, err)
		r = thttp.Request(t, method, route, tok, nil, nil)
		assert.NotEqual(t, http.StatusOK, w.Code)
	}
	const PassRoute = Endpoint + Password
	tests := []struct {
		m        string
		r        string
		checkReq bool
		h        http.HandlerFunc
	}{
		{http.MethodGet, Endpoint, false, srv.GetUser},
		{http.MethodPut, Endpoint, true, srv.UpdateUser},
		{http.MethodPut, PassRoute, true, srv.ChangePassword},
	}
	for _, test := range tests {
		testRoute(t, test.h, test.m, test.r, test.checkReq)
	}
}
