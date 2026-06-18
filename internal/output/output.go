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
		fmt.Fprintf(&builder, "  %s\n", diff.String())
	}

	fmt.Fprintln(&builder, "}")

	return builder.String()
}
