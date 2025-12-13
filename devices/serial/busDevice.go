package serial

import "github.com/gwenya/qemu-driver/config"

type BusDevice interface {
	Config(busName string) []config.Section
}
