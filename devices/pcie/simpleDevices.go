package pcie

import "github.com/gwenya/qemu-driver/config"

type simpleDevice struct {
	noHotPlug
	id     string
	driver string
}

func (d *simpleDevice) Config(alloc BusAllocation) []config.Section {
	return []config.Section{
		busDeviceConfigSection(alloc, d.id, d.driver, nil),
	}
}

func NewBalloon(id string) BusDevice {
	return &simpleDevice{
		id:     id,
		driver: "virtio-balloon-pci",
	}
}

func NewKeyboard(id string) BusDevice {
	return &simpleDevice{
		id:     id,
		driver: "virtio-keyboard-pci",
	}
}

func NewTablet(id string) BusDevice {
	return &simpleDevice{
		id:     id,
		driver: "virtio-tablet-pci",
	}
}

type VgaDriver string

const (
	CirrusVga        VgaDriver = "cirrus-vga"
	StdVga           VgaDriver = "VGA"
	StdLegacyFreeVga VgaDriver = "secondary-vga"
	QxlVga           VgaDriver = "qxl-vga"
	VirtioVga        VgaDriver = "virtio-vga"
)

func NewVga(id string, driver VgaDriver) BusDevice {
	return &simpleDevice{
		id:     id,
		driver: string(driver),
	}
}
