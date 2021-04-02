package utils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinLink(t *testing.T) {
	t.Parallel()
	const (
		example = "http://example.com"
		base    = "http://base.com"
		query1  = "?peaches=quack"
		query2  = "?quack=peaches"
	)
	tests := []struct {
		in   string
		frag string
		out  string
	}{
		{"", "/", "/"},
		{example, "/", example + "/"},
		{example + "/foo", "/foo", example + "/foo"},
		{example + "/foo", "/bar", example + "/bar"},
		{example + "/foo/bar", "/bar/foo", example + "/bar/foo"},
		{example + "/foo/bar" + query2, "/bar/foo", example + "/bar/foo"},
		{example + "/foo/bar" + query2, "/bar/foo" + query1, example + "/bar/foo" + query1},
		{example + "/foo", base + "/foo", example + "/foo"},
		{example + "/foo", base + "/bar", example + "/bar"},
		{example + "/foo/bar", base + "/bar/foo", example + "/bar/foo"},
		{example + "/foo/bar" + query2, base + "/bar/foo", example + "/bar/foo"},
		{example + "/foo/bar" + query2, base + "/bar/foo" + query1, example + "/bar/foo" + query1},
	}
	for _, test := range tests {
		out, err := JoinLink(test.in, test.frag)
		assert.NoError(t, err)
		assert.Equal(t, test.out, out)
	}
	_, err := JoinLink("\n", "/")
	assert.Error(t, err)
	_, err = JoinLink("/", "\n")
	assert.Error(t, err)
}

func TestAppendFragment(t *testing.T) {
	t.Parallel()
	const (
		example = "http://example.com"
		base    = "http://base.com"
		query1  = "?peaches=quack"
		query2  = "?quack=peaches"
	)
	p := func(s string) *url.URL { u, _ := url.Parse(s); return u }
	tests := []struct {
		in   *url.URL
		frag string
		out  *url.URL
	}{
		{p(""), "/", p("/")},
		{p(example), "/", p(example + "/")},
		{p(example + "/foo"), "/foo", p(example + "/foo")},
		{p(example + "/foo"), "/bar", p(example + "/bar")},
		{p(example + "/foo/bar"), "/bar/foo", p(example + "/bar/foo")},
		{p(example + "/foo/bar" + query2), "/bar/foo", p(example + "/bar/foo")},
		{p(example + "/foo/bar" + query2), "/bar/foo" + query1, p(example + "/bar/foo" + query1)},
		{p(example + "/foo"), base + "/foo", p(example + "/foo")},
		{p(example + "/foo"), base + "/bar", p(example + "/bar")},
		{p(example + "/foo/bar"), base + "/bar/foo", p(example + "/bar/foo")},
		{p(example + "/foo/bar" + query2), base + "/bar/foo", p(example + "/bar/foo")},
		{p(example + "/foo/bar" + query2), base + "/bar/foo" + query1, p(example + "/bar/foo" + query1)},
	}
	for _, test := range tests {
		out, err := appendFragment(test.in, test.frag)
		assert.NoError(t, err)
		assert.Equal(t, test.out, out)
	}
	_, err := appendFragment(p("/"), "\n")
	assert.Error(t, err)
}

func TestMakeRelative(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{"", "/"},
		{"http://example.com", "/"},
		{"http://example.com/foo", "/foo"},
		{"http://example.com/foo/bar", "/foo/bar"},
		{"http://example.com/foo/bar?quack=peaches", "/foo/bar?quack=peaches"},
		{"https://example.com/foo", "/foo"},
	}
	for _, test := range tests {
		out := makeRelative(test.in)
		assert.Equal(t, test.out, out)
	}
}

func TestURLValuesToMap(t *testing.T) {
	t.Parallel()
	const (
		testVal  = "test"
		testJSON = `{"test":"test"}`
	)
	v := url.Values{
		testVal: []string{testVal},
		"data":  []string{testJSON},
		"data2": []string{},
	}
	m := URLValuesToMap(v, false)
	val, ok := m[testVal]
	assert.True(t, ok)
	assert.Equal(t, testVal, val)
	val, ok = m["data"]
	assert.True(t, ok)
	assert.Equal(t, testJSON, val)
	// map data key
	m = URLValuesToMap(v, true)
	val, ok = m[testVal]
	assert.True(t, ok)
	assert.Equal(t, testVal, val)
	data, ok := m["data"].(map[string]interface{})
	assert.True(t, ok)
	val, ok = data[testVal]
	assert.True(t, ok)
	assert.Equal(t, testVal, val)
	_, ok = data["data2"]
	assert.False(t, ok)
}
