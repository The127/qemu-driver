package storage

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type rbdDrive struct{}

func (r *rbdDrive) Config() []config.Section {
	//TODO implement me
	panic("implement me")
}

func (r *rbdDrive) GetScsiHotplug(bus string) devices.HotplugDevice {
	return wrapScsiHotplug(r, bus)
}

func (r *rbdDrive) Plug(m qmp.Monitor, bus string) error {
	//TODO implement me
	panic("implement me")
}

func (r *rbdDrive) Unplug(m qmp.Monitor, bus string) error {
	//TODO implement me
	panic("implement me")
}

type RbdDrive interface {
	ScsiDrive
	BlkDrive
}

func NewRbdDrive() RbdDrive {
	return &rbdDrive{}
}
