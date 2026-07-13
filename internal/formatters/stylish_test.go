package formatters_test

import (
	"code/internal/compare"
	"code/internal/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStylishFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodes          []compare.Node
		expectedOutput string
	}{
		{
			name:           "empty diff renders empty braces",
			nodes:          nil,
			expectedOutput: "{\n}",
		},
		{
			name: "each change kind gets its marker at the root level",
			nodes: []compare.Node{
				{Kind: compare.Deleted, Key: "follow", Value: false},
				{Kind: compare.Unchanged, Key: "host", Value: "hexlet.io"},
				{Kind: compare.Added, Key: "verbose", Value: true},
			},
			expectedOutput: "{\n" +
				"  - follow: false\n" +
				"    host: hexlet.io\n" +
				"  + verbose: true\n" +
				"}",
		},
		{
			name:           "nil value renders as null",
			nodes:          []compare.Node{{Kind: compare.Added, Key: "setting3", Value: nil}},
			expectedOutput: "{\n  + setting3: null\n}",
		},
		{
			name: "updated value prints the old marker then the new one",
			nodes: []compare.Node{
				{Kind: compare.Updated, Key: "timeout", OldValue: 50, NewValue: 20},
			},
			expectedOutput: "{\n" +
				"  - timeout: 50\n" +
				"  + timeout: 20\n" +
				"}",
		},
		{
			name: "an added object renders as a sorted tree of its keys",
			nodes: []compare.Node{
				{Kind: compare.Added, Key: "settings", Value: map[string]any{"b": 2, "a": 1}},
			},
			expectedOutput: "{\n" +
				"  + settings: {\n" +
				"        a: 1\n" +
				"        b: 2\n" +
				"    }\n" +
				"}",
		},
		{
			name: "nested children indent by depth",
			nodes: []compare.Node{
				{Kind: compare.Nested, Key: "common", Children: []compare.Node{
					{Kind: compare.Added, Key: "follow", Value: false},
					{Kind: compare.Nested, Key: "sub", Children: []compare.Node{
						{Kind: compare.Deleted, Key: "x", Value: 1},
					}},
				}},
			},
			expectedOutput: "{\n" +
				"    common: {\n" +
				"      + follow: false\n" +
				"        sub: {\n" +
				"          - x: 1\n" +
				"        }\n" +
				"    }\n" +
				"}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			formatted, err := formatters.Format(testCase.nodes, formatters.Stylish)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedOutput, formatted)
		})
	}
}
