package devices

import "github.com/gwenya/qemu-driver/qmp"

type HotplugDevice interface {
	GetHotplugs() []HotplugDevice
	Plug(m qmp.Monitor) error
	Unplug(m qmp.Monitor) error
}
