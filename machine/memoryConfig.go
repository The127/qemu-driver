package machine

import (
	"fmt"

	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type memoryConfig struct {
	sizeMB    int
	maxSizeMB int
	slots     int
}

func (c *memoryConfig) Config() config.Section {
	section := config.Section{
		Name: "memory",
		Entries: map[string]string{
			"size":   "0",
			"maxmem": fmt.Sprintf("%dM", c.maxSizeMB),
		},
	}

	section.Entries["slots"] = fmt.Sprintf("%d", c.slots)

	return section
}

type memHotplug struct {
	id       string
	memdevId string
	size     uint64
}

func (h *memHotplug) Plug(m qmp.Monitor) error {
	err := m.AddMemoryBackend(h.memdevId, h.size)
	if err != nil {
		return err
	}

	err = m.AddDevice(map[string]any{
		"driver": "pc-dimm",
		"memdev": h.memdevId,
		"id":     h.id,
	})

	if err != nil {
		return err
	}

	return nil
}

func (h *memHotplug) Unplug(m qmp.Monitor) error {
	//TODO implement me
	panic("implement me")
}

func (c *memoryConfig) GetHotplugs() []devices.HotplugDevice {
	return []devices.HotplugDevice{
		&memHotplug{
			id:       "dimm0",
			memdevId: "mem0",
			size:     uint64(c.sizeMB) * 1024 * 1024,
		},
	}
}
