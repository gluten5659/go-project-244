package output_test

import (
	"code/internal/compare"
	"code/internal/output"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDiff(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		diffs          []compare.Diff
		expectedOutput string
	}{
		{
			name:           "empty diff renders empty braces",
			diffs:          nil,
			expectedOutput: "{\n}",
		},
		{
			name: "each change kind gets its marker at the root level",
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
			name: "nil value renders as null",
			diffs: []compare.Diff{
				{Change: compare.Added, Key: "setting3", Value: nil},
			},
			expectedOutput: "{\n  + setting3: null\n}",
		},
		{
			name: "nested children indent by depth",
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			formatted := output.FormatDiff(testCase.diffs)

			assert.Equal(t, testCase.expectedOutput, formatted)
		})
	}
}
