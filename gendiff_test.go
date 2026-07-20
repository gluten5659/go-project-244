package code_test

import (
	"code"
	"code/internal/formatters"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	firstFile  = "testdata/first.json"
	secondFile = "testdata/second.json"
)

func readGolden(tb testing.TB, path string) string {
	tb.Helper()

	content, err := os.ReadFile(path)
	require.NoError(tb, err)

	return strings.TrimSuffix(string(content), "\n")
}

func TestGenDiffRendersEveryFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		format     string
		goldenFile string
	}{
		{name: "stylish", format: formatters.Stylish, goldenFile: "testdata/expected/basic.stylish"},
		{name: "plain", format: formatters.Plain, goldenFile: "testdata/expected/basic.plain"},
		{name: "json", format: formatters.JSON, goldenFile: "testdata/expected/basic.json"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := code.GenDiff(firstFile, secondFile, testCase.format)

			require.NoError(t, err)
			assert.Equal(t, readGolden(t, testCase.goldenFile), result)
		})
	}
}

func TestGenDiffTreatsMatchingJSONAndYAMLAsEqual(t *testing.T) {
	t.Parallel()

	result, err := code.GenDiff(
		"testdata/identical.json",
		"testdata/identical.yaml",
		formatters.Stylish,
	)

	require.NoError(t, err)
	assert.Equal(t, readGolden(t, "testdata/expected/identical.stylish"), result)
}
