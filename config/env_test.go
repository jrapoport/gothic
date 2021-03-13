package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Env(t *testing.T) {
	c := Config{}
	dbg := c.IsDebug()
	assert.True(t, dbg)
	e := c.Env()
	assert.Equal(t, "debug", e)
}
