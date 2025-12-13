package chardev

import (
	"github.com/gwenya/qemu-driver/config"
)

type spicePortChardev struct {
	id      string
	channel string
}

func NewSpiceport(id string, channel string) Chardev {
	return &spicePortChardev{
		id:      id,
		channel: channel,
	}
}

func (c *spicePortChardev) Config() config.Section {
	return chardevConfig(c.id, "spiceport", map[string]string{
		"name": c.channel,
	})
}
