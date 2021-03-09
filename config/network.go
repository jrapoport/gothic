package config

import "net"

// HealthEndpoint is the rest health check route to use.
const HealthEndpoint = "/health"

// Network config
type Network struct {
	// Host is default adapter to listen on.
	// default: localhost
	Host string `json:"host"`
	// REST is the address for the REST server.
	// default: [Host]:8081
	REST string `json:"rest"`
	// RPC is the address for the gRPC server.
	// default: [Host]:3001
	RPC string `json:"rpc"`
	// RPCWeb is the address for the gRPC-Web server.
	// default: [Host]:6001
	RPCWeb string `json:"rpcweb" mapstructure:"rpcweb"`
	// Health is the address for the health check server.
	// default: [Host]:10001
	Health string `json:"health"`
}

func (n *Network) normalize(Service) (err error) {
	dc := networkDefaults
	if n.Host == dc.Host {
		return
	}
	updateHost := func(addr, host string) (string, error) {
		var p string
		_, p, err = net.SplitHostPort(addr)
		if err != nil {
			return "", err
		}
		return net.JoinHostPort(host, p), nil
	}
	if n.REST == dc.REST {
		n.REST, err = updateHost(n.REST, n.Host)
		if err != nil {
			return
		}
	}
	if n.RPC == dc.RPC {
		n.RPC, err = updateHost(n.RPC, n.Host)
		if err != nil {
			return
		}
	}
	if n.RPCWeb == dc.RPCWeb {
		n.RPCWeb, err = updateHost(n.RPCWeb, n.Host)
		if err != nil {
			return
		}
	}
	if n.Health == "" {
		n.Health = dc.Health
	}
	if n.Health == dc.Health {
		n.Health, err = updateHost(n.Health, n.Host)
		if err != nil {
			return
		}
	}
	return
}
