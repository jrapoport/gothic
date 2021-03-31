package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Config(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	s := Server{API: a, Logger: a.log}
	assert.Equal(t, a.config, s.Config())
}
