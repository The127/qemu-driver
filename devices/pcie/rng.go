package pcie

import "github.com/gwenya/qemu-driver/config"

type rngDevice struct {
	noHotPlug
	id  string
	rng string
}

func NewRng(id string, rng string) BusDevice {
	return &rngDevice{
		id:  id,
		rng: rng,
	}
}

func (d *rngDevice) Config(alloc BusAllocation) []config.Section {
	return []config.Section{
		busDeviceConfigSection(alloc, d.id, "virtio-rng-pci", map[string]string{
			"rng": d.rng,
		}),
	}
}
