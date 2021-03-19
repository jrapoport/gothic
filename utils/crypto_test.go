package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecureToken(t *testing.T) {
	t.Parallel()
	tok1 := SecureToken()
	assert.NotEmpty(t, tok1)
	tok2 := SecureToken()
	assert.NotEmpty(t, tok2)
	assert.NotEqual(t, tok1, tok2)
}

func TestHashPassword(t *testing.T) {
	t.Parallel()
	const pw = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	hash, err := HashPassword(pw)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestMustHashPassword(t *testing.T) {
	t.Parallel()
	const pw = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	hash := MustHashPassword(pw)
	assert.NotEmpty(t, hash)
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()
	const pw = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"
	hash, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	err = CheckPassword(hash, pw)
	assert.NoError(t, err)
	err = CheckPassword(hash, "")
	assert.Error(t, err)
}
