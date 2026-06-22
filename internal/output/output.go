package output

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

func FormatDiff(diffs []compare.Diff) string {
	builder := strings.Builder{}
	buildDiff(0, diffs, &builder)

	return builder.String()
}

func buildDiff(level int, diffs []compare.Diff, builder *strings.Builder) {
	fmt.Fprintln(builder, "{")

	for _, diff := range diffs {
		value, isSlice := diff.Value.([]compare.Diff)
		if isSlice {
			fmt.Fprint(builder, strings.Repeat("  ", level))
			fmt.Fprintf(builder, "%s %s: ", operation(diff.Change), diff.Key)

			buildDiff(level+1, value, builder)
		} else {
			fmt.Fprint(builder, strings.Repeat("  ", level))
			fmt.Fprintf(builder, "%s %s: %v\n", operation(diff.Change), diff.Key, diff.Value)
		}
	}

	fmt.Fprint(builder, strings.Repeat("  ", level))
	fmt.Fprintln(builder, "}")
}

func operation(change compare.Changes) string {
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
