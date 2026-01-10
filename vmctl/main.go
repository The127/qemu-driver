package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
)

func main() {
	var verbose bool

	cmd := &cli.Command{
		Name:  "vmctl",
		Usage: "control virtual machines",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "verbose output",
				Destination: &verbose,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "create",
				Usage:   "create VM",
				Aliases: []string{"c"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "system-id",
						Aliases:     []string{"i"},
						Value:       "",
						Usage:       "Use `UUID` as the system id of the VM",
						DefaultText: "no id",
						OnlyOnce:    true,

						Validator: func(s string) error {
							_, err := uuid.Parse(s)
							if err != nil {
								return fmt.Errorf("%q is not a valid uuid", s)
							}

							return nil
						},
					},
					&cli.Uint32Flag{
						Name:     "cpus",
						Aliases:  []string{"c"},
						Value:    1,
						Usage:    "Number of vCPUs for the VM",
						OnlyOnce: true,
						Validator: func(n uint32) error {
							if n < 1 {
								return errors.New("the VM must have at least one CPU")
							}

							return nil
						},
					},
					&cli.StringFlag{
						Name:      "memory",
						Aliases:   []string{"m"},
						Value:     "1G",
						Usage:     "Amount of RAM for the VM (suffixes K, M and G are 1024-based)",
						OnlyOnce:  true,
						Validator: validateSize,
					},
					&cli.StringFlag{
						Name:      "disk",
						Aliases:   []string{"d"},
						Value:     "10G",
						Usage:     "Root disk size for the VM (suffixes K, M and G are 1000-based)",
						OnlyOnce:  true,
						Validator: validateSize,
					},
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			println("some output")
			if verbose {
				println("verbose output")
			}
			println("more output")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func validateSize(s string) error {
	bareNum := s
	if strings.HasSuffix(s, "K") || strings.HasSuffix(s, "k") || strings.HasSuffix(s, "M") || strings.HasSuffix(s, "m") || strings.HasSuffix(s, "G") || strings.HasSuffix(s, "g") {
		bareNum = s[:len(s)-1]
	}

	if _, err := strconv.Atoi(bareNum); err != nil {
		return fmt.Errorf("%q is not a valid number", bareNum)
	}

	return nil
}
