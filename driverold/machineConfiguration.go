package driver

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kdomanski/iso9660"
)

func (m *MachineConfiguration) UnmarshalJSON(bytes []byte) error {
	var data map[string]any
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	var ok bool

	m.FirmwareSourcePath, ok = data["firmware-source-path"]

}

type networkInterfaceConfig struct {
	Type string
	Tap  *tapNetworkInterface
}

func (d *CloudInitData) CreateIso(path string) error {
	writer, err := iso9660.NewWriter()
	if err != nil {
		return fmt.Errorf("creating iso writer: %w", err)
	}

	defer writer.Cleanup()

	err = writer.AddFile(strings.NewReader(d.User), "user-data")
	if err != nil {
		return fmt.Errorf("adding user-data: %w", err)
	}

	err = writer.AddFile(strings.NewReader(d.Meta), "meta-data")
	if err != nil {
		return fmt.Errorf("adding meta-data: %w", err)
	}

	if d.Vendor != "" {
		err = writer.AddFile(strings.NewReader(d.Vendor), "vendor-data")
		if err != nil {
			return fmt.Errorf("adding vendor-data: %w", err)
		}
	}

	if d.Network != "" {
		err = writer.AddFile(strings.NewReader(d.Network), "network-config")
		if err != nil {
			return fmt.Errorf("adding network-config: %w", err)
		}
	}

	isoFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer isoFile.Close()

	err = writer.WriteTo(isoFile, "cidata")
	if err != nil {
		return fmt.Errorf("writing iso file: %w", err)
	}

	return nil
}
