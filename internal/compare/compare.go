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
	var diff []Diff

	for _, key := range sortedKeys(firstFile, secondFile) {
		firstFileValue, isKeyInFirstFile := firstFile[key]
		secondFileValue, isKeyInSecondFile := secondFile[key]
		firstFileValueMap, isFirstValueMap := firstFile[key].(map[string]any)
		secondFileValueMap, isSecondValueMap := secondFile[key].(map[string]any)

		if isFirstValueMap && isSecondValueMap {
			diff = append(
				diff,
				Diff{
					Change: NoChanges,
					Key:    key,
					Value:  Compare(firstFileValueMap, secondFileValueMap),
				},
			)

			continue
		}

		if isFirstValueMap {
			diff = append(
				diff,
				Diff{
					Change: Deleted,
					Key:    key,
					Value:  Compare(firstFileValueMap, firstFileValueMap),
				},
			)

			continue
		}

		if isSecondValueMap {
			diff = append(
				diff,
				Diff{
					Change: Added,
					Key:    key,
					Value:  Compare(secondFileValueMap, secondFileValueMap),
				},
			)

			continue
		}

		if !isKeyInFirstFile {
			diff = append(diff, Diff{Change: Added, Key: key, Value: secondFile[key]})

			continue
		}

		if !isKeyInSecondFile {
			diff = append(diff, Diff{Change: Deleted, Key: key, Value: firstFileValue})

			continue
		}

		if !reflect.DeepEqual(firstFileValue, secondFileValue) {
			diff = append(diff, Diff{Change: Deleted, Key: key, Value: firstFileValue})
			diff = append(diff, Diff{Change: Added, Key: key, Value: secondFileValue})

			continue
		}

		diff = append(diff, Diff{Change: NoChanges, Key: key, Value: firstFileValue})
	}

	return diff
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
