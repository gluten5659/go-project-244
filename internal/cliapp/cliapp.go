package cliapp

import (
	"code"
	"code/internal/formatters"
	"code/internal/parser"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

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
	var firstFilePath, secondFilePath, format string

	supportedFormats := formatters.ListSupportedNames()
	listedFormats := strings.Join(supportedFormats, ", ")

	return &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference",
		OnUsageError: func(_ context.Context, _ *cli.Command, usageErr error, _ bool) error {
			return cli.Exit(fmt.Errorf("%w: %s", errUsage, usageErr.Error()), exitUsage)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   formatters.Stylish,
				Usage:   "output format: " + listedFormats,
				Validator: func(value string) error {
					if slices.Contains(supportedFormats, value) {
						return nil
					}

					return fmt.Errorf("%w (supported: %s)",
						formatters.ErrUnsupportedFormat, listedFormats)
				},
				Destination: &format,
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "First file path",
				Destination: &firstFilePath,
			},
			&cli.StringArg{
				Name:        "Second file path",
				Destination: &secondFilePath,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			return run(cmd, firstFilePath, secondFilePath, format)
		},
	}
}

func run(cmd *cli.Command, firstFilePath, secondFilePath, format string) error {
	if firstFilePath == "" || secondFilePath == "" || cmd.NArg() > 0 {
		return cli.Exit(fmt.Errorf("%w: exactly two file paths are required", errUsage), exitUsage)
	}

	formatted, err := code.GenDiff(firstFilePath, secondFilePath, format)
	if err != nil {
		return cli.Exit(err, resolveExitCode(err))
	}

	_, err = fmt.Fprintln(cmd.Writer, formatted)
	if err != nil {
		return cli.Exit(err, exitIOErr)
	}

	return nil
}

func resolveExitCode(err error) int {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return exitNoInput
	case errors.Is(err, fs.ErrPermission):
		return exitPermission
	case errors.Is(err, parser.ErrParse):
		return exitDataErr
	case errors.Is(err, parser.ErrRead):
		return exitIOErr
	default:
		return exitGeneric
	}
}
