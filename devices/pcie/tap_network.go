package pcie

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type tapNetworkDevice struct {
	id         string
	queueCount int
}

func NewTapNetworkDevice(id string, netdevName string, queueCount int) BusDevice {
	return &tapNetworkDevice{
		id:         id,
		queueCount: queueCount,
	}
}

func (d *tapNetworkDevice) Config(_ BusAllocation) []config.Section {
	return nil
}

func (d *tapNetworkDevice) GetHotplugs(alloc BusAllocation) []devices.HotplugDevice {
	return []devices.HotplugDevice{
		hotplugWrap(d, alloc),
	}
}

func (d *tapNetworkDevice) Plug(m qmp.Monitor, alloc BusAllocation) error {
	//TODO implement me
	panic("implement me")
}

func (d *tapNetworkDevice) Unplug(m qmp.Monitor, alloc BusAllocation) error {
	//TODO implement me
	panic("implement me")
}
