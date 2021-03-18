package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	h := a.HealthCheck()
	assert.Equal(t, h.Name, a.config.Name)
	assert.Equal(t, h.Version, a.config.Version())
	assert.NotEmpty(t, h.Status)
}
