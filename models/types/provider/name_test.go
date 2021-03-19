package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	var n Name = "test"
	assert.False(t, n.IsExternal())
	n = Unknown
	assert.False(t, n.IsExternal())
	n = Google
	assert.True(t, n.IsExternal())
	assert.Equal(t, "google", n.String())
	assert.Equal(t, "0990a1ac-3962-31df-bc95-31ccf169044c", n.ID().String())
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		name     string
		provider Name
	}{
		{"", Unknown},
		{"abcdABCD1234", Name("abcdabcd1234")},
		{"#abcd!ABCD 1234", Name("abcdabcd1234")},
	}
	for _, test := range tests {
		p := NormalizeName(test.name)
		assert.Equal(t, test.provider, p)
	}
}
