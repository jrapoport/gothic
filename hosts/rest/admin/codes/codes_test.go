package codes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignupServer_CreateSignupCodes(t *testing.T) {
	t.Parallel()
	const testLen = 10
	s, _ := tsrv.RESTServer(t, false)
	srv := newSignupServer(s)
	j := srv.Config().JWT
	createCodes := func(tok string, body interface{}) *httptest.ResponseRecorder {
		r := thttp.Request(t, http.MethodPost, Codes, tok, nil, body)
		if tok != "" {
			var err error
			r, err = rest.ParseClaims(r, j, tok)
			require.NoError(t, err)
		}
		w := httptest.NewRecorder()
		srv.CreateSignupCodes(w, r)
		return w
	}
	// no user id slug
	tok := thttp.UserToken(t, j, false, false)
	res := createCodes(tok, nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// no admin id
	res = createCodes("", nil)
	assert.NotEqual(t, http.StatusOK, res.Code)
	// bad request
	res = createCodes("", []byte("\n"))
	assert.NotEqual(t, http.StatusOK, res.Code)
	// admin not found
	res = createCodes(tok, types.Map{})
	assert.NotEqual(t, http.StatusOK, res.Code)
	_, tok = tcore.TestUser(t, srv.API, "", true)
	res = createCodes(tok, types.Map{
		key.Uses:  code.InfiniteUse,
		key.Count: testLen,
	})
	assert.Equal(t, http.StatusOK, res.Code)
	var list []string
	err := json.Unmarshal(res.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, testLen)
}

func TestSignupServer_CheckSignupCode(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newSignupServer(s)
	j := srv.Config().JWT
	tok := thttp.UserToken(t, j, true, true)
	scode, err := srv.API.CreateSignupCode(context.Background(), 0)
	require.NoError(t, err)
	checkCode := func(cd string) *httptest.ResponseRecorder {
		uri := Codes + rest.Root + cd
		r := thttp.Request(t, http.MethodGet, uri, tok, nil, nil)
		if tok != "" {
			r, err = rest.ParseClaims(r, j, tok)
			require.NoError(t, err)
		}
		ctx := chi.NewRouteContext()
		ctx.URLParams = chi.RouteParams{
			Keys:   []string{key.Code},
			Values: []string{cd},
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		srv.CheckSignupCode(w, r)
		return w
	}
	// no code
	res := checkCode("")
	assert.NotEqual(t, http.StatusOK, res.Code)
	// bad code
	res = checkCode("bad")
	assert.NotEqual(t, http.StatusOK, res.Code)
	// good code
	res = checkCode(scode)
	assert.Equal(t, http.StatusOK, res.Code)
	var cr map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&cr)
	require.NoError(t, err)
	assert.Equal(t, scode, cr[key.Code])
	assert.True(t, cr[key.Valid].(bool))
	// burn the code
	conn := tconn.Conn(t, srv.Config())
	sc, err := codes.GetSignupCode(conn, scode)
	require.NoError(t, err)
	require.NotNil(t, sc)
	sc.Used = 1
	err = conn.Save(sc).Error
	require.NoError(t, err)
	// used code
	res = checkCode(scode)
	assert.Equal(t, http.StatusOK, res.Code)
	err = json.NewDecoder(res.Body).Decode(&cr)
	require.NoError(t, err)
	assert.Equal(t, scode, cr[key.Code])
	assert.False(t, cr[key.Valid].(bool))
}

func TestSignupServer_DeleteSignupCode(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RESTServer(t, false)
	srv := newSignupServer(s)
	j := srv.Config().JWT
	tok := thttp.UserToken(t, j, true, true)
	scode, err := srv.API.CreateSignupCode(context.Background(), 0)
	require.NoError(t, err)
	deleteCode := func(cd string) *httptest.ResponseRecorder {
		uri := Codes + rest.Root + cd
		r := thttp.Request(t, http.MethodDelete, uri, tok, nil, nil)
		if tok != "" {
			r, err = rest.ParseClaims(r, j, tok)
			require.NoError(t, err)
		}
		ctx := chi.NewRouteContext()
		ctx.URLParams = chi.RouteParams{
			Keys:   []string{key.Code},
			Values: []string{cd},
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		srv.DeleteSignupCode(w, r)
		return w
	}
	// bad code
	res := deleteCode("bad")
	assert.NotEqual(t, http.StatusOK, res.Code)
	// delete code
	res = deleteCode(scode)
	assert.Equal(t, http.StatusOK, res.Code)
	// code is gone
	conn := tconn.Conn(t, srv.Config())
	sc, err := codes.GetSignupCode(conn, scode)
	assert.Error(t, err)
	assert.Nil(t, sc)
}
