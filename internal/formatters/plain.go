package formatters

import (
	"code/internal/diff"
	"fmt"
	"strings"
)

type plainFormatter struct{}

func (plainFormatter) Format(nodes []diff.Node) (string, error) {
	builder := strings.Builder{}
	writePlain(&builder, nodes, "")

	return finalize(&builder), nil
}

func writePlain(builder *strings.Builder, nodes []diff.Node, parentPath string) {
	for _, node := range nodes {
		path := plainPath(parentPath, node.Key)

		switch node.Kind {
		case diff.Nested:
			writePlain(builder, node.Children, path)
		case diff.Updated:
			fmt.Fprintf(builder, "Property '%s' was updated. From %s to %s\n",
				path, plainValue(node.OldValue), plainValue(node.NewValue))
		case diff.Deleted:
			fmt.Fprintf(builder, "Property '%s' was removed\n", path)
		case diff.Added:
			fmt.Fprintf(builder, "Property '%s' was added with value: %s\n",
				path, plainValue(node.Value))
		case diff.Unchanged:
		}
	}
}

func plainPath(parentPath, key string) string {
	if parentPath == "" {
		return key
	}

	return parentPath + "." + key
}

func plainValue(value any) string {
	switch typed := value.(type) {
	case map[string]any:
		return "[complex value]"
	case nil:
		return "null"
	case string:
		return "'" + typed + "'"
	default:
		return fmt.Sprintf("%v", typed)
	}
}
