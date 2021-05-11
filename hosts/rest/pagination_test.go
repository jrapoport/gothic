package rest

import (
	"crypto/tls"
	"github.com/jrapoport/gothic/store"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jrapoport/gothic/models/types/key"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginateRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)
	req.PostForm = url.Values{
		key.Page:     []string{"5"},
		key.PageSize: []string{"5"},
	}
	page := PaginateRequest(req)
	assert.Equal(t, 5, page.Index)
	assert.Equal(t, 5, page.Size)
}

func TestPaginateResponse(t *testing.T) {
	tests := []struct {
		scheme   string
		tls      *tls.ConnectionState
		forward  bool
		expected string
	}{
		{"http", nil, false, "http"},
		{"https", nil, false, "https"},
		{"http", &tls.ConnectionState{}, false, "https"},
		{"http", &tls.ConnectionState{}, true, "foo"},
	}
	for _, test := range tests {
		req, err := http.NewRequest(http.MethodGet, test.scheme+"://example.com", nil)
		require.NoError(t, err)
		req.TLS = test.tls
		if test.forward {
			req.Header.Set(ForwardedProto, test.expected)
		}
		rec := httptest.NewRecorder()
		PaginateResponse(rec, req, &store.Pagination{
			Index:  5,
			Size:   10,
			Prev:   5,
			Next:   6,
			Count:  10,
			Items:  [10]int{},
			Length: 10,
			Total:  100,
		})
		const link = "<SCHEME://example.com?page=1&per_page=10>; rel=\"first\"," +
			"<SCHEME://example.com?page=5&per_page=10>; rel=\"prev\"," +
			"<SCHEME://example.com?page=6&per_page=10>; rel=\"next\"," +
			"<SCHEME://example.com?page=10&per_page=10>; rel=\"last\""
		expected := strings.ReplaceAll(link, "SCHEME", test.expected)
		assert.Equal(t, expected, rec.Header().Get(Link))
	}
}
