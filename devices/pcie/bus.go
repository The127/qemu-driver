package pcie

import (
	"fmt"

	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
)

type Bus interface {
	AddDevice(device BusDevice)
}

func NewAllocator() Allocator {
	return &busAllocator{}
}

type busAllocator struct {
	devices.BusImpl[BusDevice]
}

type rootPort struct {
	noHotPlug
	id    string
	index int
}

func (p *rootPort) Config(alloc BusAllocation) []config.Section {
	return []config.Section{
		busDeviceConfigSection(alloc, p.id, "pcie-root-port", map[string]string{
			"chassis": fmt.Sprintf("%d", p.index),
		}),
	}
}

type RootPort interface{}

func (r *busAllocator) Allocate() []Allocation {
	hotplugDevices := make([]BusDevice, 0, len(r.Devices))
	noHotplugDevices := make([]BusDevice, 0, len(r.Devices))

	for _, device := range r.Devices {
		if device.IsHotplug() {
			hotplugDevices = append(hotplugDevices, device)
		} else {
			noHotplugDevices = append(noHotplugDevices, device)
		}
	}

	rootBus := "pcie.0"

	allocations := make([]Allocation, 0, len(r.Devices))

	i := 64 // first pcie slot (8 functions * 8 functions) is used by DRAM controller

	for _, device := range noHotplugDevices {
		portNum := i / 8
		portName := fmt.Sprintf("pcie_port_%d", portNum)
		deviceFn := i % 8

		if deviceFn == 0 {
			port := portNum / 8
			portFn := portNum % 8
			allocations = append(allocations, Allocation{
				Device: &rootPort{
					id:    portName,
					index: portNum,
				},
				Bus:           rootBus,
				Addr:          fmt.Sprintf("%x.%d", port, portFn),
				Multifunction: portFn == 0,
			})
		}

		allocations = append(allocations, Allocation{
			Device:        device,
			Bus:           portName,
			Addr:          fmt.Sprintf("00.%d", deviceFn),
			Multifunction: deviceFn == 0,
		})

		i += 1
	}

	if i%8 != 0 { // if
		i += 8 - (i % 8)
	}

	for _, device := range hotplugDevices {
		portNum := i / 8
		portName := fmt.Sprintf("pcie_port_%d", portNum)

		port := portNum / 8
		portFn := portNum % 8
		allocations = append(allocations, Allocation{
			Device: &rootPort{
				id:    portName,
				index: portNum,
			},
			Bus:           rootBus,
			Addr:          fmt.Sprintf("%x.%d", port, portFn),
			Multifunction: portFn == 0,
		})

		allocations = append(allocations, Allocation{
			Device:        device,
			Bus:           portName,
			Addr:          "00.0",
			Multifunction: false,
		})

		i += 8
	}

	// add some root ports for hot-plugging additional devices, the amount should be configurable, maybe different way of doing it
	for range 8 {
		portNum := i / 8
		portName := fmt.Sprintf("pcie_port_%d", portNum)

		port := portNum / 8
		portFn := portNum % 8
		allocations = append(allocations, Allocation{
			Device: &rootPort{
				id:    portName,
				index: portNum,
			},
			Bus:           rootBus,
			Addr:          fmt.Sprintf("%x.%d", port, portFn),
			Multifunction: portFn == 0,
		})

		i += 8
	}

	return allocations
}

type Allocator interface {
	Bus
	Allocate() []Allocation
}

type Allocation struct {
	Device        BusDevice
	Bus           string
	Addr          string
	Multifunction bool
}
