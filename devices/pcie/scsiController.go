package pcie

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/devices/storage"
	"github.com/gwenya/qemu-driver/qmp"
)

type scsiController struct {
	id    string
	disks []storage.ScsiDrive
}

type ScsiController interface {
	BusDevice
	AddDisk(disk storage.ScsiDrive)
}

func NewScsiController(id string) ScsiController {
	return &scsiController{
		id: id,
	}
}

func (s *scsiController) Config(alloc BusAllocation) []config.Section {
	// TODO: devices?
	return []config.Section{
		busDeviceConfigSection(alloc, s.id, "virtio-scsi-pci", nil),
	}
}

func (s *scsiController) GetHotplugs() []devices.HotplugDevice {
	var devs []devices.HotplugDevice
	for _, disk := range s.disks {
		devs = append(devs, disk.GetHotplugs()...)
	}

	return devs
}

func (s *scsiController) Plug(q qmp.Monitor) error {
	panic("not supported")
}

func (*scsiController) Unplug(q qmp.Monitor) error {
	panic("not supported")
}

func (s *scsiController) AddDisk(disk storage.ScsiDrive) {
	s.disks = append(s.disks, disk)
}
