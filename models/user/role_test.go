package user

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	tests := []struct {
		role  Role
		name  string
		Valid assert.BoolAssertionFunc
	}{
		{-2, "", assert.False},
		{RoleSystem, "system", assert.True},
		{RoleUser, "user", assert.True},
		{RoleAdmin, "admin", assert.True},
		{RoleSuper, "super", assert.True},
		{math.MaxInt8, "", assert.False},
	}
	for _, test := range tests {
		r := test.role
		assert.Equal(t, test.name, r.String())
		test.Valid(t, r.Valid())
	}
}
