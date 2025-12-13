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

func (d *tapNetworkDevice) Config(BusAllocation) []config.Section {
	return nil
}

func (d *tapNetworkDevice) GetHotplugs() []devices.HotplugDevice {
	return []devices.HotplugDevice{d}
}

func (d *tapNetworkDevice) Plug(m qmp.Monitor) error {
	//TODO implement me
	panic("implement me")
}

func (d *tapNetworkDevice) Unplug(m qmp.Monitor) error {
	//TODO implement me
	panic("implement me")
}
