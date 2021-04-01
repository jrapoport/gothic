package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutboundIP(t *testing.T) {
	t.Parallel()
	ip, err := OutboundIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
	dns = "255.255.255.255"
	_, err = OutboundIP()
	assert.Error(t, err)
}

func TestMakeAddress(t *testing.T) {
	addr := MakeAddress("127.0.0.1", 999)
	assert.Equal(t, "127.0.0.1:999", addr)
}
