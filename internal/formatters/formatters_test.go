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
			name:           "plain renders nothing for an empty diff",
			format:         formatters.Plain,
			diffs:          nil,
			expectedOutput: "",
		},
		{
			name:   "plain reports added values with quoting rules",
			format: formatters.Plain,
			diffs: []compare.Diff{
				{Change: compare.Added, Key: "flag", Value: false},
				{Change: compare.Added, Key: "name", Value: "bob"},
				{Change: compare.Added, Key: "opt", Value: nil},
				{Change: compare.Added, Key: "empty", Value: ""},
			},
			expectedOutput: "Property 'flag' was added with value: false\n" +
				"Property 'name' was added with value: 'bob'\n" +
				"Property 'opt' was added with value: null\n" +
				"Property 'empty' was added with value: ''",
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
			name:           "json renders an empty array for an empty diff",
			format:         formatters.JSON,
			diffs:          nil,
			expectedOutput: "[]",
		},
		{
			name:   "json represents each change kind as a typed node",
			format: formatters.JSON,
			diffs: []compare.Diff{
				{Change: compare.Deleted, Key: "gone", Value: 5},
				{Change: compare.Deleted, Key: "x", Value: 1},
				{Change: compare.Added, Key: "x", Value: 2},
				{Change: compare.Added, Key: "y", Value: true},
				{Change: compare.Added, Key: "nothing", Value: nil},
				{Change: compare.NoChanges, Key: "z", Value: "keep"},
			},
			expectedOutput: `[
  {
    "key": "gone",
    "type": "removed",
    "value": 5
  },
  {
    "key": "x",
    "newValue": 2,
    "oldValue": 1,
    "type": "updated"
  },
  {
    "key": "y",
    "type": "added",
    "value": true
  },
  {
    "key": "nothing",
    "type": "added",
    "value": null
  },
  {
    "key": "z",
    "type": "unchanged",
    "value": "keep"
  }
]`,
		},
		{
			name:   "json nests objects with changes and collapses whole values",
			format: formatters.JSON,
			diffs: []compare.Diff{
				{Change: compare.NoChanges, Key: "parent", Value: []compare.Diff{
					{Change: compare.Added, Key: "leaf", Value: 1},
				}},
				{Change: compare.Added, Key: "obj", Value: []compare.Diff{
					{Change: compare.NoChanges, Key: "inner", Value: 2},
				}},
			},
			expectedOutput: `[
  {
    "children": [
      {
        "key": "leaf",
        "type": "added",
        "value": 1
      }
    ],
    "key": "parent",
    "type": "nested"
  },
  {
    "key": "obj",
    "type": "added",
    "value": {
      "inner": 2
    }
  }
]`,
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

func TestSupportedNames(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		[]string{formatters.JSON, formatters.Plain, formatters.Stylish},
		formatters.SupportedNames(),
	)
}
