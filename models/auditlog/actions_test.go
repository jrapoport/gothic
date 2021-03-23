package auditlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAction_Type(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a Action
		t Type
	}{
		{Banned, Account},
		{CodeSent, Account},
		{ConfirmSent, Account},
		{Confirmed, Account},
		{Deleted, Account},
		{Signup, Account},
		{Startup, System},
		{Shutdown, System},
		{Granted, Token},
		{Refreshed, Token},
		{Revoked, Token},
		{RevokedAll, Token},
		{ChangeRole, User},
		{Email, User},
		{Linked, User},
		{Login, User},
		{Logout, User},
		{Password, User},
		{Updated, User},
		{"", Unknown},
	}
	for _, test := range tests {
		typ := test.a.Type()
		assert.Equal(t, test.t, typ)
	}
	assert.Equal(t, string(Banned), Banned.String())
}
