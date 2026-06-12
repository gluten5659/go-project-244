package main

import (
	"code/internal/cliapp"
	"context"
	"fmt"
	"os"
)

func main() {
	err := cliapp.NewCommand().Run(context.Background(), os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
