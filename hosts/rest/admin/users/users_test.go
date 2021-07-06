package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconn"
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

func testUser(t *testing.T, srv *usersServer, admin bool) (*user.User, string) {
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
	if admin {
		conn := tconn.Conn(t, srv.Config())
		u.Role = user.RoleSuper
		err = conn.Save(u).Error
		require.NoError(t, err)
	}
	bt, err := srv.GrantBearerToken(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, bt)
	return u, bt.String()
}

func TestUserServer_GetUser(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	srv.Config().Signup.Default.Color = false
	j := srv.Config().JWT
	u, _ := testUser(t, srv, false)
	test := &rest.UserResponse{
		UserID:   u.ID.String(),
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: u.Username,
		Data:     u.Data,
	}
	uri := Users + rest.Root + u.ID.String()
	getUser := func(tok string, useCtx bool, testID uuid.UUID) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodGet, uri, tok, nil, nil)
		if tok != "" {
			var err error
			r, err = rest.ParseClaims(r, srv.Config().JWT, tok)
			require.NoError(t, err)
		}
		uid := u.ID
		if testID != uuid.Nil {
			uid = testID
		}
		if useCtx {
			ctx := chi.NewRouteContext()
			ctx.URLParams = chi.RouteParams{
				Keys:   []string{key.UserID},
				Values: []string{uid.String()},
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		}
		w := httptest.NewRecorder()
		srv.GetUser(w, r)
		return w
	}
	// no user id slug
	tok := thttp.UserToken(t, j, false, false)
	res := getUser(tok, false, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// no admin id
	res = getUser("", true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin not found
	tok = thttp.UserToken(t, j, false, false)
	res = getUser(tok, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// not admin
	_, tok = testUser(t, srv, false)
	res = getUser(tok, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin
	_, tok = testUser(t, srv, true)
	res = getUser(tok, true, uuid.Nil)
	assert.Equal(t, http.StatusOK, res.Code)
	ur := userResponse(t, res.Body.String())
	assert.Equal(t, test, ur)
	// user not found
	res = getUser(tok, true, uuid.New())
	assert.NotEqual(t, http.StatusOK, res.Code)
}

func TestUserServer_UpdateUser(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	srv.Config().Signup.Default.Color = false
	j := srv.Config().JWT
	u, _ := testUser(t, srv, false)
	uri := Users + rest.Root + u.ID.String()
	updateUser := func(tok string, body interface{}, useCtx bool, testID uuid.UUID) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodPut, uri, tok, nil, body)
		if tok != "" {
			var err error
			r, err = rest.ParseClaims(r, srv.Config().JWT, tok)
			require.NoError(t, err)
		}
		uid := u.ID
		if testID != uuid.Nil {
			uid = testID
		}
		if useCtx {
			ctx := chi.NewRouteContext()
			ctx.URLParams = chi.RouteParams{
				Keys:   []string{key.UserID},
				Values: []string{uid.String()},
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		}
		w := httptest.NewRecorder()
		srv.UpdateUser(w, r)
		return w
	}
	// no user id slug
	tok := thttp.UserToken(t, j, false, false)
	res := updateUser(tok, nil, false, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// no admin id
	res = updateUser("", nil, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin not found
	tok = thttp.UserToken(t, j, false, false)
	res = updateUser(tok, nil, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// not admin
	_, tok = testUser(t, srv, false)
	res = updateUser(tok, nil, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	req := &Request{
		Username: "peaches",
		Data: types.Map{
			"foo":   "bar",
			"tasty": "salad",
		},
	}
	// admin
	_, tok = testUser(t, srv, true)
	// invalid username
	srv.Config().Validation.UsernameRegex = "0"
	res = updateUser(tok, req, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// bad req
	res = updateUser(tok, []byte("\n"), true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	srv.Config().Validation.UsernameRegex = ""
	res = updateUser(tok, req, true, uuid.Nil)
	assert.Equal(t, http.StatusOK, res.Code)
	ur := userResponse(t, res.Body.String())
	err := mergo.Map(&u.Data, req.Data)
	assert.NoError(t, err)
	test := &rest.UserResponse{
		UserID:   u.ID.String(),
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: req.Username,
		Data:     u.Data,
	}
	assert.Equal(t, test, ur)
}

func TestUserServer_DeleteUser(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	srv.Config().Signup.Default.Color = false
	j := srv.Config().JWT
	u, _ := testUser(t, srv, false)
	uri := Users + rest.Root + u.ID.String()
	deleteUser := func(tok string, useCtx bool, testID uuid.UUID) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodGet, uri, tok, nil, nil)
		if tok != "" {
			var err error
			r, err = rest.ParseClaims(r, srv.Config().JWT, tok)
			require.NoError(t, err)
		}
		uid := u.ID
		if testID != uuid.Nil {
			uid = testID
		}
		if useCtx {
			ctx := chi.NewRouteContext()
			ctx.URLParams = chi.RouteParams{
				Keys:   []string{key.UserID},
				Values: []string{uid.String()},
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		}
		w := httptest.NewRecorder()
		srv.AdminDeleteUser(w, r)
		return w
	}
	// no user id slug
	tok := thttp.UserToken(t, j, false, false)
	res := deleteUser(tok, false, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// no admin id
	res = deleteUser("", true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin not found
	tok = thttp.UserToken(t, j, false, false)
	res = deleteUser(tok, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// not admin
	_, tok = testUser(t, srv, false)
	res = deleteUser(tok, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin
	_, tok = testUser(t, srv, true)
	res = deleteUser(tok, true, uuid.Nil)
	assert.Equal(t, http.StatusOK, res.Code)
	// user should come back as not found
	_, err := srv.API.GetUser(u.ID)
	assert.Error(t, err)
	// user not found
	res = deleteUser(tok, true, uuid.New())
	assert.NotEqual(t, http.StatusOK, res.Code)
}

func TestUserServer_PromoteUser(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	srv.Config().Signup.AutoConfirm = true
	srv.Config().MaskEmails = false
	srv.Config().Signup.Default.Color = false
	j := srv.Config().JWT
	u, _ := testUser(t, srv, false)
	test := &rest.UserResponse{
		UserID:   u.ID.String(),
		Role:     user.RoleAdmin.String(),
		Email:    u.Email,
		Username: u.Username,
		Data:     u.Data,
	}
	uri := Users + rest.Root + u.ID.String() + Promote
	promoteUser := func(tok string, useCtx bool, testID uuid.UUID) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodGet, uri, tok, nil, nil)
		if tok != "" {
			var err error
			r, err = rest.ParseClaims(r, srv.Config().JWT, tok)
			require.NoError(t, err)
		}
		uid := u.ID
		if testID != uuid.Nil {
			uid = testID
		}
		if useCtx {
			ctx := chi.NewRouteContext()
			ctx.URLParams = chi.RouteParams{
				Keys:   []string{key.UserID},
				Values: []string{uid.String()},
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		}
		w := httptest.NewRecorder()
		srv.AdminPromoteUser(w, r)
		return w
	}
	// no user id slug
	tok := thttp.UserToken(t, j, false, false)
	res := promoteUser(tok, false, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// no admin id
	res = promoteUser("", true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin not found
	tok = thttp.UserToken(t, j, false, false)
	res = promoteUser(tok, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// not admin
	_, tok = testUser(t, srv, false)
	res = promoteUser(tok, true, uuid.Nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin
	_, tok = testUser(t, srv, true)
	res = promoteUser(tok, true, uuid.Nil)
	assert.Equal(t, http.StatusOK, res.Code)
	ur := userResponse(t, res.Body.String())
	assert.Equal(t, test, ur)
	// user not found
	res = promoteUser(tok, true, uuid.New())
	assert.NotEqual(t, http.StatusOK, res.Code)
}
