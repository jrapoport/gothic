package users

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createUserTest(t *testing.T) (url.Values, *rest.UserResponse) {
	const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	em := tutils.RandomEmail()
	un := utils.RandomUsername()
	data := types.Map{
		"foo":   "bar",
		"tasty": "salad",
	}
	d, err := data.JSON()
	require.NoError(t, err)
	v := url.Values{}
	v.Set(key.Email, em)
	v.Set(key.Username, un)
	v.Set(key.Password, testPass)
	v.Set(key.Data, string(d))
	ur := &rest.UserResponse{
		Role:     user.RoleUser.String(),
		Email:    em,
		Username: un,
		Data:     data,
	}
	return v, ur
}

func TestUsersServer_AdminCreateUser(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newUserServer(s)
	v, test := createUserTest(t)
	createUser := func(tok string) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodPost, Endpoint, tok, v, nil)
		if tok != "" {
			var err error
			r, err = rest.ParseClaims(r, srv.Config().JWT, tok)
			require.NoError(t, err)
		}
		w := httptest.NewRecorder()
		srv.AdminCreateUser(w, r)
		return w
	}
	res := createUser("")
	assert.NotEqual(t, http.StatusOK, res.Code)
	_, tok := testUser(t, srv, false)
	res = createUser(tok)
	assert.NotEqual(t, http.StatusOK, res.Code)
	_, tok = testUser(t, srv, true)
	res = createUser(tok)
	assert.Equal(t, http.StatusOK, res.Code)
	ur := userResponse(t, res.Body.String())
	assert.Equal(t, test, ur)
	r := thttp.Request(t, http.MethodPost, Endpoint, tok, nil, nil)
	r, err := rest.ParseClaims(r, srv.Config().JWT, tok)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	srv.AdminCreateUser(w, r)
	assert.NotEqual(t, http.StatusOK, w.Code)
	// invalid req
	r = thttp.Request(t, http.MethodPost, Endpoint, tok, nil, []byte("\n"))
	r, err = rest.ParseClaims(r, srv.Config().JWT, tok)
	require.NoError(t, err)
	w = httptest.NewRecorder()
	srv.AdminCreateUser(w, r)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
