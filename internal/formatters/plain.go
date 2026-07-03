package formatters

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

func writePlain(builder *strings.Builder, diffs []compare.Diff, parentPath string) {
	for _, entry := range mergeUpdates(diffs) {
		path := plainPath(parentPath, entry.Key)

		switch {
		case entry.updated:
			fmt.Fprintf(builder, "Property '%s' was updated. From %s to %s\n",
				path, plainValue(entry.Value), plainValue(entry.newValue))
		case entry.Kind == compare.Unchanged:
			if children, isNested := entry.Value.([]compare.Diff); isNested {
				writePlain(builder, children, path)
			}
		case entry.Kind == compare.Deleted:
			fmt.Fprintf(builder, "Property '%s' was removed\n", path)
		case entry.Kind == compare.Added:
			fmt.Fprintf(
				builder,
				"Property '%s' was added with value: %s\n",
				path,
				plainValue(entry.Value),
			)
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
