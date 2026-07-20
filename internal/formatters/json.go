package formatters

import (
	"code/internal/diff"
	"encoding/json"
	"fmt"
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

type jsonFormatter struct{}

func (jsonFormatter) Format(nodes []diff.Node) (string, error) {
	document := map[string]any{fieldDiff: jsonNodes(nodes)}

	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json diff: %w", err)
	}

	return string(encoded), nil
}

func jsonNodes(nodes []diff.Node) []map[string]any {
	encoded := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		encoded = append(encoded, jsonNode(node))
	}

	return encoded
}

func jsonNode(node diff.Node) map[string]any {
	switch node.Kind {
	case diff.Nested:
		return map[string]any{
			fieldKey:      node.Key,
			fieldType:     nodeNested,
			fieldChildren: jsonNodes(node.Children),
		}
	case diff.Updated:
		return map[string]any{
			fieldKey:      node.Key,
			fieldType:     nodeUpdated,
			fieldOldValue: node.OldValue,
			fieldNewValue: node.NewValue,
		}
	case diff.Added:
		return map[string]any{
			fieldKey:   node.Key,
			fieldType:  nodeAdded,
			fieldValue: node.Value,
		}
	case diff.Deleted:
		return map[string]any{
			fieldKey:   node.Key,
			fieldType:  nodeRemoved,
			fieldValue: node.Value,
		}
	case diff.Unchanged:
		return map[string]any{
			fieldKey:   node.Key,
			fieldType:  nodeUnchanged,
			fieldValue: node.Value,
		}
	}

	return nil
}
