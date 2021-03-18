package health

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	c := tconf.Config(t)
	h := Check(c)
	assert.Equal(t, h.Name, c.Name)
	assert.Equal(t, h.Version, c.Version())
	assert.NotEmpty(t, h.Status)
}
