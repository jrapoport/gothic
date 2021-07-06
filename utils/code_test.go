package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPINCode(t *testing.T) {
	pin := PINCode()
	assert.NotEmpty(t, pin)
}

func TestIsValidPIN(t *testing.T) {
	pin := PINCode()
	require.NotEmpty(t, pin)
	is := IsValidCode(pin)
	assert.True(t, is)
	is = IsValidCode("")
	assert.False(t, is)
}
