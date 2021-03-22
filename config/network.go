package config

import "net"

// HealthCheck is the health check route to use.
const HealthCheck = "/health"

// Network config
type Network struct {
	// Host is default adapter to listen on.
	// default: localhost
	Host string `json:"host"`
	// Health is the address for the health check server.
	// default: [Host]:7720
	Health string `json:"health"`
	// RPC is the address for the gRPC server.
	// default: [Host]:7721
	RPC string `json:"rpc"`
	// Admin is the address for the admin server.
	// default: [Host]:7722
	Admin string `json:"admin"`
	// HTTP is the address for the HTTP server.
	// default: [Host]:7727
	REST string `json:"rest"`
	// RPCWeb is the address for the gRPC-Web server.
	// default: [Host]:7729
	RPCWeb string `json:"rpcweb" mapstructure:"rpcweb"`

	// TODO: use RequestID
	// RequestID is the request id to use
	RequestID string `json:"request_id" yaml:"request_id" mapstructure:"request_id"`
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
	if n.RPC == dc.RPC {
		n.RPC, err = updateHost(n.RPC, n.Host)
		if err != nil {
			return
		}
	}
	if n.Admin == dc.Admin {
		n.Admin, err = updateHost(n.Admin, n.Host)
		if err != nil {
			return
		}
	}
	if n.REST == dc.REST {
		n.REST, err = updateHost(n.REST, n.Host)
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
