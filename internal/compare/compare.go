package compare

import (
	"cmp"
	"maps"
	"reflect"
	"slices"
)

type Changes int

const (
	NoChanges Changes = iota
	Added
	Deleted
)

type Diff struct {
	Change Changes
	Key    string
	Value  any
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
	firstFileValueMap, isFirstValueMap := firstFile[key].(map[string]any)
	secondFileValueMap, isSecondValueMap := secondFile[key].(map[string]any)

	if isFirstValueMap && isSecondValueMap {
		return []Diff{{
			Change: NoChanges,
			Key:    key,
			Value:  Compare(firstFileValueMap, secondFileValueMap),
		}}
	}

	if isFirstValueMap {
		return []Diff{{
			Change: Deleted,
			Key:    key,
			Value:  Compare(firstFileValueMap, firstFileValueMap),
		}}
	}

	if isSecondValueMap {
		return []Diff{{
			Change: Added,
			Key:    key,
			Value:  Compare(secondFileValueMap, secondFileValueMap),
		}}
	}

	if !isKeyInFirstFile {
		return []Diff{{Change: Added, Key: key, Value: secondFile[key]}}
	}

	if !isKeyInSecondFile {
		return []Diff{{Change: Deleted, Key: key, Value: firstFileValue}}
	}

	if !reflect.DeepEqual(firstFileValue, secondFileValue) {
		return []Diff{
			{Change: Deleted, Key: key, Value: firstFileValue},
			{Change: Added, Key: key, Value: secondFileValue},
		}
	}

	return []Diff{{Change: NoChanges, Key: key, Value: firstFileValue}}
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
