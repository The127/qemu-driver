package driver

import "net"

type NetworkAdapter struct {
	name    string
	options any
}

type tapNetworkAdapterOptions struct {
	macAddress net.HardwareAddr
}

func NewTapNetworkAdapter(name string, macAddress net.HardwareAddr) NetworkAdapter {
	return NetworkAdapter{
		name: name,
		options: tapNetworkAdapterOptions{
			macAddress: macAddress,
		},
	}
}
