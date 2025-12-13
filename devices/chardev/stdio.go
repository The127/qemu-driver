package chardev

import (
	"github.com/gwenya/qemu-driver/config"
	"github.com/gwenya/qemu-driver/util"
)

type stdioChardev struct {
	id     string
	signal bool
}

func NewStdio(id string, enableSignals bool) Chardev {
	return &stdioChardev{
		id:     id,
		signal: enableSignals,
	}
}

func (c *stdioChardev) Config() config.Section {
	return chardevConfig(c.id, "stdio", map[string]string{
		"signal": util.BoolToOnOff(c.signal),
	})
}
