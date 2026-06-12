package cliapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

var errUsage = errors.New("usage error")

const (
	exitGeneric    = 1
	exitUsage      = 64
	exitDataErr    = 65
	exitNoInput    = 66
	exitIOErr      = 74
	exitPermission = 77
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference",
		OnUsageError: func(_ context.Context, _ *cli.Command, usageErr error, _ bool) error {
			return cli.Exit(fmt.Errorf("%w: %s", errUsage, usageErr.Error()), exitUsage)
		},
	}
}
