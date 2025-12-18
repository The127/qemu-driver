package main

import (
	"os"
	"path"

	"github.com/gwenya/qemu-driver/driver"

	"github.com/google/uuid"
)

func main() {
	id := uuid.MustParse("804859f4-343b-4a0f-97b1-75d04aee531d")

	storagePath := path.Join("/tmp", id.String())

	imageSource := "/var/home/gwen/Downloads/Arch-Linux-x86_64-cloudimg-20251201.460866.qcow2"
	firmwareSource := "/usr/share/edk2/ovmf/OVMF_CODE.fd"
	nvramSource := "/usr/share/edk2/ovmf/OVMF_VARS.fd"

	err := os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	d, err := driver.New(driver.MachineConfiguration{
		Id:                 id,
		StoragePath:        storagePath,
		ImageSourcePath:    imageSource,
		FirmwareSourcePath: firmwareSource,
		NvramSourcePath:    nvramSource,
		CpuCount:           1,
		MemorySize:         1024,
		DiskSize:           10,
		NetworkInterfaces:  nil,
		Volumes:            nil,
	})

	if err != nil {
		panic(err)
	}

	_ = d

	err = d.Start()
	if err != nil {
		panic(err)
	}
}
