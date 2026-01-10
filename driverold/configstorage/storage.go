package configstorage

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteToFile(config MachineConfiguration, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("opening %q for writing: %w", path, err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	encoder := json.NewEncoder(f)

	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("serializing machine configuration: %w", err)
	}

	return nil
}

func ReadFromFile(path string) (MachineConfiguration, error) {
	f, err := os.Open(path)
	if err != nil {
		return MachineConfiguration{}, fmt.Errorf("opening %q for reading: %w", path, err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	decoder := json.NewDecoder(f)

	var result MachineConfiguration

	err = decoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("deserializing machine configuration: %w", err)
	}

	return result, nil
}

type MachineConfiguration struct {
	FirmwareSourcePath string             `json:"firmware-source-path"`
	NvramSourcePath    string             `json:"nvram-source-path"`
	CpuCount           uint32             `json:"cpu-count"`
	MemorySize         uint64             `json:"memory-size"`
	DiskSize           uint64             `json:"disk-size"`
	NetworkInterfaces  []NetworkInterface `json:"network-interfaces"`
	Volumes            []Volume           `json:"volumes"`
	VsockCid           uint32             `json:"vsock-cid"`
	CloudInit          CloudInit          `json:"cloud-init"`
}

type NetworkInterface struct {
	Type     string                           `json:"type"`
	Tap      *TapNetworkInterfaceOptions      `json:"tap"`
	Physical *PhysicalNetworkInterfaceOptions `json:"physical"`
}

type TapNetworkInterfaceOptions struct {
	Name       string `json:"name"`
	MacAddress string `json:"mac-address"`
}

type PhysicalNetworkInterfaceOptions struct {
	Name string `json:"name"`
}

type Volume struct {
	Type string            `json:"type"`
	Rbd  *RbdVolumeOptions `json:"rbd"`
}

type RbdVolumeOptions struct {
	Serial string `json:"serial"`
	Pool   string `json:"pool"`
	Name   string `json:"name"`
}

type CloudInit struct {
	Meta    string `json:"meta"`
	Network string `json:"network"`
	Vendor  string `json:"vendor"`
	User    string `json:"user"`
}
