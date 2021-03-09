package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Config(t *testing.T) {
	a := apiWithTempDB(t)
	s := Server{API: a, FieldLogger: a.log}
	assert.Equal(t, a.config, s.Config())
}
