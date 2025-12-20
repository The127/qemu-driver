package pcie

import (
	"fmt"
	"maps"

	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/qmp"
)

type BusAllocation struct {
	Bus           string
	Address       string
	Multifunction bool
}

type BusDevice interface {
	IsHotplug() bool
	Config(alloc BusAllocation) []config.Section
	GetHotplugs(alloc BusAllocation) []devices.HotplugDevice
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

func (*noHotPlug) IsHotplug() bool {
	return false
}

func (*noHotPlug) GetHotplugs(_ BusAllocation) []devices.HotplugDevice {
	return nil
}

type hotPlug struct{}

func (*hotPlug) IsHotplug() bool {
	return true
}

type wrappedHotpluggable interface {
	Plug(m qmp.Monitor, alloc BusAllocation) error
	Unplug(m qmp.Monitor, alloc BusAllocation) error
}

func hotplugWrap(device wrappedHotpluggable, allocation BusAllocation) devices.HotplugDevice {
	return &hotplugWrapper{
		wrapped:    device,
		allocation: allocation,
	}
}

type hotplugWrapper struct {
	wrapped    wrappedHotpluggable
	allocation BusAllocation
}

func (h *hotplugWrapper) Plug(m qmp.Monitor) error {
	return h.wrapped.Plug(m, h.allocation)
}

func (h *hotplugWrapper) Unplug(m qmp.Monitor) error {
	return h.wrapped.Unplug(m, h.allocation)
}
