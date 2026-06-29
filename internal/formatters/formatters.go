package formatters

import (
	"code/internal/compare"
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

func writers() map[string]func(*strings.Builder, []compare.Diff) {
	return map[string]func(*strings.Builder, []compare.Diff){
		Stylish: writeStylishRoot,
		Plain:   writePlainRoot,
		JSON:    writeJSON,
	}
}

func Format(diffs []compare.Diff, name string) (string, error) {
	write, isSupported := writers()[name]
	if !isSupported {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedFormat, name)
	}

	builder := strings.Builder{}
	write(&builder, diffs)

	return strings.TrimRight(builder.String(), "\n"), nil
}

func SupportedNames() []string {
	names := slices.Collect(maps.Keys(writers()))
	slices.Sort(names)

	return names
}

func writeStylishRoot(builder *strings.Builder, diffs []compare.Diff) {
	writeStylish(builder, diffs, 0)
}

func writePlainRoot(builder *strings.Builder, diffs []compare.Diff) {
	writePlain(builder, diffs, "")
}
