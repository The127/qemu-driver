package storage

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
)

type ScsiDrive interface {
	Config() []config.Section
	devices.HotplugDevice
}
