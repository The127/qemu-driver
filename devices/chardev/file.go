package chardev

import (
	"github.com/gwenya/qemu-driver/config"
)

type fileChardev struct {
	id        string
	path      string
	inputPath string
}

func NewFile(id string, path string) Chardev {
	return &fileChardev{
		id:   id,
		path: path,
	}
}

func NewFileWithInput(id string, outputPath string, inputPath string) Chardev {
	return &fileChardev{
		id:        id,
		path:      outputPath,
		inputPath: inputPath,
	}
}

func (c *fileChardev) Config() config.Section {
	options := map[string]string{
		"path": c.path,
	}

	if c.inputPath != "" {
		options["input-path"] = c.inputPath
	}

	return chardevConfig(c.id, "file", options)
}
