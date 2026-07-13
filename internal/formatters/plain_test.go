package formatters_test

import (
	"code/internal/compare"
	"code/internal/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlainFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodes          []compare.Node
		expectedOutput string
	}{
		{
			name:           "empty diff renders nothing",
			nodes:          nil,
			expectedOutput: "",
		},
		{
			name: "added values follow the quoting rules",
			nodes: []compare.Node{
				{Kind: compare.Added, Key: "flag", Value: false},
				{Kind: compare.Added, Key: "name", Value: "bob"},
				{Kind: compare.Added, Key: "opt", Value: nil},
				{Kind: compare.Added, Key: "empty", Value: ""},
			},
			expectedOutput: "Property 'flag' was added with value: false\n" +
				"Property 'name' was added with value: 'bob'\n" +
				"Property 'opt' was added with value: null\n" +
				"Property 'empty' was added with value: ''",
		},
		{
			name:           "removed property is reported",
			nodes:          []compare.Node{{Kind: compare.Deleted, Key: "old", Value: 1}},
			expectedOutput: "Property 'old' was removed",
		},
		{
			name: "updated property reports both values",
			nodes: []compare.Node{
				{Kind: compare.Updated, Key: "x", OldValue: 1, NewValue: 2},
			},
			expectedOutput: "Property 'x' was updated. From 1 to 2",
		},
		{
			name: "updated object value collapses to a complex placeholder",
			nodes: []compare.Node{
				{Kind: compare.Updated, Key: "nest", OldValue: map[string]any{"k": "v"}, NewValue: "str"},
			},
			expectedOutput: "Property 'nest' was updated. From [complex value] to 'str'",
		},
		{
			name: "added object renders as a complex placeholder",
			nodes: []compare.Node{
				{Kind: compare.Added, Key: "obj", Value: map[string]any{"a": 1}},
			},
			expectedOutput: "Property 'obj' was added with value: [complex value]",
		},
		{
			name: "nested paths are dotted and unchanged leaves are skipped",
			nodes: []compare.Node{
				{Kind: compare.Nested, Key: "common", Children: []compare.Node{
					{Kind: compare.Unchanged, Key: "keep", Value: "v"},
					{Kind: compare.Added, Key: "new", Value: true},
				}},
			},
			expectedOutput: "Property 'common.new' was added with value: true",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			formatted, err := formatters.Format(testCase.nodes, formatters.Plain)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedOutput, formatted)
		})
	}
}
