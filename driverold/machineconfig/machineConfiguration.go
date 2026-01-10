package machineconfig

import (
	"fmt"
	"net"

	"github.com/gwenya/qemu-driver/driverold/configstorage"
)

type MachineConfiguration struct {
	FirmwareSourcePath string
	NvramSourcePath    string
	CpuCount           uint32
	MemorySize         uint64
	DiskSize           uint64
	NetworkInterfaces  []NetworkInterface
	Volumes            []Volume
	VsockCid           uint32
	CloudInit          CloudInitData
}

func ConfigToStorage(m MachineConfiguration) configstorage.MachineConfiguration {
	var networkInterfaces []configstorage.NetworkInterface
	for i := range m.NetworkInterfaces {
		var mapped configstorage.NetworkInterface
		switch ni := m.NetworkInterfaces[i].(type) {
		case physicalNetworkInterface:
			mapped.Type = "physical"
			mapped.Physical = &configstorage.PhysicalNetworkInterfaceOptions{}

		case tapNetworkInterface:
			mapped.Type = "tap"
			mapped.Tap = &configstorage.TapNetworkInterfaceOptions{
				Name:       ni.name,
				MacAddress: ni.macAddress.String(),
			}

		default:
			panic("unimplemented storage mapping")
		}

		networkInterfaces = append(networkInterfaces, mapped)
	}

	var volumes []configstorage.Volume
	for i := range m.Volumes {
		var mapped configstorage.Volume
		switch v := m.Volumes[i].(type) {
		case cephVolume:
			mapped.Type = "rbd"
			mapped.Rbd = &configstorage.RbdVolumeOptions{
				Serial: v.Serial,
				Pool:   v.Pool,
				Name:   v.Name,
			}

		default:
			panic("unimplemented volume mapping")
		}

		volumes = append(volumes, mapped)
	}

	return configstorage.MachineConfiguration{
		FirmwareSourcePath: m.FirmwareSourcePath,
		NvramSourcePath:    m.NvramSourcePath,
		CpuCount:           m.CpuCount,
		MemorySize:         m.MemorySize,
		DiskSize:           m.DiskSize,
		NetworkInterfaces:  networkInterfaces,
		Volumes:            volumes,
		VsockCid:           m.VsockCid,
		CloudInit: configstorage.CloudInit{
			Meta:    m.CloudInit.Meta,
			Vendor:  m.CloudInit.Vendor,
			User:    m.CloudInit.User,
			Network: m.CloudInit.Network,
		},
	}
}

func MapConfigFromStorage(m configstorage.MachineConfiguration) (MachineConfiguration, error) {
	var networkInterfaces []NetworkInterface
	for _, ni := range m.NetworkInterfaces {
		var mapped NetworkInterface
		switch ni.Type {
		case "tap":
			macAddress, err := net.ParseMAC(ni.Tap.MacAddress)
			if err != nil {
				return MachineConfiguration{}, fmt.Errorf("parsing mac for tap network interface %s", ni.Tap.Name)
			}

			mapped = NewTapNetworkInterface(ni.Tap.Name, macAddress)

		case "physical":
			mapped = NewPhysicalNetworkInterface()

		default:
			return MachineConfiguration{}, fmt.Errorf("unimplemented network mapping: %s", ni.Type)
		}

		networkInterfaces = append(networkInterfaces, mapped)
	}

	var volumes []Volume
	for _, v := range m.Volumes {
		var mapped Volume
		switch v.Type {
		case "rbd":
			mapped = NewCephVolume(v.Rbd.Serial, v.Rbd.Pool, v.Rbd.Name)

		default:
			return MachineConfiguration{}, fmt.Errorf("unimplemented volume mapping: %s", v.Type)
		}

		volumes = append(volumes, mapped)
	}

	return MachineConfiguration{
		FirmwareSourcePath: m.FirmwareSourcePath,
		NvramSourcePath:    m.NvramSourcePath,
		CpuCount:           m.CpuCount,
		MemorySize:         m.MemorySize,
		DiskSize:           m.DiskSize,
		NetworkInterfaces:  networkInterfaces,
		Volumes:            volumes,
		VsockCid:           m.VsockCid,
		CloudInit: CloudInitData{
			Meta:    m.CloudInit.Meta,
			Vendor:  m.CloudInit.Vendor,
			User:    m.CloudInit.User,
			Network: m.CloudInit.Network,
		},
	}, nil
}
