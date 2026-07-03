package compare_test

import (
	"code/internal/compare"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	keyHost     = "host"
	keyTimeout  = "timeout"
	keyList     = "list"
	keyNested   = "nested"
	hostValue   = "hexlet.io"
	scalarValue = "str"
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
			expectedDiffs: []compare.Diff{},
		},
		{
			name:       "key only in first file is deleted with its value",
			firstFile:  map[string]any{"follow": false},
			secondFile: map[string]any{},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Deleted, Key: "follow", Value: false},
			},
		},
		{
			name:       "key only in second file is added with its value",
			firstFile:  map[string]any{},
			secondFile: map[string]any{"verbose": true},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Added, Key: "verbose", Value: true},
			},
		},
		{
			name:       "changed value is deleted then added",
			firstFile:  map[string]any{keyTimeout: 50},
			secondFile: map[string]any{keyTimeout: 20},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Deleted, Key: keyTimeout, Value: 50},
				{Kind: compare.Added, Key: keyTimeout, Value: 20},
			},
		},
		{
			name:       "equal value stays unchanged",
			firstFile:  map[string]any{keyHost: hostValue},
			secondFile: map[string]any{keyHost: hostValue},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Unchanged, Key: keyHost, Value: hostValue},
			},
		},
		{
			name:       "equal slice values stay unchanged",
			firstFile:  map[string]any{keyList: []any{1, 2, 3}},
			secondFile: map[string]any{keyList: []any{1, 2, 3}},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Unchanged, Key: keyList, Value: []any{1, 2, 3}},
			},
		},
		{
			name:       "different slice values are deleted then added",
			firstFile:  map[string]any{keyList: []any{1, 2, 3}},
			secondFile: map[string]any{keyList: []any{1, 2, 4}},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Deleted, Key: keyList, Value: []any{1, 2, 3}},
				{Kind: compare.Added, Key: keyList, Value: []any{1, 2, 4}},
			},
		},
		{
			name:       "different nested values recurse into a child diff",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{keyNested: map[string]any{"x": 2}},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Unchanged, Key: keyNested, Value: []compare.Diff{
					{Kind: compare.Deleted, Key: "x", Value: 1},
					{Kind: compare.Added, Key: "x", Value: 2},
				}},
			},
		},
		{
			name:       "object only in first file is deleted as a nested tree",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Deleted, Key: keyNested, Value: []compare.Diff{
					{Kind: compare.Unchanged, Key: "x", Value: 1},
				}},
			},
		},
		{
			name:       "value changed from object to scalar deletes the tree and adds the scalar",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{keyNested: scalarValue},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Deleted, Key: keyNested, Value: []compare.Diff{
					{Kind: compare.Unchanged, Key: "x", Value: 1},
				}},
				{Kind: compare.Added, Key: keyNested, Value: scalarValue},
			},
		},
		{
			name:       "value changed from scalar to object deletes the scalar and adds the tree",
			firstFile:  map[string]any{keyNested: scalarValue},
			secondFile: map[string]any{keyNested: map[string]any{"x": 1}},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Deleted, Key: keyNested, Value: scalarValue},
				{Kind: compare.Added, Key: keyNested, Value: []compare.Diff{
					{Kind: compare.Unchanged, Key: "x", Value: 1},
				}},
			},
		},
		{
			name:      "object only in second file is added as a sorted nested tree",
			firstFile: map[string]any{},
			secondFile: map[string]any{
				keyNested: map[string]any{"alpha": 1, "beta": map[string]any{"gamma": 2}},
			},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Added, Key: keyNested, Value: []compare.Diff{
					{Kind: compare.Unchanged, Key: "alpha", Value: 1},
					{Kind: compare.Unchanged, Key: "beta", Value: []compare.Diff{
						{Kind: compare.Unchanged, Key: "gamma", Value: 2},
					}},
				}},
			},
		},
		{
			name:       "keys are sorted across both files",
			firstFile:  map[string]any{"b": 1, "a": 1},
			secondFile: map[string]any{"c": 1, "a": 1},
			expectedDiffs: []compare.Diff{
				{Kind: compare.Unchanged, Key: "a", Value: 1},
				{Kind: compare.Deleted, Key: "b", Value: 1},
				{Kind: compare.Added, Key: "c", Value: 1},
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
