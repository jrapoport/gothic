package thttp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/models/user"
	"github.com/stretchr/testify/require"
)

// Do does an http request
func Do(t *testing.T, s *httptest.Server,
	method, path string, v url.Values, body interface{}) (*http.Response, error) {
	return DoAuth(t, s, method, path, "", v, body)
}

// DoAuth does an authorized http request
func DoAuth(t *testing.T, s *httptest.Server,
	method, path, token string, v url.Values, body interface{}) (*http.Response, error) {
	req := Request(t, method, s.URL+path, token, v, body)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		var b string
		if res.Body != nil {
			buf, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = res.Body.Close()
				require.NoError(t, err)
			})
			b = string(buf)
		}
		if len(b) > 0 {
			b = strings.TrimRight(b, "\n")
			err = fmt.Errorf("%w: %s", err, b)
		}
		t.Log(err)
		return nil, err
	}
	return res, nil
}

// DoRequest makes an HTTP server request.
func DoRequest(t *testing.T, s *httptest.Server, method, path string, v url.Values, body interface{}) (string, error) {
	return DoAuthRequest(t, s, method, path, "", v, body)
}

// DoAuthRequest makes an authenticated HTTP server request.
func DoAuthRequest(t *testing.T, s *httptest.Server, method, path, token string, v url.Values, body interface{}) (string, error) {
	res, err := DoAuth(t, s, method, path, token, v, body)
	if err != nil {
		return "", err
	}
	var b string
	if res.Body != nil {
		buf, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		t.Cleanup(func() {
			err = res.Body.Close()
			require.NoError(t, err)
		})
		b = string(buf)
	}
	return b, err
}

// FmtError formats an error for tests.
func FmtError(code int) error {
	msg := http.StatusText(code)
	return fmt.Errorf("%d %s: %s", code, msg, msg)
}

// UserToken returns a dummy token for tests.
func UserToken(t *testing.T, c config.JWT, confirmed, admin bool) string {
	return token(t, c, uuid.New(), confirmed, admin)
}

// BadToken returns a bad token for tests.
func BadToken(t *testing.T, c config.JWT) string {
	return token(t, c, uuid.Nil, false, false)
}

func token(t *testing.T, c config.JWT, uid uuid.UUID, confirmed, admin bool) string {
	if admin {
		confirmed = true
	}
	claims := jwt.NewUserClaims(&user.User{})
	claims.StandardClaims = *jwt.NewStandardClaims(uid.String())
	claims.Admin = admin
	claims.Confirmed = confirmed
	tk, err := jwt.NewToken(c, claims).Bearer()
	require.NoError(t, err)
	return tk
}
