package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	x := Max(1, 10)
	assert.Equal(t, 10, x)
	x = Max(20, 10)
	assert.Equal(t, 20, x)
}

func TestClamp(t *testing.T) {
	tests := []struct {
		x        int
		min      int
		max      int
		expected int
	}{
		{0, 0, 0, 0},
		{0, 1, 0, 1},
		{0, 0, 1, 0},
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{1, 0, 1, 1},
		{3, 2, 5, 3},
		{10, 1, 0, 0},
		{10, 20, 0, 20},
		{10, 0, 5, 5},
		{10, 20, 5, 20},
	}
	for _, test := range tests {
		n := Clamp(test.x, test.min, test.max)
		assert.Equal(t, test.expected, n)
	}
}
