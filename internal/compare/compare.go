package compare

import (
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

	keys := append(slices.Collect(maps.Keys(firstFile)), slices.Collect(maps.Keys(secondFile))...)
	slices.Sort(keys)
	keys = slices.Compact(keys)

	for _, key := range keys {
		firstFileValue, isKeyInFirstFile := firstFile[key]

		if !isKeyInFirstFile {
			diff = append(diff, Diff{Change: Added, Key: key, Value: secondFile[key]})

			continue
		}

		secondFileValue, isKeyInSecondFile := secondFile[key]

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
