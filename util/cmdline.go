package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func GetCmdline(pid int) ([]string, error) {
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdlineBytes, err := os.ReadFile(cmdlineFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", cmdlineFile, err)
	}

	buffer := bytes.NewBuffer(cmdlineBytes)

	scanner := bufio.NewScanner(buffer)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i += 1 {
			if data[i] == 0 {
				return i + 1, data[0:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > 0 {
			return len(data), data[:], nil
		}
		// Request more data.
		return 0, nil, nil
	})

	cmdline := make([]string, 0)
	for scanner.Scan() {
		cmdline = append(cmdline, scanner.Text())
	}

	return cmdline, nil
}
