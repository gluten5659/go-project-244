package main

import (
	"code/internal/cliapp"
	"context"
	"fmt"
	"os"
)

func main() {
	command := cliapp.NewCommand()

	err := command.Run(context.Background(), os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
