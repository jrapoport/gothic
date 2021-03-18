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
