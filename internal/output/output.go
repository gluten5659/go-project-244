package output

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

func FormatDiff(diffs []compare.Diff) {
	builder := strings.Builder{}
	fmt.Fprintln(&builder, "{")

	for _, diff := range diffs {
		fmt.Fprintf(&builder, "  %s", diff.String())
	}

	fmt.Fprintln(&builder, "}")
}
