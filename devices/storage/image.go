package storage

import (
	"fmt"
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type imageDrive struct {
	id   string
	path string
}

type ImageDrive interface {
	ScsiDrive
	BlkDrive
}

func NewImageDrive(id string, path string) ImageDrive {
	return &imageDrive{
		id:   id,
		path: path,
	}
}

func (d *imageDrive) Config() []config.Section {
	return nil
}

func (d *imageDrive) GetHotplugs() []devices.HotplugDevice {
	return []devices.HotplugDevice{d}
}

func (d *imageDrive) Plug(m qmp.Monitor) error {
	nodeName := "node-" + d.id

	err := m.AddBlockDevice(map[string]any{
		"aio": "io_uring",
		"cache": map[string]any{
			"direct":   true,
			"no-flush": false,
		},
		"discard":   "unmap", // Forward as an unmap request. This is the same as `discard=on` in the qemu config file.
		"driver":    "file",
		"node-name": nodeName,
		"read-only": false,
		"locking":   "off",
		"filename":  d.path,
	})

	if err != nil {
		return fmt.Errorf("adding block device: %w", err)
	}

	err = m.AddDevice(map[string]any{
		"id":        d.id,
		"drive":     nodeName,
		"serial":    d.id,
		"device_id": nodeName,
		"channel":   0,
		"lun":       1,
		"bus":       "scsi.0", // this needs to be passed in from the scsi bus somehow => make a wrapper type for hotplugging
		"driver":    "scsi-hd",
	})

	if err != nil {
		return fmt.Errorf("adding device: %w", err)
	}

	return nil
}

func (d *imageDrive) Unplug(m qmp.Monitor) error {
	panic("implement me")
}
