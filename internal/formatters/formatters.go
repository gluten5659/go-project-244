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

type writer func(*strings.Builder, []compare.Diff) error

func writers() map[string]writer {
	return map[string]writer{
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

	err := write(&builder, diffs)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(builder.String(), "\n"), nil
}

func SupportedNames() []string {
	names := slices.Collect(maps.Keys(writers()))
	slices.Sort(names)

	return names
}

func writeStylishRoot(builder *strings.Builder, diffs []compare.Diff) error {
	writeStylish(builder, diffs, 0)

	return nil
}

func writePlainRoot(builder *strings.Builder, diffs []compare.Diff) error {
	writePlain(builder, diffs, "")

	return nil
}

type mergedDiff struct {
	compare.Diff

	newValue any
	updated  bool
}

func mergeUpdates(diffs []compare.Diff) []mergedDiff {
	merged := make([]mergedDiff, 0, len(diffs))

	for index := 0; index < len(diffs); index++ {
		diff := diffs[index]

		if diff.Kind == compare.Deleted && index+1 < len(diffs) &&
			isUpdatedTo(diffs[index+1], diff.Key) {
			merged = append(merged, mergedDiff{Diff: diff, newValue: diffs[index+1].Value, updated: true})
			index++

			continue
		}

		merged = append(merged, mergedDiff{Diff: diff})
	}

	return merged
}

func isUpdatedTo(next compare.Diff, key string) bool {
	return next.Kind == compare.Added && next.Key == key
}
