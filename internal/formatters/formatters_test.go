package formatters_test

import (
	"code/internal/compare"
	"code/internal/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		format         string
		diffs          []compare.Diff
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "stylish renders empty braces for an empty diff",
			format:         formatters.Stylish,
			diffs:          nil,
			expectedOutput: "{\n}",
		},
		{
			name:   "stylish gives each change kind its marker at the root level",
			format: formatters.Stylish,
			diffs: []compare.Diff{
				{Change: compare.Deleted, Key: "follow", Value: false},
				{Change: compare.NoChanges, Key: "host", Value: "hexlet.io"},
				{Change: compare.Added, Key: "verbose", Value: true},
			},
			expectedOutput: "{\n" +
				"  - follow: false\n" +
				"    host: hexlet.io\n" +
				"  + verbose: true\n" +
				"}",
		},
		{
			name:           "stylish renders a nil value as null",
			format:         formatters.Stylish,
			diffs:          []compare.Diff{{Change: compare.Added, Key: "setting3", Value: nil}},
			expectedOutput: "{\n  + setting3: null\n}",
		},
		{
			name:   "stylish indents nested children by depth",
			format: formatters.Stylish,
			diffs: []compare.Diff{
				{Change: compare.NoChanges, Key: "common", Value: []compare.Diff{
					{Change: compare.Added, Key: "follow", Value: false},
					{Change: compare.NoChanges, Key: "sub", Value: []compare.Diff{
						{Change: compare.Deleted, Key: "x", Value: 1},
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
		{
			name:   "plain reports added values with quoting rules",
			format: formatters.Plain,
			diffs: []compare.Diff{
				{Change: compare.Added, Key: "flag", Value: false},
				{Change: compare.Added, Key: "name", Value: "bob"},
				{Change: compare.Added, Key: "opt", Value: nil},
			},
			expectedOutput: "Property 'flag' was added with value: false\n" +
				"Property 'name' was added with value: 'bob'\n" +
				"Property 'opt' was added with value: null",
		},
		{
			name:           "plain reports a removed property",
			format:         formatters.Plain,
			diffs:          []compare.Diff{{Change: compare.Deleted, Key: "old", Value: 1}},
			expectedOutput: "Property 'old' was removed",
		},
		{
			name:   "plain joins a deleted and added pair into an update",
			format: formatters.Plain,
			diffs: []compare.Diff{
				{Change: compare.Deleted, Key: "x", Value: 1},
				{Change: compare.Added, Key: "x", Value: 2},
			},
			expectedOutput: "Property 'x' was updated. From 1 to 2",
		},
		{
			name:   "plain renders nested values as a complex placeholder",
			format: formatters.Plain,
			diffs: []compare.Diff{
				{Change: compare.Added, Key: "obj", Value: []compare.Diff{
					{Change: compare.NoChanges, Key: "a", Value: 1},
				}},
			},
			expectedOutput: "Property 'obj' was added with value: [complex value]",
		},
		{
			name:   "plain dots the path and skips unchanged leaves",
			format: formatters.Plain,
			diffs: []compare.Diff{
				{Change: compare.NoChanges, Key: "common", Value: []compare.Diff{
					{Change: compare.NoChanges, Key: "keep", Value: "v"},
					{Change: compare.Added, Key: "new", Value: true},
				}},
			},
			expectedOutput: "Property 'common.new' was added with value: true",
		},
		{
			name:        "unsupported format returns an error",
			format:      "bogus",
			diffs:       []compare.Diff{{Change: compare.Added, Key: "x", Value: 1}},
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			formatted, err := formatters.Format(testCase.diffs, testCase.format)

			if testCase.expectError {
				require.Error(t, err)
				assert.Empty(t, formatted)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedOutput, formatted)
		})
	}
}
