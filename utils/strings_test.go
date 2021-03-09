package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespaced(t *testing.T) {
	const name = "bar"
	const namespace = "foo"
	test := Namespaced("", name)
	assert.Equal(t, name, test)
	test = Namespaced(namespace, "")
	assert.Equal(t, namespace+"_", test)
	test = Namespaced(namespace, name)
	assert.Equal(t, namespace+"_"+name, test)
}

func TestMaskString(t *testing.T) {
	tests := []struct {
		in  string
		n   int
		out string
	}{
		{"", 0, ""},
		{"", 1, ""},
		{"a", 0, "*"},
		{"a", 1, "*"},
		{"abcd", 2, "ab**"},
	}
	for _, test := range tests {
		out := MaskString(test.in, test.n)
		assert.Equal(t, test.out, out)
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", ""},
		{"abcd", ""},
		{"foo@bar", ""},
		{"foo@bar.co", "f*o@b**.co"},
		{"foo@bar.com", "f*o@b**.com"},
		{"quackfoo@example.com", "qu******@e******.com"},
	}
	for _, test := range tests {
		out := MaskEmail(test.in)
		assert.Equal(t, test.out, out)
	}
}
