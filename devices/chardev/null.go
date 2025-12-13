package chardev

import (
	"github.com/gwenya/qemu-driver/config"
)

type nullChardev struct {
	id string
}

func NewNull(id string) Chardev {
	return &nullChardev{
		id: id,
	}
}

func (c *nullChardev) Config() config.Section {
	return chardevConfig(c.id, "null", nil)
}
