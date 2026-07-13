package compare

import (
	"cmp"
	"maps"
	"reflect"
	"slices"
)

type Kind int

const (
	Unchanged Kind = iota
	Added
	Deleted
	Updated
	Nested
)

type Node struct {
	Key      string
	Kind     Kind
	Value    any
	OldValue any
	NewValue any
	Children []Node
}

func Compare(firstFile, secondFile map[string]any) []Node {
	keys := collectSortedKeys(firstFile, secondFile)

	nodes := make([]Node, 0, len(keys))
	for _, key := range keys {
		nodes = append(nodes, compareKey(key, firstFile, secondFile))
	}

	return nodes
}

func compareKey(key string, firstFile, secondFile map[string]any) Node {
	firstValue, isKeyInFirstFile := firstFile[key]
	secondValue, isKeyInSecondFile := secondFile[key]

	switch {
	case !isKeyInFirstFile:
		return Node{Key: key, Kind: Added, Value: secondValue}
	case !isKeyInSecondFile:
		return Node{Key: key, Kind: Deleted, Value: firstValue}
	}

	firstObject, isFirstObject := firstValue.(map[string]any)
	secondObject, isSecondObject := secondValue.(map[string]any)

	if isFirstObject && isSecondObject {
		return Node{Key: key, Kind: Nested, Children: Compare(firstObject, secondObject)}
	}

	if reflect.DeepEqual(firstValue, secondValue) {
		return Node{Key: key, Kind: Unchanged, Value: firstValue}
	}

	return Node{Key: key, Kind: Updated, OldValue: firstValue, NewValue: secondValue}
}

func collectSortedKeys[Key cmp.Ordered, Value any](sources ...map[Key]Value) []Key {
	capacity := 0
	for _, source := range sources {
		capacity += len(source)
	}

	keys := make([]Key, 0, capacity)
	for _, source := range sources {
		keys = append(keys, slices.Collect(maps.Keys(source))...)
	}

	slices.Sort(keys)

	return slices.Compact(keys)
}
