package formatters

import (
	"code/internal/diff"
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

type Formatter interface {
	Format(nodes []diff.Node) (string, error)
}

func New(name string) (Formatter, error) {
	switch name {
	case Stylish:
		return NewStylish(), nil
	case Plain:
		return NewPlain(), nil
	case JSON:
		return NewJSON(), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedFormat, name)
	}
}

func ListSupportedNames() []string {
	return []string{JSON, Plain, Stylish}
}

func finalize(builder *strings.Builder) string {
	return strings.TrimRight(builder.String(), "\n")
}
