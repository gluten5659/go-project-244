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
)

var ErrUnsupportedFormat = errors.New("unsupported output format")

func Format(diffs []compare.Diff, name string) (string, error) {
	builder := strings.Builder{}

	switch name {
	case Stylish:
		writeStylish(&builder, diffs, 0)
	case Plain:
		writePlain(&builder, diffs, "")
	default:
		return "", fmt.Errorf("%w: %q", ErrUnsupportedFormat, name)
	}

	return strings.TrimRight(builder.String(), "\n"), nil
}
