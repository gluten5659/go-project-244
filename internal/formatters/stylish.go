package formatters

import (
	"code/internal/diff"
	"fmt"
	"maps"
	"slices"
	"strings"
)

const indentUnit = "    "

type stylishFormatter struct{}

func NewStylish() Formatter {
	return stylishFormatter{}
}

func (stylishFormatter) Format(nodes []diff.Node) (string, error) {
	builder := strings.Builder{}
	writeStylish(&builder, nodes, 0)

	return finalize(&builder), nil
}

func writeStylish(builder *strings.Builder, nodes []diff.Node, level int) {
	fmt.Fprintln(builder, "{")

	for _, node := range nodes {
		writeStylishNode(builder, node, level)
	}

	fmt.Fprintf(builder, "%s}\n", strings.Repeat(indentUnit, level))
}

func writeStylishNode(builder *strings.Builder, node diff.Node, level int) {
	switch node.Kind {
	case diff.Nested:
		writeStylishKey(builder, level, " ", node.Key)
		writeStylish(builder, node.Children, level+1)
	case diff.Updated:
		writeStylishKey(builder, level, "-", node.Key)
		writeStylishValue(builder, node.OldValue, level)
		writeStylishKey(builder, level, "+", node.Key)
		writeStylishValue(builder, node.NewValue, level)
	case diff.Added:
		writeStylishKey(builder, level, "+", node.Key)
		writeStylishValue(builder, node.NewValue, level)
	case diff.Deleted:
		writeStylishKey(builder, level, "-", node.Key)
		writeStylishValue(builder, node.OldValue, level)
	case diff.Unchanged:
		writeStylishKey(builder, level, " ", node.Key)
		writeStylishValue(builder, node.OldValue, level)
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

	for _, key := range collectSortedObjectKeys(object) {
		writeStylishKey(builder, level+1, " ", key)
		writeStylishValue(builder, object[key], level+1)
	}

	fmt.Fprintf(builder, "%s}\n", strings.Repeat(indentUnit, level+1))
}

func collectSortedObjectKeys(object map[string]any) []string {
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
