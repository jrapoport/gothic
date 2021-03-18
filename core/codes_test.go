package core

import (
	"testing"

	"github.com/jrapoport/gothic/models/code"
	"github.com/stretchr/testify/assert"
)

func TestAPI_CreateCode(t *testing.T) {
	t.Parallel()
	a := createAPI(t)
	ctx := testContext(a)
	c, err := a.CreateSignupCode(ctx, code.SingleUse)
	assert.NoError(t, err)
	assert.NotEmpty(t, c)
}

func TestAPI_CreateCodes(t *testing.T) {
	t.Parallel()
	const count = 100
	a := createAPI(t)
	ctx := testContext(a)
	list, err := a.CreateSignupCodes(ctx, code.SingleUse, count)
	assert.NoError(t, err)
	assert.Len(t, list, count)
	assertUnique(t, list)
}

func TestAPI_CheckSignupCode(t *testing.T) {
	t.Parallel()
	a := createAPI(t)
	ctx := testContext(a)
	test, err := a.CreateSignupCode(ctx, code.SingleUse)
	assert.NoError(t, err)
	assert.NotEmpty(t, test)
	sc, err := a.CheckSignupCode(test)
	assert.NoError(t, err)
	assert.Equal(t, test, sc.Token)
	_, err = a.CheckSignupCode("")
	assert.Error(t, err)
}

func assertUnique(t *testing.T, s []string) {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	assert.Exactly(t, s, us)
}
