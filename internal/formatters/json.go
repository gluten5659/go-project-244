package formatters

import (
	"code/internal/compare"
	"encoding/json"
	"fmt"
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

func writeJSON(builder *strings.Builder, nodes []compare.Node) error {
	document := map[string]any{fieldDiff: jsonNodes(nodes)}

	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json diff: %w", err)
	}

	builder.Write(encoded)

	return nil
}

func jsonNodes(nodes []compare.Node) []map[string]any {
	encoded := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		encoded = append(encoded, jsonNode(node))
	}

	return encoded
}

func jsonNode(node compare.Node) map[string]any {
	switch node.Kind {
	case compare.Nested:
		return map[string]any{
			fieldKey:      node.Key,
			fieldType:     nodeNested,
			fieldChildren: jsonNodes(node.Children),
		}
	case compare.Updated:
		return map[string]any{
			fieldKey:      node.Key,
			fieldType:     nodeUpdated,
			fieldOldValue: node.OldValue,
			fieldNewValue: node.NewValue,
		}
	case compare.Added:
		return map[string]any{
			fieldKey:   node.Key,
			fieldType:  nodeAdded,
			fieldValue: node.Value,
		}
	case compare.Deleted:
		return map[string]any{
			fieldKey:   node.Key,
			fieldType:  nodeRemoved,
			fieldValue: node.Value,
		}
	case compare.Unchanged:
		return map[string]any{
			fieldKey:   node.Key,
			fieldType:  nodeUnchanged,
			fieldValue: node.Value,
		}
	}

	return nil
}
