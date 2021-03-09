package jwt

import (
	"testing"

	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseValue(t *testing.T) {
	c := tconf.Config(t)
	in := types.Map{
		key.Token: utils.SecureToken(),
		key.Email: tutils.RandomEmail(),
	}
	tok, err := NewSignedData(c.JWT, in)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)
	out, err := ParseData(c.JWT, tok)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out[key.Token], in[key.Token])
	assert.Equal(t, out[key.Email], in[key.Email])
}
