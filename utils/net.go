package utils

import (
	"fmt"
	"net"
)

var dns = "8.8.8.8:80"

// OutboundIP preferred outbound ip of this machine
func OutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", dns)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// MakeAddress returns a host:port address
func MakeAddress(host string, port int) string {
	return fmt.Sprint(host, ":", port)
}
