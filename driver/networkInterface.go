package driver

import (
	"net"
)

type NetworkInterface interface {
	__networkInterfaceMarker()
}

type tapNetworkInterface struct {
	macAddress net.HardwareAddr
	name       string
}

func (tapNetworkInterface) __networkInterfaceMarker() {}

func NewTapNetworkInterface(name string, macAddress net.HardwareAddr) NetworkInterface {
	return tapNetworkInterface{
		macAddress: macAddress,
		name:       name,
	}
}

type physicalNetworkInterface struct {
}

func (physicalNetworkInterface) __networkInterfaceMarker() {}
