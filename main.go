package main

import (
	"fmt"
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

	//hwaddr := net.HardwareAddr{
	//	(byte(rand.Int31n(256)) & ^byte(0b1)) | byte(0b10),
	//	byte(rand.Int31n(256)),
	//	byte(rand.Int31n(256)),
	//	byte(rand.Int31n(256)),
	//	byte(rand.Int31n(256)),
	//	byte(rand.Int31n(256)),
	//}

	userData := `#cloud-config
users:
  - name: root
    lock_passwd: false
    hashed_passwd: "$6$rounds=4096$nqxzsCUB62RiUjKp$YX0V8FDfz/9LdV6d5s0UGKBT8tAH2svzBILoS9Z/rWXjcny9Z9.ANt5XI6PU87268UrJrWeqtmH1lupgZtKZI/"

  - name: arch
    lock_passwd: false
    hashed_passwd: "$6$ekjadUze3yUXluSP$KTBA960c5FiFQIVHz7WQ8/9paorVjLWbnQ./NpUG8sE99ehX4SELqrMPEFq/yFKCB55i9gw6xbg.75i49WRlh/"

runcmd:
  - echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
`

	d, err := driver.New("/usr/bin/qemu-system-x86_64", driver.MachineConfiguration{
		Id:                 id,
		StorageDirectory:   storagePath,
		RuntimeDirectory:   storagePath,
		ImageSourcePath:    imageSource,
		FirmwareSourcePath: firmwareSource,
		NvramSourcePath:    nvramSource,
		CpuCount:           1,
		MemorySize:         1024 * 1024 * 1024,
		DiskSize:           50_000_000_000,
		NetworkInterfaces:  []driver.NetworkInterface{
			//driver.NewTapNetworkInterface("test-tap", hwaddr),
		},
		Volumes: []driver.Volume{
			driver.NewCephVolume("sample-volume", "dev", "33333333-549a-42d3-87c9-090451088b24"),
		},
		VsockCid: 4,
		CloudInit: driver.CloudInitData{
			Vendor:  fmt.Sprintf("instance-id: %s", id),
			Meta:    "",
			User:    string(userData),
			Network: "",
		},
	})

	if err != nil {
		panic(err)
	}

	err = d.Create()
	if err != nil {
		panic(err)
	}

	err = d.Start()
	if err != nil {
		panic(err)
	}
}
