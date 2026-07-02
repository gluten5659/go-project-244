package code_test

import (
	"code"
	"code/internal/formatters"
	"code/internal/testutil"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	firstConfig  = `{"a": 1, "b": 2}`
	secondConfig = `{"a": 1, "b": 3}`
	stylishDiff  = "{\n" +
		"    a: 1\n" +
		"  - b: 2\n" +
		"  + b: 3\n" +
		"}"
)

func writeConfigs(tb testing.TB) (string, string) {
	tb.Helper()

	return testutil.WriteTempFileNamed(tb, "first.json", firstConfig),
		testutil.WriteTempFileNamed(tb, "second.json", secondConfig)
}

func TestGenDiffRendersStylishDiff(t *testing.T) {
	t.Parallel()

	firstPath, secondPath := writeConfigs(t)

	result, err := code.GenDiff(firstPath, secondPath, "stylish")

	require.NoError(t, err)
	assert.Equal(t, stylishDiff, result)
}

func TestGenDiffSupportsEveryFormat(t *testing.T) {
	t.Parallel()

	firstPath, secondPath := writeConfigs(t)

	for _, format := range formatters.SupportedNames() {
		t.Run(format, func(t *testing.T) {
			t.Parallel()

			result, err := code.GenDiff(firstPath, secondPath, format)

			require.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}

func TestGenDiffRejectsUnsupportedFormat(t *testing.T) {
	t.Parallel()

	firstPath, secondPath := writeConfigs(t)

	_, err := code.GenDiff(firstPath, secondPath, "bogus")

	require.ErrorIs(t, err, formatters.ErrUnsupportedFormat)
}

func TestGenDiffReportsMissingFile(t *testing.T) {
	t.Parallel()

	firstPath, secondPath := writeConfigs(t)
	missingPath := filepath.Join(t.TempDir(), "missing.json")

	testCases := map[string]struct {
		firstPath  string
		secondPath string
	}{
		"missing first file":  {firstPath: missingPath, secondPath: secondPath},
		"missing second file": {firstPath: firstPath, secondPath: missingPath},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := code.GenDiff(testCase.firstPath, testCase.secondPath, "stylish")

			require.ErrorIs(t, err, fs.ErrNotExist)
		})
	}
}
