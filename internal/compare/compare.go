package compare

import (
	"cmp"
	"maps"
	"reflect"
	"slices"
)

type Change int

const (
	Unchanged Change = iota
	Added
	Deleted
)

type Diff struct {
	Kind  Change
	Key   string
	Value any
}

func Compare(firstFile, secondFile map[string]any) []Diff {
	keys := sortedKeys(firstFile, secondFile)

	diff := make([]Diff, 0, len(keys))
	for _, key := range keys {
		diff = append(diff, compareKey(key, firstFile, secondFile)...)
	}

	return diff
}

func compareKey(key string, firstFile, secondFile map[string]any) []Diff {
	firstFileValue, isKeyInFirstFile := firstFile[key]
	secondFileValue, isKeyInSecondFile := secondFile[key]

	switch {
	case !isKeyInFirstFile:
		return []Diff{added(key, secondFileValue)}
	case !isKeyInSecondFile:
		return []Diff{deleted(key, firstFileValue)}
	}

	firstFileValueMap, isFirstValueMap := firstFileValue.(map[string]any)
	secondFileValueMap, isSecondValueMap := secondFileValue.(map[string]any)

	if isFirstValueMap && isSecondValueMap {
		return []Diff{unchanged(key, Compare(firstFileValueMap, secondFileValueMap))}
	}

	if reflect.DeepEqual(firstFileValue, secondFileValue) {
		return []Diff{unchanged(key, firstFileValue)}
	}

	return []Diff{deleted(key, firstFileValue), added(key, secondFileValue)}
}

func added(key string, value any) Diff {
	return Diff{Kind: Added, Key: key, Value: expand(value)}
}

func deleted(key string, value any) Diff {
	return Diff{Kind: Deleted, Key: key, Value: expand(value)}
}

func unchanged(key string, value any) Diff {
	return Diff{Kind: Unchanged, Key: key, Value: value}
}

func expand(value any) any {
	valueMap, isValueMap := value.(map[string]any)
	if !isValueMap {
		return value
	}

	return Compare(valueMap, valueMap)
}

func sortedKeys[Key cmp.Ordered, Value any](sources ...map[Key]Value) []Key {
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
