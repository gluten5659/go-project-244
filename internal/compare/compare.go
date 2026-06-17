package compare

import (
	"maps"
	"slices"
)

type Changes int

const (
	NoChanges = iota
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
		firstFileValue, isKeyInMap := firstFile[key]

		if !isKeyInMap {
			diff = append(diff, Diff{Change: Added, Key: key, Value: firstFileValue})

			continue
		}

		secondFileValue, isKeyInMap := secondFile[key]

		if !isKeyInMap {
			diff = append(diff, Diff{Change: Deleted, Key: key, Value: secondFileValue})

			continue
		}

		if firstFileValue != secondFileValue {
			diff = append(diff, Diff{Change: Deleted, Key: key, Value: firstFileValue})
			diff = append(diff, Diff{Change: Added, Key: key, Value: secondFileValue})

			continue
		}

		diff = append(diff, Diff{Change: NoChanges, Key: key, Value: firstFileValue})
	}

	return diff
}
