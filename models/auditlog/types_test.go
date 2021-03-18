package auditlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s string
		t Type
	}{
		{"system", System},
		{"account", Account},
		{"token", Token},
		{"user", User},
		{"", Unknown},
	}
	for _, test := range tests {
		typ := TypeFromString(test.s)
		assert.Equal(t, test.t, typ)
		assert.Equal(t, test.s, typ.String())
	}
}
