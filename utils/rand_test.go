package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomUsername(t *testing.T) {
	un := RandomUsername()
	assert.NotEmpty(t, un)
}

func TestRandomUsernameN(t *testing.T) {
	tests := []int{10, 20, 100}
	for _, test := range tests {
		un := RandomUsernameN(test)
		assert.LessOrEqual(t, len(un), test)
	}
	un := RandomUsernameN(0)
	assert.NotEmpty(t, un)
}

func TestRandomColor(t *testing.T) {
	clr := RandomColor()
	assert.Len(t, clr, 7)
	assert.True(t, strings.HasPrefix(clr, "#"))
}

func TestRandomPIN(t *testing.T) {
	tests := []int{10, 20, 100}
	for _, test := range tests {
		pin := RandomPIN(test)
		assert.LessOrEqual(t, len(pin), test)
	}
	pin := RandomPIN(0)
	assert.Empty(t, pin)
}
