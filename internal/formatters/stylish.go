package formatters

import (
	"code/internal/compare"
	"fmt"
	"maps"
	"slices"
	"strings"
)

const indentUnit = "    "

func writeStylishRoot(builder *strings.Builder, nodes []compare.Node) error {
	writeStylish(builder, nodes, 0)

	return nil
}

func writeStylish(builder *strings.Builder, nodes []compare.Node, level int) {
	fmt.Fprintln(builder, "{")

	for _, node := range nodes {
		writeStylishNode(builder, node, level)
	}

	fmt.Fprintf(builder, "%s}\n", strings.Repeat(indentUnit, level))
}

func writeStylishNode(builder *strings.Builder, node compare.Node, level int) {
	switch node.Kind {
	case compare.Nested:
		writeStylishKey(builder, level, " ", node.Key)
		writeStylish(builder, node.Children, level+1)
	case compare.Updated:
		writeStylishKey(builder, level, "-", node.Key)
		writeStylishValue(builder, node.OldValue, level)
		writeStylishKey(builder, level, "+", node.Key)
		writeStylishValue(builder, node.NewValue, level)
	case compare.Added:
		writeStylishKey(builder, level, "+", node.Key)
		writeStylishValue(builder, node.Value, level)
	case compare.Deleted:
		writeStylishKey(builder, level, "-", node.Key)
		writeStylishValue(builder, node.Value, level)
	case compare.Unchanged:
		writeStylishKey(builder, level, " ", node.Key)
		writeStylishValue(builder, node.Value, level)
	}
}

func writeStylishKey(builder *strings.Builder, level int, marker, key string) {
	entryIndent := strings.Repeat(indentUnit, level) + "  "
	fmt.Fprintf(builder, "%s%s %s: ", entryIndent, marker, key)
}

func writeStylishValue(builder *strings.Builder, value any, level int) {
	object, isObject := value.(map[string]any)
	if !isObject {
		fmt.Fprintln(builder, formatValue(value))

		return
	}

	fmt.Fprintln(builder, "{")

	for _, key := range sortedObjectKeys(object) {
		writeStylishKey(builder, level+1, " ", key)
		writeStylishValue(builder, object[key], level+1)
	}

	fmt.Fprintf(builder, "%s}\n", strings.Repeat(indentUnit, level+1))
}

func sortedObjectKeys(object map[string]any) []string {
	keys := slices.Collect(maps.Keys(object))
	slices.Sort(keys)

	return keys
}

func formatValue(value any) string {
	if value == nil {
		return "null"
	}

	return fmt.Sprintf("%v", value)
}
