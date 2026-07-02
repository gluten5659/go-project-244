package cliapp

import (
	"code"
	"code/internal/files"
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
				Usage:   "output format: " + strings.Join(formatters.SupportedNames(), ", "),
				Validator: func(value string) error {
					if slices.Contains(formatters.SupportedNames(), value) {
						return nil
					}

					return cli.Exit(fmt.Errorf("%w: %q (supported: %s)",
						formatters.ErrUnsupportedFormat, value,
						strings.Join(formatters.SupportedNames(), ", ")), exitUsage)
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
	formatted, err := code.GenDiff(firstFilePath, secondFilePath, format)
	if err != nil {
		return cli.Exit(err, exitCodeFor(err))
	}

	_, err = fmt.Fprintln(cmd.Writer, formatted)
	if err != nil {
		return cli.Exit(err, exitIOErr)
	}

	return nil
}

func exitCodeFor(err error) int {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return exitNoInput
	case errors.Is(err, fs.ErrPermission):
		return exitPermission
	case errors.Is(err, parser.ErrParse):
		return exitDataErr
	case errors.Is(err, files.ErrRead):
		return exitIOErr
	default:
		return exitGeneric
	}
}
