package formatters

import (
	"code/internal/diff"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

const (
	Stylish = "stylish"
	Plain   = "plain"
	JSON    = "json"
)

var ErrUnsupportedFormat = errors.New("unsupported output format")

type Formatter interface {
	Format(nodes []diff.Node) (string, error)
}

func registry() map[string]func() Formatter {
	return map[string]func() Formatter{
		Stylish: func() Formatter { return stylishFormatter{} },
		Plain:   func() Formatter { return plainFormatter{} },
		JSON:    func() Formatter { return jsonFormatter{} },
	}
}

func New(name string) (Formatter, error) {
	build, ok := registry()[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedFormat, name)
	}

	return build(), nil
}

func SupportedNames() []string {
	names := slices.Collect(maps.Keys(registry()))
	slices.Sort(names)

	return names
}

func finalize(builder *strings.Builder) string {
	return strings.TrimRight(builder.String(), "\n")
}
