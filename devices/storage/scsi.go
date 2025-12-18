package storage

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type ScsiDrive interface {
	Config() []config.Section
	GetScsiHotplug(bus string) devices.HotplugDevice
}

type scsiHotpluggable interface {
	Plug(m qmp.Monitor, bus string) error
	Unplug(m qmp.Monitor, bus string) error
}

func wrapScsiHotplug(device scsiHotpluggable, bus string) devices.HotplugDevice {
	return &scsiHotplugWrap{
		device: device,
		bus:    bus,
	}
}

type scsiHotplugWrap struct {
	device scsiHotpluggable
	bus    string
}

func (s *scsiHotplugWrap) Plug(m qmp.Monitor) error {
	return s.device.Plug(m, s.bus)
}

func (s *scsiHotplugWrap) Unplug(m qmp.Monitor) error {
	return s.device.Unplug(m, s.bus)
}
