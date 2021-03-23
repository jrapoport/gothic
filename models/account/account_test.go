package account

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_BeforeSave(t *testing.T) {
	var (
		name = provider.Google
		aid  = uuid.New().String()
	)
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	la := &Account{}
	err := conn.Create(la).Error
	assert.Error(t, err)
	la.Type = Auth
	err = conn.Create(la).Error
	assert.Error(t, err)
	la.Provider = name
	err = conn.Create(la).Error
	assert.Error(t, err)
	la.AccountID = aid
	err = conn.Create(la).Error
	assert.NoError(t, err)
	// Provider + AccountID must be unique
	la = &Account{}
	la.Type = Auth
	la.Provider = name
	la.AccountID = aid
	err = conn.Create(la).Error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestAccount_Valid(t *testing.T) {
	t.Parallel()
	la := NewAccount(provider.Unknown, "", "", nil)
	err := la.Valid()
	assert.Error(t, err)
	la.Type = None
	err = la.Valid()
	assert.Error(t, err)
	la.Type = Auth
	err = la.Valid()
	assert.Error(t, err)
	la.Provider = "test provider"
	err = la.Valid()
	assert.Error(t, err)
	const accountID = "test-id"
	la.AccountID = accountID
	err = la.Valid()
	assert.NoError(t, err)
	la.CreatedAt = time.Now()
	err = la.Valid()
	assert.Error(t, err)
	la.UserID = uuid.New()
	err = la.Valid()
	assert.NoError(t, err)
}

func TestAccount_Type(t *testing.T) {
	tests := []struct {
		name      provider.Name
		isAuth    assert.BoolAssertionFunc
		isPayment assert.BoolAssertionFunc
		isWallet  assert.BoolAssertionFunc
	}{
		{provider.Unknown, assert.False, assert.False, assert.False},
		{provider.Google, assert.True, assert.False, assert.False},
		{provider.Stripe, assert.True, assert.True, assert.False},
		{provider.PayPal, assert.True, assert.True, assert.True},
	}
	for _, test := range tests {
		la := NewAccount(test.name, "", "", nil)
		test.isAuth(t, la.HasType(Auth))
		test.isPayment(t, la.HasType(Payment))
		test.isWallet(t, la.HasType(Wallet))
	}
}

func TestProviderType(t *testing.T) {
	tests := []struct {
		name      provider.Name
		isAuth    assert.BoolAssertionFunc
		isPayment assert.BoolAssertionFunc
		isWallet  assert.BoolAssertionFunc
	}{
		{provider.PayPal, assert.True, assert.True, assert.True},
		{provider.Stripe, assert.True, assert.True, assert.False},
		{provider.Google, assert.True, assert.False, assert.False},
	}
	for _, test := range tests {
		typ := providerType(test.name)
		test.isAuth(t, typ.Has(Auth))
		test.isPayment(t, typ.Has(Payment))
		test.isWallet(t, typ.Has(Wallet))
	}
}
