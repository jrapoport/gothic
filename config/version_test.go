package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVersion(t *testing.T) {
	ver = version()
	assert.Equal(t, "debug", BuildVersion())
	Version = "1.0"
	ver = version()
	assert.Equal(t, "1.0", BuildVersion())
	Build = "x100"
	ver = version()
	assert.Equal(t, "1.0 (x100)", BuildVersion())
	Version = ""
	ver = version()
	assert.Equal(t, "x100", BuildVersion())
}
