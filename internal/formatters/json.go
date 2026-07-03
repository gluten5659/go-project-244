package formatters

import (
	"code/internal/compare"
	"encoding/json"
	"strings"
)

const (
	fieldKey      = "key"
	fieldType     = "type"
	fieldValue    = "value"
	fieldOldValue = "oldValue"
	fieldNewValue = "newValue"
	fieldChildren = "children"
	fieldDiff     = "diff"

	nodeAdded     = "added"
	nodeRemoved   = "removed"
	nodeUpdated   = "updated"
	nodeUnchanged = "unchanged"
	nodeNested    = "nested"
)

func writeJSON(builder *strings.Builder, diffs []compare.Diff) {
	document := map[string]any{fieldDiff: jsonNodes(diffs)}

	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		panic(err)
	}

	builder.Write(encoded)
}

func jsonNodes(diffs []compare.Diff) []map[string]any {
	merged := mergeUpdates(diffs)

	nodes := make([]map[string]any, 0, len(merged))
	for _, entry := range merged {
		if entry.updated {
			nodes = append(nodes, updatedNode(entry.Key, entry.Value, entry.newValue))

			continue
		}

		nodes = append(nodes, changeNode(entry.Diff))
	}

	return nodes
}

func changeNode(diff compare.Diff) map[string]any {
	children, isTree := diff.Value.([]compare.Diff)

	switch {
	case diff.Kind == compare.Unchanged && isTree:
		return map[string]any{
			fieldKey:      diff.Key,
			fieldType:     nodeNested,
			fieldChildren: jsonNodes(children),
		}
	case diff.Kind == compare.Unchanged:
		return map[string]any{
			fieldKey:   diff.Key,
			fieldType:  nodeUnchanged,
			fieldValue: diff.Value,
		}
	case diff.Kind == compare.Added:
		return map[string]any{
			fieldKey:   diff.Key,
			fieldType:  nodeAdded,
			fieldValue: collapse(diff.Value),
		}
	default:
		return map[string]any{
			fieldKey:   diff.Key,
			fieldType:  nodeRemoved,
			fieldValue: collapse(diff.Value),
		}
	}
}

func updatedNode(key string, oldValue, newValue any) map[string]any {
	return map[string]any{
		fieldKey:      key,
		fieldType:     nodeUpdated,
		fieldOldValue: collapse(oldValue),
		fieldNewValue: collapse(newValue),
	}
}

func collapse(value any) any {
	children, isTree := value.([]compare.Diff)
	if !isTree {
		return value
	}

	object := make(map[string]any, len(children))
	for _, child := range children {
		object[child.Key] = collapse(child.Value)
	}

	return object
}
