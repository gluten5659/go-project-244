package formatters

import (
	"code/internal/diff"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	nodeAdded     = "added"
	nodeRemoved   = "removed"
	nodeUpdated   = "updated"
	nodeUnchanged = "unchanged"
	nodeNested    = "nested"
)

var errUnknownChangeKind = errors.New("unknown change kind")

type jsonDiff struct {
	Diff []jsonNode `json:"diff"`
}

type jsonNode struct {
	Children *[]jsonNode `json:"children,omitempty"`
	Key      string      `json:"key"`
	NewValue *any        `json:"newValue,omitempty"`
	OldValue *any        `json:"oldValue,omitempty"`
	Type     string      `json:"type"`
	Value    *any        `json:"value,omitempty"`
}

func newNestedNode(key string, children []jsonNode) jsonNode {
	return jsonNode{Key: key, Type: nodeNested, Children: &children}
}

func newUpdatedNode(key string, oldValue, newValue any) jsonNode {
	return jsonNode{Key: key, Type: nodeUpdated, OldValue: &oldValue, NewValue: &newValue}
}

func newValueNode(nodeType, key string, value any) jsonNode {
	return jsonNode{Key: key, Type: nodeType, Value: &value}
}

type jsonFormatter struct{}

func NewJSON() Formatter {
	return jsonFormatter{}
}

func (jsonFormatter) Format(nodes []diff.Node) (string, error) {
	built, err := buildJSONNodes(nodes)
	if err != nil {
		return "", err
	}

	encoded, err := json.MarshalIndent(jsonDiff{Diff: built}, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json diff: %w", err)
	}

	return string(encoded), nil
}

func buildJSONNodes(nodes []diff.Node) ([]jsonNode, error) {
	built := make([]jsonNode, 0, len(nodes))

	for _, node := range nodes {
		encoded, err := buildJSONNode(node)
		if err != nil {
			return nil, err
		}

		built = append(built, encoded)
	}

	return built, nil
}

func buildJSONNode(node diff.Node) (jsonNode, error) {
	switch node.Kind {
	case diff.Nested:
		children, err := buildJSONNodes(node.Children)
		if err != nil {
			return jsonNode{}, err
		}

		return newNestedNode(node.Key, children), nil
	case diff.Updated:
		return newUpdatedNode(node.Key, node.OldValue, node.NewValue), nil
	case diff.Added:
		return newValueNode(nodeAdded, node.Key, node.NewValue), nil
	case diff.Deleted:
		return newValueNode(nodeRemoved, node.Key, node.OldValue), nil
	case diff.Unchanged:
		return newValueNode(nodeUnchanged, node.Key, node.OldValue), nil
	default:
		return jsonNode{}, fmt.Errorf("%w: %d", errUnknownChangeKind, node.Kind)
	}
}
