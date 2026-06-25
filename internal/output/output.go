package output

import (
	"code/internal/compare"
	"errors"
	"fmt"
	"strings"
)

const (
	indentUnit    = "    "
	stylishFormat = "stylish"
	plainFormat   = "plain"
)

var ErrUnsupportedFormat = errors.New("unsupported output format")

func FormatDiff(diffs []compare.Diff, format string) (string, error) {
	builder := strings.Builder{}

	err := writeBlockByFormat(format, &builder, diffs)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(builder.String(), "\n"), nil
}

func writeBlockByFormat(format string, builder *strings.Builder, diffs []compare.Diff) error {
	switch format {
	case stylishFormat:
		writeBlockStylish(builder, diffs, 0)
	case plainFormat:
		writeBlockPlain(builder, diffs, "")
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedFormat, format)
	}

	return nil
}

func writeBlockStylish(builder *strings.Builder, diffs []compare.Diff, level int) {
	fmt.Fprintln(builder, "{")

	entryIndent := strings.Repeat(indentUnit, level) + "  "
	for _, diff := range diffs {
		fmt.Fprintf(builder, "%s%s %s: ", entryIndent, operationStylish(diff.Change), diff.Key)

		if children, isNested := diff.Value.([]compare.Diff); isNested {
			writeBlockStylish(builder, children, level+1)
		} else {
			fmt.Fprintln(builder, formatValue(diff.Value))
		}
	}

	fmt.Fprintf(builder, "%s}\n", strings.Repeat(indentUnit, level))
}

func operationStylish(change compare.Changes) string {
	switch change {
	case compare.Added:
		return "+"
	case compare.Deleted:
		return "-"
	case compare.NoChanges:
		return " "
	}

	return " "
}

func writeBlockPlain(builder *strings.Builder, diffs []compare.Diff, parentPath string) {
	for index := 0; index < len(diffs); index++ {
		diff := diffs[index]
		path := plainPath(parentPath, diff.Key)

		switch diff.Change {
		case compare.NoChanges:
			if children, isNested := diff.Value.([]compare.Diff); isNested {
				writeBlockPlain(builder, children, path)
			}
		case compare.Deleted:
			if index+1 < len(diffs) && isUpdatedTo(diffs[index+1], diff.Key) {
				fmt.Fprintf(builder, "Property '%s' was updated. From %s to %s\n",
					path, plainValue(diff.Value), plainValue(diffs[index+1].Value))
				index++
			} else {
				fmt.Fprintf(builder, "Property '%s' was removed\n", path)
			}
		case compare.Added:
			fmt.Fprintf(builder, "Property '%s' was added with value: %s\n", path, plainValue(diff.Value))
		}
	}
}

func plainPath(parentPath, key string) string {
	if parentPath == "" {
		return key
	}

	return parentPath + "." + key
}

func isUpdatedTo(next compare.Diff, key string) bool {
	return next.Change == compare.Added && next.Key == key
}

func plainValue(value any) string {
	switch typed := value.(type) {
	case []compare.Diff:
		return "[complex value]"
	case nil:
		return "null"
	case string:
		return "'" + typed + "'"
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func formatValue(value any) string {
	if value == nil {
		return "null"
	}

	return fmt.Sprintf("%v", value)
}
