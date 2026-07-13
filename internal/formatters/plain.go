package formatters

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

func writePlainRoot(builder *strings.Builder, nodes []compare.Node) error {
	writePlain(builder, nodes, "")

	return nil
}

func writePlain(builder *strings.Builder, nodes []compare.Node, parentPath string) {
	for _, node := range nodes {
		path := plainPath(parentPath, node.Key)

		switch node.Kind {
		case compare.Nested:
			writePlain(builder, node.Children, path)
		case compare.Updated:
			fmt.Fprintf(builder, "Property '%s' was updated. From %s to %s\n",
				path, plainValue(node.OldValue), plainValue(node.NewValue))
		case compare.Deleted:
			fmt.Fprintf(builder, "Property '%s' was removed\n", path)
		case compare.Added:
			fmt.Fprintf(builder, "Property '%s' was added with value: %s\n",
				path, plainValue(node.Value))
		case compare.Unchanged:
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
