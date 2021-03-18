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
