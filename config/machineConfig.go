package config

import (
	"fmt"
	"strings"
)

type MachineConfig struct {
	sections []Section
}

type Section struct {
	Name    string
	Entries map[string]string
}

func (c *MachineConfig) AddSection(section Section) {
	c.sections = append(c.sections, section)
}

func (c *MachineConfig) ToString() string {
	b := strings.Builder{}

	for _, section := range c.sections {
		b.WriteString(fmt.Sprintf("[%s]\n", section.Name))
		for key, value := range section.Entries {
			b.WriteString(fmt.Sprintf(`%s = "%s"`, key, value))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
