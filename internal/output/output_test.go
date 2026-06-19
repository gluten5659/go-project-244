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
			name: "added entry is prefixed with a plus",
			diffs: []compare.Diff{
				{Change: compare.Added, Key: "verbose", Value: true},
			},
			expectedOutput: "{\n  + verbose: true\n}",
		},
		{
			name: "deleted entry is prefixed with a minus",
			diffs: []compare.Diff{
				{Change: compare.Deleted, Key: "proxy", Value: "123.234.53.22"},
			},
			expectedOutput: "{\n  - proxy: 123.234.53.22\n}",
		},
		{
			name: "unchanged entry is prefixed with spaces",
			diffs: []compare.Diff{
				{Change: compare.NoChanges, Key: "host", Value: "hexlet.io"},
			},
			expectedOutput: "{\n    host: hexlet.io\n}",
		},
		{
			name: "all change kinds render together in order",
			diffs: []compare.Diff{
				{Change: compare.Deleted, Key: "follow", Value: false},
				{Change: compare.NoChanges, Key: "host", Value: "hexlet.io"},
				{Change: compare.Deleted, Key: "timeout", Value: 50},
				{Change: compare.Added, Key: "timeout", Value: 20},
			},
			expectedOutput: "{\n  - follow: false\n    host: hexlet.io\n  - timeout: 50\n  + timeout: 20\n}",
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
