package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPINCode(t *testing.T) {
	t.Parallel()
	pin := PINCode()
	assert.NotEmpty(t, pin)
}

func TestIsValidPIN(t *testing.T) {
	t.Parallel()
	pin := PINCode()
	require.NotEmpty(t, pin)
	is := IsValidCode(pin)
	assert.True(t, is)
	is = IsValidCode("")
	assert.False(t, is)
}
