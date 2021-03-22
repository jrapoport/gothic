package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	host       = "example.com"
	restPort   = ":80"
	rpcPort    = ":90"
	webPort    = ":100"
	healthPort = ":110"
	requestID  = "foobar"
)

func TestNetwork(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		n := c.Network
		assert.Equal(t, host+test.mark, n.Host)
		assert.Equal(t, host+test.mark+restPort, n.REST)
		assert.Equal(t, host+test.mark+rpcPort, n.RPC)
		assert.Equal(t, host+test.mark+webPort, n.RPCWeb)
		assert.Equal(t, host+test.mark+healthPort, n.Health)
		assert.Equal(t, requestID+test.mark, n.RequestID)
	})
}

// tests the ENV vars are correctly taking precedence
func TestNetwork_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			n := c.Network
			assert.Equal(t, host, n.Host)
			assert.Equal(t, host+restPort, n.REST)
			assert.Equal(t, host+rpcPort, n.RPC)
			assert.Equal(t, host+webPort, n.RPCWeb)
			assert.Equal(t, host+healthPort, n.Health)
			assert.Equal(t, requestID, n.RequestID)
		})
	}
}

// test the *un-normalized* defaults with load
func TestNetwork_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := networkDefaults
	n := c.Network
	assert.Equal(t, def, n)
}

func TestNetwork_Normalization(t *testing.T) {
	const host2 = "peaches"
	def := networkDefaults
	t.Cleanup(func() {
		networkDefaults = def
	})
	netTests := []struct {
		uHost   string
		uRest   string
		uRPC    string
		uRPCWeb string
		nHost   string
		nRest   string
		nRPC    string
		nRPCWeb string
	}{
		{
			host, host + restPort, host + rpcPort, host + webPort,
			host, host + restPort, host + rpcPort, host + webPort,
		},
		{
			host2, "", "", "",
			host2, "", "", "",
		},
		{
			host2, host + restPort, host + rpcPort, host + webPort,
			host2, host + restPort, host + rpcPort, host + webPort,
		},
		{
			host2, def.REST, def.RPC, def.RPCWeb,
			host2, host2 + ":7727", host2 + ":7721", host2 + ":7729",
		},
		{
			host2, host + restPort, def.RPC, def.RPCWeb,
			host2, host + restPort, host2 + ":7721", host2 + ":7729",
		},
		{
			host2, host + restPort, host + restPort, host + restPort,
			host2, host + restPort, host + restPort, host + restPort,
		},
		{
			def.Host, def.REST, host + restPort, host + restPort,
			def.Host, def.REST, host + restPort, host + restPort,
		},
	}
	for _, test := range netTests {
		n := Network{}
		n.Host = test.uHost
		n.REST = test.uRest
		n.RPC = test.uRPC
		n.RPCWeb = test.uRPCWeb
		err := n.normalize(serviceDefaults)
		assert.NoError(t, err)
		assert.Equal(t, test.nHost, n.Host)
		assert.Equal(t, test.nRest, n.REST)
		assert.Equal(t, test.nRPC, n.RPC)
		assert.Equal(t, test.nRPCWeb, n.RPCWeb)
	}
	n := Network{}
	networkDefaults.Health = "::"
	networkDefaults.RPCWeb = "::"
	networkDefaults.RPC = "::"
	networkDefaults.REST = "::"
	n.Host = "::"
	n.Health = "::"
	err := n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.RPCWeb = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.RPC = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.REST = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
}
