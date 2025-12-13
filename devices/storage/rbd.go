package storage

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type rbdDrive struct{}

func (r rbdDrive) Config() []config.Section {
	//TODO implement me
	panic("implement me")
}

func (r rbdDrive) GetHotplugs() []devices.HotplugDevice {
	//TODO implement me
	panic("implement me")
}

func (r rbdDrive) Plug(m qmp.Monitor) error {
	//TODO implement me
	panic("implement me")
}

func (r rbdDrive) Unplug(m qmp.Monitor) error {
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
