package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testHost   = "example.com"
	adminPort  = ":91"
	restPort   = ":80"
	rpcPort    = ":90"
	webPort    = ":100"
	healthPort = ":110"
	requestID  = "foobar"
)

func TestNetwork(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		n := c.Network
		assert.Equal(t, testHost+test.mark, n.Host)
		assert.Equal(t, testHost+test.mark+adminPort, n.AdminAddress)
		assert.Equal(t, testHost+test.mark+restPort, n.RESTAddress)
		assert.Equal(t, testHost+test.mark+rpcPort, n.RPCAddress)
		assert.Equal(t, testHost+test.mark+webPort, n.RPCWebAddress)
		assert.Equal(t, testHost+test.mark+healthPort, n.HealthAddress)
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
			assert.Equal(t, testHost, n.Host)
			assert.Equal(t, testHost+adminPort, n.AdminAddress)
			assert.Equal(t, testHost+restPort, n.RESTAddress)
			assert.Equal(t, testHost+rpcPort, n.RPCAddress)
			assert.Equal(t, testHost+webPort, n.RPCWebAddress)
			assert.Equal(t, testHost+healthPort, n.HealthAddress)
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
	networkDefaults.HealthAddress = def.RPCAddress
	networkDefaults.AdminAddress = def.RPCAddress
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
			testHost, testHost + restPort, testHost + rpcPort, testHost + webPort,
			testHost, testHost + restPort, testHost + rpcPort, testHost + webPort,
		},
		{
			host2, "", "", "",
			host2, "", "", "",
		},
		{
			host2, testHost + restPort, testHost + rpcPort, testHost + webPort,
			host2, testHost + restPort, testHost + rpcPort, testHost + webPort,
		},
		{
			host2, def.RESTAddress, def.RPCAddress, def.RPCWebAddress,
			host2, host2 + ":7727", host2 + ":7721", host2 + ":7729",
		},
		{
			host2, testHost + restPort, def.RPCAddress, def.RPCWebAddress,
			host2, testHost + restPort, host2 + ":7721", host2 + ":7729",
		},
		{
			host2, testHost + restPort, testHost + restPort, testHost + restPort,
			host2, testHost + restPort, testHost + restPort, testHost + restPort,
		},
		{
			def.Host, def.RESTAddress, testHost + restPort, testHost + restPort,
			def.Host, def.RESTAddress, testHost + restPort, testHost + restPort,
		},
	}
	for _, test := range netTests {
		n := Network{}
		n.Host = test.uHost
		n.RESTAddress = test.uRest
		n.RPCAddress = test.uRPC
		n.RPCWebAddress = test.uRPCWeb
		n.AdminAddress = test.uRPC
		n.HealthAddress = test.uRPC
		err := n.normalize(serviceDefaults)
		assert.NoError(t, err)
		assert.Equal(t, test.nHost, n.Host)
		assert.Equal(t, test.nRest, n.RESTAddress)
		assert.Equal(t, test.nRPC, n.RPCAddress)
		assert.Equal(t, test.nRPCWeb, n.RPCWebAddress)
		assert.Equal(t, test.nRPC, n.AdminAddress)
		assert.Equal(t, test.nRPC, n.HealthAddress)
	}
	n := Network{}
	networkDefaults.HealthAddress = "::"
	networkDefaults.RPCWebAddress = "::"
	networkDefaults.RPCAddress = "::"
	networkDefaults.AdminAddress = "::"
	networkDefaults.RESTAddress = "::"
	n.Host = "::"
	n.HealthAddress = "::"
	err := n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.RPCWebAddress = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.RPCAddress = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.AdminAddress = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
	n.RESTAddress = "::"
	err = n.normalize(serviceDefaults)
	assert.Error(t, err)
}
