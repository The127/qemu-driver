package driver

import (
	"github.com/google/uuid"
	"github.com/gwenya/qemu-driver/pidfd"
)

type Options struct {
	SystemId         uuid.UUID
	StorageDirectory string
	RuntimeDirectory string

	QemuPath    string
	Logger      Logger
	PidFdWaiter pidfd.Waiter
}
