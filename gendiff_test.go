package code_test

import (
	"code"
	"code/internal/formatters"
	"code/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	firstConfig  = `{"a": 1, "b": 2}`
	secondConfig = `{"a": 1, "b": 3}`

	stylishDiff = "{\n" +
		"    a: 1\n" +
		"  - b: 2\n" +
		"  + b: 3\n" +
		"}"

	plainDiff = "Property 'b' was updated. From 2 to 3"

	jsonDiff = "{\n" +
		"  \"diff\": [\n" +
		"    {\n" +
		"      \"key\": \"a\",\n" +
		"      \"type\": \"unchanged\",\n" +
		"      \"value\": 1\n" +
		"    },\n" +
		"    {\n" +
		"      \"key\": \"b\",\n" +
		"      \"newValue\": 3,\n" +
		"      \"oldValue\": 2,\n" +
		"      \"type\": \"updated\"\n" +
		"    }\n" +
		"  ]\n" +
		"}"

	identicalDiff = "{\n" +
		"    host: x\n" +
		"    timeout: 50\n" +
		"}"
)

func writeConfigs(tb testing.TB) (string, string) {
	tb.Helper()

	return testutil.WriteTempFileNamed(tb, "first.json", firstConfig),
		testutil.WriteTempFileNamed(tb, "second.json", secondConfig)
}

func TestGenDiffRendersEveryFormat(t *testing.T) {
	t.Parallel()

	firstPath, secondPath := writeConfigs(t)

	testCases := map[string]string{
		formatters.Stylish: stylishDiff,
		formatters.Plain:   plainDiff,
		formatters.JSON:    jsonDiff,
	}

	for format, expectedOutput := range testCases {
		t.Run(format, func(t *testing.T) {
			t.Parallel()

			result, err := code.GenDiff(firstPath, secondPath, format)

			require.NoError(t, err)
			assert.Equal(t, expectedOutput, result)
		})
	}
}

func TestGenDiffTreatsMatchingJSONAndYAMLAsEqual(t *testing.T) {
	t.Parallel()

	jsonPath := testutil.WriteTempFileNamed(t, "config.json", `{"host": "x", "timeout": 50}`)
	yamlPath := testutil.WriteTempFileNamed(t, "config.yaml", "host: x\ntimeout: 50\n")

	result, err := code.GenDiff(jsonPath, yamlPath, formatters.Stylish)

	require.NoError(t, err)
	assert.Equal(t, identicalDiff, result)
}
