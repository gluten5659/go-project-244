package compare_test

import (
	"code/internal/compare"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	keyHost    = "host"
	keyTimeout = "timeout"
	keyList    = "list"
	keyNested  = "nested"
	hostValue  = "hexlet.io"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		firstFile     map[string]any
		secondFile    map[string]any
		expectedDiffs []compare.Diff
	}{
		{
			name:          "both files empty",
			firstFile:     map[string]any{},
			secondFile:    map[string]any{},
			expectedDiffs: nil,
		},
		{
			name:       "key only in first file is deleted with its value",
			firstFile:  map[string]any{"follow": false},
			secondFile: map[string]any{},
			expectedDiffs: []compare.Diff{
				{Change: compare.Deleted, Key: "follow", Value: false},
			},
		},
		{
			name:       "key only in second file is added with its value",
			firstFile:  map[string]any{},
			secondFile: map[string]any{"verbose": true},
			expectedDiffs: []compare.Diff{
				{Change: compare.Added, Key: "verbose", Value: true},
			},
		},
		{
			name:       "changed value is deleted then added",
			firstFile:  map[string]any{keyTimeout: 50},
			secondFile: map[string]any{keyTimeout: 20},
			expectedDiffs: []compare.Diff{
				{Change: compare.Deleted, Key: keyTimeout, Value: 50},
				{Change: compare.Added, Key: keyTimeout, Value: 20},
			},
		},
		{
			name:       "equal value stays unchanged",
			firstFile:  map[string]any{keyHost: hostValue},
			secondFile: map[string]any{keyHost: hostValue},
			expectedDiffs: []compare.Diff{
				{Change: compare.NoChanges, Key: keyHost, Value: hostValue},
			},
		},
		{
			name:       "equal nested values stay unchanged without panic",
			firstFile:  map[string]any{keyList: []any{1, 2, 3}},
			secondFile: map[string]any{keyList: []any{1, 2, 3}},
			expectedDiffs: []compare.Diff{
				{Change: compare.NoChanges, Key: keyList, Value: []any{1, 2, 3}},
			},
		},
		{
			name:       "different nested values are deleted then added",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{keyNested: map[string]any{"x": 2}},
			expectedDiffs: []compare.Diff{
				{Change: compare.Deleted, Key: keyNested, Value: map[string]any{"x": 1}},
				{Change: compare.Added, Key: keyNested, Value: map[string]any{"x": 2}},
			},
		},
		{
			name:       "keys are sorted across both files",
			firstFile:  map[string]any{"b": 1, "a": 1},
			secondFile: map[string]any{"c": 1, "a": 1},
			expectedDiffs: []compare.Diff{
				{Change: compare.NoChanges, Key: "a", Value: 1},
				{Change: compare.Deleted, Key: "b", Value: 1},
				{Change: compare.Added, Key: "c", Value: 1},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			diffs := compare.Compare(testCase.firstFile, testCase.secondFile)

			assert.Equal(t, testCase.expectedDiffs, diffs)
		})
	}
}
