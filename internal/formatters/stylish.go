package formatters

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

const indentUnit = "    "

func writeStylish(builder *strings.Builder, diffs []compare.Diff, level int) {
	fmt.Fprintln(builder, "{")

	entryIndent := strings.Repeat(indentUnit, level) + "  "
	for _, diff := range diffs {
		fmt.Fprintf(builder, "%s%s %s: ", entryIndent, operationStylish(diff.Kind), diff.Key)

		if children, isNested := diff.Value.([]compare.Diff); isNested {
			writeStylish(builder, children, level+1)
		} else {
			fmt.Fprintln(builder, formatValue(diff.Value))
		}
	}

	fmt.Fprintf(builder, "%s}\n", strings.Repeat(indentUnit, level))
}

func operationStylish(kind compare.Change) string {
	switch kind {
	case compare.Added:
		return "+"
	case compare.Deleted:
		return "-"
	case compare.Unchanged:
		return " "
	}

	return " "
}

func formatValue(value any) string {
	if value == nil {
		return "null"
	}

	return fmt.Sprintf("%v", value)
}
