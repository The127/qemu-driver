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
	return &allocator{}
}

type allocator struct {
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

func newRootPort(id string, index int) BusDevice {
	return &rootPort{
		id:    id,
		index: index,
	}
}

func (r *allocator) Allocate() []Allocation {
	// for now we don't care about the order or how much space they take up
	// but later we want to allocate multiple devices into one port via multi-function devices to save limited pcie space
	// and also we maybe to sort them in some way so that network devices always get the same address and therefore same predictable interface cid
	// or at least the main network device should be like that

	allocations := make([]Allocation, 0, len(r.Devices))

	port := 0
	fn := 0
	bridgeDev := 1 // address 0 is used by the DRAM controller
	bridgeFn := 0

	for _, device := range r.Devices {
		portName := fmt.Sprintf("qemu_pcie%d", port)
		if fn == 0 {
			allocations = append(allocations, Allocation{
				Device:        newRootPort(portName, port),
				Bus:           "pcie.0",
				Addr:          fmt.Sprintf("%x.%d", bridgeDev, bridgeFn),
				Multifunction: bridgeFn == 0,
			})

			if bridgeFn == 7 {
				bridgeDev++
				bridgeFn = 0
			} else {
				bridgeFn++
			}

		}

		allocations = append(allocations, Allocation{
			Device:        device,
			Bus:           portName,
			Addr:          fmt.Sprintf("00.%d", fn),
			Multifunction: fn == 0,
		})

		if fn == 7 {
			fn = 0
		} else {
			fn++
		}
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
