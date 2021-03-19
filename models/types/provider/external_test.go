package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddExternal(t *testing.T) {
	var n Name = "new_provider"
	assert.False(t, IsExternal(n))
	AddExternal(n)
	assert.True(t, IsExternal(n))
}

func TestIsExternal(t *testing.T) {
	var n Name = "test"
	assert.False(t, IsExternal(n))
	n = Unknown
	assert.False(t, IsExternal(n))
	n = Google
	assert.True(t, IsExternal(n))
}
