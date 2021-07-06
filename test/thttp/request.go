package thttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/jrapoport/gothic/utils"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/require"
)

// Request makes an HTTP server request w/ authentication.
// IF values != nil && v == nil => url query string request
// IF values != nil && v != nil => x-www-form-urlencoded request
// IF values == nil && v != nil => json request
func Request(t *testing.T, method, path, token string, v url.Values, body interface{}) *http.Request {
	const (
		Authorization = "Authorization"
		Bearer        = "Bearer "
		ContentType   = "Content-Type"
		FORM          = "application/x-www-form-urlencoded; param=value"
		JSON          = "application/json"
	)
	var reqBody io.Reader
	if v != nil && body != nil {
		reqBody = bytes.NewBufferString(v.Encode())
	} else if body != nil {
		switch typ := body.(type) {
		case url.Values:
			body = utils.URLValuesToMap(typ, true)
		default:
			break
		}
		buffer, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(buffer)
	}
	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(t, err)
	if token != "" {
		req.Header.Set(Authorization, Bearer+token)
	}
	if body != nil {
		if v != nil {
			// This makes it work
			req.Header.Set(ContentType, FORM)
		} else {
			req.Header.Set(ContentType, JSON)
		}
	} else if v != nil {
		req.URL.RawQuery = v.Encode()
	}
	return req
}
