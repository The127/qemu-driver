package chardev

import (
	"github.com/gwenya/qemu-driver/config"
)

type spicevmcChardev struct {
	id      string
	channel string
}

func NewSpicevmc(id string, channel string) Chardev {
	return &spicevmcChardev{
		id:      id,
		channel: channel,
	}
}

func (c *spicevmcChardev) Config() config.Section {
	return chardevConfig(c.id, "spicevmc", map[string]string{
		"name": c.channel,
	})
}
