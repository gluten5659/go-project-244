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

		if !isKeyInFirstFile {
			diff = append(diff, Diff{Change: Added, Key: key, Value: secondFile[key]})

			continue
		}

		secondFileValue, isKeyInSecondFile := secondFile[key]

		if !isKeyInSecondFile {
			diff = append(diff, Diff{Change: Deleted, Key: key, Value: firstFileValue})

			continue
		}

		typeOfFirstFileValue := reflect.ValueOf(firstFileValue).Kind()
		typeOfSecondFileValue := reflect.ValueOf(secondFileValue).Kind()

		if typeOfFirstFileValue == reflect.Map && typeOfSecondFileValue == reflect.Map {
			diff = append(
				diff,
				Diff{
					Change: NoChanges,
					Key:    key,
					Value:  Compare(firstFile[key].(map[string]any), secondFile[key].(map[string]any)),
				},
			)

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
