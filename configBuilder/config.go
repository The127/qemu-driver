package configBuilder

import (
	"fmt"
)

type ConfigBuilder interface {
}

type BusBuilder interface {
}

type configBuilder struct {
	sections   []configSection
	busDevices []BusDevice
}

type configSection struct {
	name    string
	entries map[string]string
}

type busDeviceConfig struct {
	name        string
	driver      string
	extraConfig map[string]string
}

func New() ConfigBuilder {
	return &configBuilder{}
}

func (c *configBuilder) CpuMemConfig() {
	cpuSockets := 1
	cpuCores := 1
	cpuThreads := 1
	cpuCount := 1
	memoryBytes := 1024 * 1024 * 1024

	c.sections = append(c.sections, configSection{
		name: "smp-opts",
		entries: map[string]string{
			"sockets": fmt.Sprintf("%d", cpuSockets),
			"cores":   fmt.Sprintf("%d", cpuCores),
			"threads": fmt.Sprintf("%d", cpuThreads),
			"cpus":    fmt.Sprintf("%d", cpuCount),
		},
	}, configSection{
		name: `object "mem0"`,
		entries: map[string]string{
			"qom-type": "memory-backend-memfd",
			"size":     fmt.Sprintf("%d", memoryBytes),
		},
	}, configSection{
		name: "numa",
		entries: map[string]string{
			"type":   "node",
			"nodeid": "0",
			"memdev": "mem0",
			"share":  "on",
		},
	}, configSection{
		name: "memory",
		entries: map[string]string{
			"size":   fmt.Sprintf("%d", memoryBytes),
			"maxmem": fmt.Sprintf("%d", memoryBytes),
		},
	})
}

func (c *configBuilder) MachineConfig() {
	c.sections = append(c.sections, configSection{
		name: "machine",
		entries: map[string]string{
			"graphics":       "off",
			"type":           "q35",
			"gic-version":    "",
			"cap-large-decr": "",
			"accel":          "kvm",
			"usb":            "off",
		},
	}, configSection{
		name: "global",
		entries: map[string]string{
			"driver":   "ICH9-LPC",
			"property": "disable_s3",
			"value":    "1",
		},
	}, configSection{
		name: "global",
		entries: map[string]string{
			"driver":   "ICH9-LPC",
			"property": "disable_s4",
			"value":    "0",
		},
	}, configSection{
		name:    "boot-opts",
		entries: map[string]string{"strict": "on"},
	})
}

func (c *configBuilder) UefiConfig(firmwarePath string, nvramPath string) {
	c.sections = append(c.sections, configSection{
		name: "drive",
		entries: map[string]string{
			"file":     firmwarePath,
			"if":       "pflash",
			"format":   "raw",
			"unit":     "0",
			"readonly": "on",
		},
	}, configSection{
		name: "drive",
		entries: map[string]string{
			"file":   nvramPath,
			"if":     "pflash",
			"format": "raw",
			"unit":   "1",
		},
	})
}

func (c *configBuilder) QmpConfig(socketPath string) {
	c.sections = append(c.sections, configSection{
		name: `chardev "monitor"`,
		entries: map[string]string{
			"backend": "socket",
			"path":    socketPath,
			"server":  "on",
			"wait":    "off",
		},
	}, configSection{
		name: "mon",
		entries: map[string]string{
			"chardev": "monitor",
			"mode":    "control",
		},
	})
}

func (c *configBuilder) ConsoleConfig() {
	c.sections = append(c.sections, configSection{
		name: `chardev "console"`,
		entries: map[string]string{
			"backend": "ringbuf",
			"size":    "1048576",
		},
	})
}

func (c *configBuilder) QemuCoreInfoConfig() {
	c.sections = append(c.sections, configSection{
		name: `device "qemu_vmcoreinfo"`,
		entries: map[string]string{
			"driver": "vmcoreinfo",
		},
	})
}

func (c *configBuilder) SerialConfig() {

	c.sections = append(c.sections, configSection{
		name: `chardev "qemu_spice-chardev"`,
		entries: map[string]string{
			"backend": "spicevmc",
			"name":    "vdagent",
		},
	}, configSection{
		name: `device "qemu_spice"`,
		entries: map[string]string{
			"driver":  "virtserialport",
			"name":    "com.redhat.spice.0",
			"chardev": "qemu_spice-chardev",
			"bus":     "dev-qemu_serial.0",
		},
	}, configSection{
		name: `chardev "qemu_spicedir-chardev"`,
		entries: map[string]string{
			"backend": "spiceport",
			"name":    "org.spice-space.webdav.0",
		},
	}, configSection{
		name: `device "qemu_spicedir"`,
		entries: map[string]string{
			"driver":  "virtserialport",
			"name":    "org.spice-space.webdav.0",
			"chardev": "qemu_spicedir-chardev",
			"bus":     "dev-qemu_serial.0",
		},
	})
}
