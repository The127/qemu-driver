package machine

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/devices/chardev"
	"github.com/gwenya/qemu-driver/devices/pcie"
)

type Description struct {
	cpu            cpuConfig
	memory         memoryConfig
	monitorChardev string
	pcie           pcie.Allocator
	scsi           pcie.ScsiController
	serial         pcie.SerialBus
	chardevs       []chardev.Chardev
	firmwarePath   string
	nvramPath      string
}

func (d *Description) Firmware(firmwarePath string, nvramPath string) {
	d.firmwarePath = firmwarePath
	d.nvramPath = nvramPath
}

func (d *Description) Cpu(sockets int, cores int, threads int) {
	d.cpu = cpuConfig{
		sockets: sockets,
		cores:   cores,
		threads: threads,
	}
}

func (d *Description) Memory(sizeMB int, maxSizeMB int) {
	slots := 8
	if sizeMB == maxSizeMB {
		slots = 1
	}

	d.memory = memoryConfig{
		sizeMB:    sizeMB,
		maxSizeMB: maxSizeMB,
		slots:     slots,
	}
}

func (d *Description) Monitor(path string) {
	d.monitorChardev = "monitor"

	d.AddChardev(chardev.NewSocket(d.monitorChardev, chardev.SocketOpts{
		Unix: chardev.SocketOptsUnix{
			Path: path,
		},
		Server: true,
	}))
}

func (d *Description) AddChardev(device chardev.Chardev) {
	d.chardevs = append(d.chardevs, device)
}

func (d *Description) Pcie() pcie.Bus {
	if d.pcie == nil {
		d.pcie = pcie.NewAllocator()
	}
	return d.pcie
}

func (d *Description) Scsi() pcie.ScsiController {
	if d.scsi == nil {
		d.scsi = pcie.NewScsiController("scsi")
		d.Pcie().AddDevice(d.scsi)
	}

	return d.scsi
}

func (d *Description) Serial() pcie.SerialBus {
	if d.serial == nil {
		d.serial = pcie.NewSerialBus("serial")
		d.Pcie().AddDevice(d.serial)
	}

	return d.serial
}

func addBaseConfig(cfg *config.MachineConfig) {
	cfg.AddSection(config.Section{
		Name: "machine",
		Entries: map[string]string{
			"graphics": "off",
			"type":     "q35",
			//"gic-version":    "",
			//"cap-large-decr": "",
			"accel": "kvm",
			"usb":   "off",
		},
	})

	cfg.AddSection(config.Section{
		Name: "global",
		Entries: map[string]string{
			"driver":   "ICH9-LPC",
			"property": "disable_s3",
			"value":    "1",
		},
	})

	cfg.AddSection(config.Section{
		Name: "global",
		Entries: map[string]string{
			"driver":   "ICH9-LPC",
			"property": "disable_s4",
			"value":    "1",
		},
	})

	cfg.AddSection(config.Section{
		Name: "boot-opts",
		Entries: map[string]string{
			"strict": "on",
		},
	})
}

func (d *Description) BuildConfig() (config.MachineConfig, []devices.HotplugDevice) {
	cfg := config.MachineConfig{}
	var hotplugDevices []devices.HotplugDevice

	addBaseConfig(&cfg)

	if d.memory != (memoryConfig{}) {
		cfg.AddSection(d.memory.Config())
	}

	if d.cpu != (cpuConfig{}) {
		cfg.AddSection(d.cpu.Config())
	}

	if d.monitorChardev != "" {
		cfg.AddSection(config.Section{
			Name: "mon",
			Entries: map[string]string{
				"chardev": d.monitorChardev,
				"mode":    "control",
			},
		})
	}

	for _, device := range d.chardevs {
		cfg.AddSection(device.Config())
	}

	pcieAllocations := d.pcie.Allocate()
	for _, allocation := range pcieAllocations {
		pciAlloc := pcie.BusAllocation{
			Bus:           allocation.Bus,
			Address:       allocation.Addr,
			Multifunction: allocation.Multifunction,
		}

		sections := allocation.Device.Config(pciAlloc)

		for _, section := range sections {
			cfg.AddSection(section)
		}

		hotplugDevices = append(hotplugDevices, allocation.Device.GetHotplugs(pciAlloc)...)
	}

	if d.firmwarePath != "" {
		cfg.AddSection(config.Section{
			Name: "drive",
			Entries: map[string]string{
				"file":     d.firmwarePath,
				"if":       "pflash",
				"format":   "raw",
				"unit":     "0",
				"readonly": "on",
			},
		})
	}

	if d.nvramPath != "" {
		cfg.AddSection(config.Section{
			Name: "drive",
			Entries: map[string]string{
				"file":   d.nvramPath,
				"if":     "pflash",
				"format": "raw",
				"unit":   "1",
			},
		})
	}

	return cfg, hotplugDevices
}
