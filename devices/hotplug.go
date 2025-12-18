package devices

import "github.com/gwenya/qemu-driver/qmp"

type HotplugDevice interface {
	Plug(m qmp.Monitor) error
	Unplug(m qmp.Monitor) error
}
