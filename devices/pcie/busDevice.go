package pcie

import (
	"fmt"
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
	"maps"
)

type BusAllocation struct {
	Bus           string
	Address       string
	Multifunction bool
}

type BusDevice interface {
	Config(alloc BusAllocation) []config.Section
	devices.HotplugDevice
}

func busDeviceConfigSection(alloc BusAllocation, id string, driver string, extraConfig map[string]string) config.Section {
	var entries map[string]string
	if extraConfig == nil {
		entries = make(map[string]string)
	} else {
		entries = maps.Clone(extraConfig)
	}

	entries["driver"] = driver
	entries["bus"] = alloc.Bus
	entries["addr"] = alloc.Address
	if alloc.Multifunction {
		entries["multifunction"] = "on"
	}

	return config.Section{
		Name:    fmt.Sprintf(`device "%s"`, id),
		Entries: entries,
	}
}

type noHotPlug struct{}

func (*noHotPlug) GetHotplugs() []devices.HotplugDevice {
	return nil
}

func (*noHotPlug) Plug(m qmp.Monitor) error {
	return NoHotplugError
}

func (*noHotPlug) Unplug(m qmp.Monitor) error {
	return NoHotplugError
}
