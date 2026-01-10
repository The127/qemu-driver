package machine

import (
	"fmt"

	"github.com/gwenya/qemu-driver/config"
)

type cpuConfig struct {
	cpus    int
	maxCpus int
}

func (c *cpuConfig) Config() config.Section {
	return config.Section{
		Name: "smp-opts",
		Entries: map[string]string{
			"cpus":    fmt.Sprintf("%d", c.cpus),
			"maxcpus": fmt.Sprintf("%d", c.maxCpus),
		},
	}
}
