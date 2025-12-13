package pcie

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type physicalNetworkDevice struct {
	id string
}

func NewPhysicalNetworkDevice(id string, netdevName string) BusDevice {
	return &physicalNetworkDevice{
		id: id,
	}
}

func (d *physicalNetworkDevice) Config(BusAllocation) []config.Section {
	return nil
}

func (d *physicalNetworkDevice) GetHotplugs() []devices.HotplugDevice {
	return []devices.HotplugDevice{
		d,
	}
}

func (d *physicalNetworkDevice) Plug(m qmp.Monitor) error {
	//TODO implement me
	panic("implement me")
}

func (d *physicalNetworkDevice) Unplug(m qmp.Monitor) error {
	//TODO implement me
	panic("implement me")
}
