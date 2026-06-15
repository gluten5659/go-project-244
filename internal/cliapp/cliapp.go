package cliapp

import (
	"code/internal/core"
	"code/internal/files"
	"context"
	"errors"
	"fmt"
	"io/fs"

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
	var firstFilePath, secondFilePath string

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
				Value:   "stylish",
				Usage:   "output format",
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
			return run(cmd, firstFilePath, secondFilePath)
		},
	}
}

func run(cmd *cli.Command, firstFilePath, secondFilePath string) error {
	firstConfig, err := load(firstFilePath)
	if err != nil {
		return cli.Exit(err, exitCodeFor(err))
	}

	secondConfig, err := load(secondFilePath)
	if err != nil {
		return cli.Exit(err, exitCodeFor(err))
	}

	_, err = fmt.Fprintf(cmd.Writer, "%v\n%v\n", firstConfig, secondConfig)
	if err != nil {
		return cli.Exit(err, exitIOErr)
	}

	return nil
}

func load(path string) (map[string]any, error) {
	content, err := files.Read(path)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	config, err := core.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	return config, nil
}

func exitCodeFor(err error) int {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return exitNoInput
	case errors.Is(err, fs.ErrPermission):
		return exitPermission
	case errors.Is(err, core.ErrParse):
		return exitDataErr
	case errors.Is(err, files.ErrRead):
		return exitIOErr
	default:
		return exitGeneric
	}
}
