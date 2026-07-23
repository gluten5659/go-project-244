package diff_test

import (
	"code/internal/diff"
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
		expectedNodes []diff.Node
	}{
		{
			name:          "both files empty",
			firstFile:     map[string]any{},
			secondFile:    map[string]any{},
			expectedNodes: []diff.Node{},
		},
		{
			name:       "key only in first file is deleted with its old value",
			firstFile:  map[string]any{"follow": false},
			secondFile: map[string]any{},
			expectedNodes: []diff.Node{
				{Kind: diff.Deleted, Key: "follow", OldValue: false},
			},
		},
		{
			name:       "key only in second file is added with its new value",
			firstFile:  map[string]any{},
			secondFile: map[string]any{"verbose": true},
			expectedNodes: []diff.Node{
				{Kind: diff.Added, Key: "verbose", NewValue: true},
			},
		},
		{
			name:       "changed scalar becomes a single update with both values",
			firstFile:  map[string]any{keyTimeout: 50},
			secondFile: map[string]any{keyTimeout: 20},
			expectedNodes: []diff.Node{
				{Kind: diff.Updated, Key: keyTimeout, OldValue: 50, NewValue: 20},
			},
		},
		{
			name:       "equal value stays unchanged and carries both sides",
			firstFile:  map[string]any{keyHost: hostValue},
			secondFile: map[string]any{keyHost: hostValue},
			expectedNodes: []diff.Node{
				{Kind: diff.Unchanged, Key: keyHost, OldValue: hostValue, NewValue: hostValue},
			},
		},
		{
			name:       "equal slice values stay unchanged",
			firstFile:  map[string]any{keyList: []any{1, 2, 3}},
			secondFile: map[string]any{keyList: []any{1, 2, 3}},
			expectedNodes: []diff.Node{
				{
					Kind:     diff.Unchanged,
					Key:      keyList,
					OldValue: []any{1, 2, 3},
					NewValue: []any{1, 2, 3},
				},
			},
		},
		{
			name:       "different slice values become a single update",
			firstFile:  map[string]any{keyList: []any{1, 2, 3}},
			secondFile: map[string]any{keyList: []any{1, 2, 4}},
			expectedNodes: []diff.Node{
				{Kind: diff.Updated, Key: keyList, OldValue: []any{1, 2, 3}, NewValue: []any{1, 2, 4}},
			},
		},
		{
			name:       "two objects recurse into a nested node with child diffs",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{keyNested: map[string]any{"x": 2}},
			expectedNodes: []diff.Node{
				{Kind: diff.Nested, Key: keyNested, Children: []diff.Node{
					{Kind: diff.Updated, Key: "x", OldValue: 1, NewValue: 2},
				}},
			},
		},
		{
			name:       "object only in first file is deleted with its raw value",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{},
			expectedNodes: []diff.Node{
				{Kind: diff.Deleted, Key: keyNested, OldValue: map[string]any{"x": 1}},
			},
		},
		{
			name:       "value changed from object to scalar is an update carrying both",
			firstFile:  map[string]any{keyNested: map[string]any{"x": 1}},
			secondFile: map[string]any{keyNested: scalarValue},
			expectedNodes: []diff.Node{
				{
					Kind:     diff.Updated,
					Key:      keyNested,
					OldValue: map[string]any{"x": 1},
					NewValue: scalarValue,
				},
			},
		},
		{
			name:       "value changed from scalar to object is an update carrying both",
			firstFile:  map[string]any{keyNested: scalarValue},
			secondFile: map[string]any{keyNested: map[string]any{"x": 1}},
			expectedNodes: []diff.Node{
				{
					Kind:     diff.Updated,
					Key:      keyNested,
					OldValue: scalarValue,
					NewValue: map[string]any{"x": 1},
				},
			},
		},
		{
			name:      "object only in second file is added with its raw value",
			firstFile: map[string]any{},
			secondFile: map[string]any{
				keyNested: map[string]any{"alpha": 1, "beta": map[string]any{"gamma": 2}},
			},
			expectedNodes: []diff.Node{
				{Kind: diff.Added, Key: keyNested, NewValue: map[string]any{
					"alpha": 1, "beta": map[string]any{"gamma": 2},
				}},
			},
		},
		{
			name:       "keys are sorted across both files",
			firstFile:  map[string]any{"b": 1, "a": 1},
			secondFile: map[string]any{"c": 1, "a": 1},
			expectedNodes: []diff.Node{
				{Kind: diff.Unchanged, Key: "a", OldValue: 1, NewValue: 1},
				{Kind: diff.Deleted, Key: "b", OldValue: 1},
				{Kind: diff.Added, Key: "c", NewValue: 1},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			nodes := diff.Compare(testCase.firstFile, testCase.secondFile)

			assert.Equal(t, testCase.expectedNodes, nodes)
		})
	}
}
