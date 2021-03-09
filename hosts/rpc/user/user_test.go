package user

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/test/tcore"
	"testing"

	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"

func testServer(t *testing.T) *userServer {
	srv, _ := tsrv.RPCServer(t, false)
	srv.Config().Signup.AutoConfirm = true
	return newUserServer(srv)
}

func TestUserServer_GetUser(t *testing.T) {
	srv := testServer(t)
	srv.Config().MaskEmails = false
	req := &GetUserRequest{}
	ctx := context.Background()
	// no id
	_, err := srv.GetUser(ctx, req)
	assert.Error(t, err)
	// unmasked
	u, tok := tcore.TestUser(t, srv.API, "", false)
	ctx = tsrv.RPCAuthContext(t, srv.Config(), tok)
	res, err := srv.GetUser(ctx, req)
	assert.NoError(t, err)
	test, err := rpc.NewUserResponse(u)
	assert.NoError(t, err)
	assert.Equal(t, test.Email, res.Email)
	assert.Equal(t, test.Username, res.Username)
	assert.Equal(t, test.Data.AsMap(), res.Data.AsMap())
	// masked
	srv.Config().MaskEmails = true
	test.Email = utils.MaskEmail(u.Email)
	res, err = srv.GetUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, test.Email, res.Email)
	assert.Equal(t, test.Username, res.Username)
	assert.Equal(t, test.Data.AsMap(), res.Data.AsMap())
}

func TestUserServer_UpdateUser(t *testing.T) {
	srv := testServer(t)
	u, tok := tcore.TestUser(t, srv.API, "", false)
	ctx := tsrv.RPCAuthContext(t, srv.Config(), tok)
	data, err := structpb.NewStruct(types.Map{
		"foo":   "bar",
		"tasty": "salad",
	})
	require.NoError(t, err)
	req := &UpdateUserRequest{
		Username: "peaches",
		Data:     data,
	}
	// invalid username
	srv.Config().Validation.UsernameRegex = "0"
	_, err = srv.UpdateUser(ctx, req)
	assert.Error(t, err)
	// unmasked
	srv.Config().MaskEmails = false
	srv.Config().Validation.UsernameRegex = ""
	res, err := srv.UpdateUser(ctx, req)
	assert.NoError(t, err)
	err = mergo.Map(&u.Data, req.Data.AsMap())
	assert.NoError(t, err)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, req.Username, res.Username)
	assert.EqualValues(t, u.Data, res.Data.AsMap())
	// masked
	srv.Config().MaskEmails = true
	data, err = structpb.NewStruct(types.Map{
		"quack":   99.0,
		"peaches": "happy",
	})
	require.NoError(t, err)
	req = &UpdateUserRequest{
		Username: "mario",
		Data:     data,
	}
	res, err = srv.UpdateUser(ctx, req)
	assert.NoError(t, err)
	err = mergo.Map(&u.Data, req.Data.AsMap())
	assert.NoError(t, err)
	assert.Equal(t, utils.MaskEmail(u.Email), res.Email)
	assert.Equal(t, req.Username, res.Username)
	assert.EqualValues(t, u.Data, res.Data.AsMap())
}

func TestUserServer_ChangePassword(t *testing.T) {
	const newPassword = "gj8#xtg#yrabxpnno!p5f3t8na!hd3?4jq7majxs"
	srv := testServer(t)
	u, tok := tcore.TestUser(t, srv.API, "", false)
	ctx := tsrv.RPCAuthContext(t, srv.Config(), tok)
	req := &ChangePasswordRequest{
		Password: newPassword,
	}
	// invalid password
	srv.Config().Validation.PasswordRegex = "0"
	_, err := srv.ChangePassword(ctx, req)
	assert.Error(t, err)
	// wrong old password
	srv.Config().Validation.PasswordRegex = ""
	_, err = srv.ChangePassword(ctx, req)
	assert.Error(t, err)
	// success
	req.OldPassword = testPass
	res, err := srv.ChangePassword(ctx, req)
	assert.NoError(t, err)
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Access)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, u.ID.String(), claims.Subject)
	u, err = srv.API.GetUser(u.ID)
	assert.NoError(t, err)
	err = u.Authenticate(newPassword)
	assert.NoError(t, err)
}

func TestRequestErrors(t *testing.T) {
	srv := testServer(t)
	ctx := context.Background()
	_, err := srv.GetUser(ctx, &GetUserRequest{})
	assert.Error(t, err)
	_, err = srv.UpdateUser(ctx, &UpdateUserRequest{})
	assert.Error(t, err)
	_, err = srv.ChangePassword(ctx, &ChangePasswordRequest{})
	assert.Error(t, err)
	claims := jwt.UserClaims{}
	claims.Subject = uuid.New().String()
	ctx = context.WithContext(rpc.WithClaims(ctx, claims))
	_, err = srv.GetUser(ctx, &GetUserRequest{})
	assert.Error(t, err)
	_, err = srv.UpdateUser(ctx, &UpdateUserRequest{})
	assert.Error(t, err)
	_, err = srv.ChangePassword(ctx, &ChangePasswordRequest{})
	assert.Error(t, err)
}

/*
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
		tok := thttp.BadToken(t, j, false, true)
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
		u, tok := tsrv.TestUser(t, srv.Server.Server, "", false)
		_, err = srv.BanUser(context.Background(), u.ID)
		require.NoError(t, err)
		r, err = rest.ParseClaims(r, j, tok)
		require.NoError(t, err)
		r = thttp.Request(t, method, route, tok, nil, nil)
		assert.NotEqual(t, http.StatusOK, w.Code)
	}
	const PassRoute = Endpoint + Password
	tests := []struct {
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

*/
