package formatters

import (
	"code/internal/compare"
	"fmt"
	"strings"
)

func writePlain(builder *strings.Builder, diffs []compare.Diff, parentPath string) {
	for index := 0; index < len(diffs); index++ {
		diff := diffs[index]
		path := plainPath(parentPath, diff.Key)

		switch diff.Change {
		case compare.NoChanges:
			if children, isNested := diff.Value.([]compare.Diff); isNested {
				writePlain(builder, children, path)
			}
		case compare.Deleted:
			if index+1 < len(diffs) && isUpdatedTo(diffs[index+1], diff.Key) {
				fmt.Fprintf(builder, "Property '%s' was updated. From %s to %s\n",
					path, plainValue(diff.Value), plainValue(diffs[index+1].Value))
				index++
			} else {
				fmt.Fprintf(builder, "Property '%s' was removed\n", path)
			}
		case compare.Added:
			fmt.Fprintf(builder, "Property '%s' was added with value: %s\n", path, plainValue(diff.Value))
		}
	}
}

func plainPath(parentPath, key string) string {
	if parentPath == "" {
		return key
	}

	return parentPath + "." + key
}

func isUpdatedTo(next compare.Diff, key string) bool {
	return next.Change == compare.Added && next.Key == key
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
