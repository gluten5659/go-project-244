package output

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

func FormatDiff(diffs []compare.Diff) string {
	builder := strings.Builder{}
	fmt.Fprintln(&builder, "{")

	for _, diff := range diffs {
		fmt.Fprintf(&builder, "  %s %s: %v\n", operation(diff.Change), diff.Key, diff.Value)
	}

	fmt.Fprint(&builder, "}")

	return builder.String()
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
