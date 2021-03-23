package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	typ := None
	assert.False(t, typ.Has(Auth))
	typ = typ.Set(Auth)
	assert.True(t, typ.Has(Auth))
	assert.True(t, typ.Has(All))
	assert.False(t, typ.Has(Payment))
	typ = typ.Toggle(Payment)
	assert.True(t, typ.Has(Payment))
	typ = typ.Clear(Auth)
	assert.False(t, typ.Has(Auth))
	assert.True(t, typ.Has(Payment))
	typ = typ.Clear(All)
	assert.False(t, typ.Has(Payment))
	assert.False(t, typ.Has(All))
}

func TestType_String(t *testing.T) {
	typ := None
	assert.Equal(t, "", typ.String())
	typ = Auth | Payment | Wallet
	assert.Equal(t, "auth,payment,wallet", typ.String())
}
