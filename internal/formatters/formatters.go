package formatters

import (
	"code/internal/compare"
	"errors"
	"fmt"
	"strings"
)

const (
	Stylish = "stylish"
	Plain   = "plain"
	JSON    = "json"
)

var ErrUnsupportedFormat = errors.New("unsupported output format")

type writer func(*strings.Builder, []compare.Node) error

type namedFormatter struct {
	name  string
	write writer
}

func supportedFormatters() []namedFormatter {
	return []namedFormatter{
		{name: JSON, write: writeJSON},
		{name: Plain, write: writePlainRoot},
		{name: Stylish, write: writeStylishRoot},
	}
}

func Format(nodes []compare.Node, name string) (string, error) {
	write := writerFor(name)
	if write == nil {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedFormat, name)
	}

	builder := strings.Builder{}

	err := write(&builder, nodes)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(builder.String(), "\n"), nil
}

func writerFor(name string) writer {
	for _, formatter := range supportedFormatters() {
		if formatter.name == name {
			return formatter.write
		}
	}

	return nil
}

func SupportedNames() []string {
	available := supportedFormatters()

	names := make([]string, 0, len(available))
	for _, formatter := range available {
		names = append(names, formatter.name)
	}

	return names
}
