package chardev

import "github.com/gwenya/qemu-driver/config"

type pipeChardev struct {
	id   string
	path string
}

func NewPipe(id string, path string) Chardev {
	return &pipeChardev{
		id:   id,
		path: path,
	}
}

func (c *pipeChardev) Config() config.Section {
	return chardevConfig(c.id, "pipe", map[string]string{
		"path": c.path,
	})
}
