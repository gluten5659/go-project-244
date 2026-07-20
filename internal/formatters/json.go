package formatters

import (
	"code/internal/diff"
	"encoding/json"
	"fmt"
)

const (
	nodeAdded     = "added"
	nodeRemoved   = "removed"
	nodeUpdated   = "updated"
	nodeUnchanged = "unchanged"
	nodeNested    = "nested"
)

type jsonDiff struct {
	Diff []any `json:"diff"`
}

type jsonValueNode struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type jsonUpdatedNode struct {
	Key      string `json:"key"`
	NewValue any    `json:"newValue"`
	OldValue any    `json:"oldValue"`
	Type     string `json:"type"`
}

type jsonNestedNode struct {
	Children []any  `json:"children"`
	Key      string `json:"key"`
	Type     string `json:"type"`
}

type jsonFormatter struct{}

func (jsonFormatter) Format(nodes []diff.Node) (string, error) {
	encoded, err := json.MarshalIndent(jsonDiff{Diff: jsonNodes(nodes)}, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json diff: %w", err)
	}

	return string(encoded), nil
}

func jsonNodes(nodes []diff.Node) []any {
	encoded := make([]any, 0, len(nodes))
	for _, node := range nodes {
		encoded = append(encoded, jsonNode(node))
	}

	return encoded
}

func jsonNode(node diff.Node) any {
	switch node.Kind {
	case diff.Nested:
		return jsonNestedNode{Children: jsonNodes(node.Children), Key: node.Key, Type: nodeNested}
	case diff.Updated:
		return jsonUpdatedNode{
			Key:      node.Key,
			NewValue: node.NewValue,
			OldValue: node.OldValue,
			Type:     nodeUpdated,
		}
	case diff.Added:
		return jsonValueNode{Key: node.Key, Type: nodeAdded, Value: node.Value}
	case diff.Deleted:
		return jsonValueNode{Key: node.Key, Type: nodeRemoved, Value: node.Value}
	case diff.Unchanged:
		return jsonValueNode{Key: node.Key, Type: nodeUnchanged, Value: node.Value}
	}

	return nil
}
