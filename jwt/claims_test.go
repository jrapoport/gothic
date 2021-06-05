package jwt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStandardClaims(t *testing.T) {
	sub1 := uuid.New().String()
	claims := NewStandardClaims(sub1)
	assert.Equal(t, sub1, claims.Subject())
	assert.Equal(t, "", claims.Issuer())
	sub2 := uuid.New().String()
	claims.SetSubject(sub2)
	assert.Equal(t, sub2, claims.Subject())
	assert.Equal(t, []string{}, claims.Scope())
	err := claims.Set(ScopeKey, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{}, claims.Scope())
}
