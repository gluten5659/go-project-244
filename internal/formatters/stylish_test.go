package formatters_test

import (
	"code/internal/diff"
	"code/internal/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStylishFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodes          []diff.Node
		expectedOutput string
	}{
		{
			name:           "empty diff renders empty braces",
			nodes:          nil,
			expectedOutput: "{\n}",
		},
		{
			name: "each change kind gets its marker at the root level",
			nodes: []diff.Node{
				{Kind: diff.Deleted, Key: "follow", OldValue: false},
				{Kind: diff.Unchanged, Key: "host", OldValue: "hexlet.io", NewValue: "hexlet.io"},
				{Kind: diff.Added, Key: "verbose", NewValue: true},
			},
			expectedOutput: "{\n" +
				"  - follow: false\n" +
				"    host: hexlet.io\n" +
				"  + verbose: true\n" +
				"}",
		},
		{
			name:           "nil value renders as null",
			nodes:          []diff.Node{{Kind: diff.Added, Key: "setting3", NewValue: nil}},
			expectedOutput: "{\n  + setting3: null\n}",
		},
		{
			name: "updated value prints the old marker then the new one",
			nodes: []diff.Node{
				{Kind: diff.Updated, Key: "timeout", OldValue: 50, NewValue: 20},
			},
			expectedOutput: "{\n" +
				"  - timeout: 50\n" +
				"  + timeout: 20\n" +
				"}",
		},
		{
			name: "an added object renders as a sorted tree of its keys",
			nodes: []diff.Node{
				{Kind: diff.Added, Key: "settings", NewValue: map[string]any{"b": 2, "a": 1}},
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
			nodes: []diff.Node{
				{Kind: diff.Nested, Key: "common", Children: []diff.Node{
					{Kind: diff.Added, Key: "follow", NewValue: false},
					{Kind: diff.Nested, Key: "sub", Children: []diff.Node{
						{Kind: diff.Deleted, Key: "x", OldValue: 1},
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

			formatter, err := formatters.New(formatters.Stylish)
			require.NoError(t, err)

			formatted, err := formatter.Format(testCase.nodes)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedOutput, formatted)
		})
	}
}
