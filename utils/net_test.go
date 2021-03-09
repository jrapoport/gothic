package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutboundIP(t *testing.T) {
	ip, err := OutboundIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
}
