package loader_test

import (
	"code/internal/loader"
	"code/internal/testutil"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	jsonConfigName = "config.json"
	yamlConfigName = "config.yaml"
	hostKey        = "host"
	hostValue      = "hexlet.io"
	settingsKey    = "settings"
	timeoutKey     = "timeout"
)

func TestFromFileParsesContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		fileName      string
		content       string
		expectedValue map[string]any
	}{
		{
			name:          "json flat object normalizes numbers",
			fileName:      jsonConfigName,
			content:       `{"host": "hexlet.io", "timeout": 50}`,
			expectedValue: map[string]any{hostKey: hostValue, timeoutKey: loader.IntNumber(50)},
		},
		{
			name:          "json empty object",
			fileName:      jsonConfigName,
			content:       `{}`,
			expectedValue: map[string]any{},
		},
		{
			name:          "json nested object",
			fileName:      jsonConfigName,
			content:       `{"settings": {"timeout": 50}}`,
			expectedValue: map[string]any{settingsKey: map[string]any{timeoutKey: loader.IntNumber(50)}},
		},
		{
			name:          "yaml flat object normalizes numbers the same as json",
			fileName:      yamlConfigName,
			content:       "host: hexlet.io\ntimeout: 50",
			expectedValue: map[string]any{hostKey: hostValue, timeoutKey: loader.IntNumber(50)},
		},
		{
			name:          "yml extension parses like yaml",
			fileName:      "config.yml",
			content:       "host: hexlet.io\ntimeout: 50",
			expectedValue: map[string]any{hostKey: hostValue, timeoutKey: loader.IntNumber(50)},
		},
		{
			name:          "yaml nested object",
			fileName:      yamlConfigName,
			content:       "settings:\n  timeout: 50",
			expectedValue: map[string]any{settingsKey: map[string]any{timeoutKey: loader.IntNumber(50)}},
		},
		{
			name:          "yaml non-string keys normalize into string-keyed maps",
			fileName:      yamlConfigName,
			content:       "settings:\n  1: one\n  2: two",
			expectedValue: map[string]any{settingsKey: map[string]any{"1": "one", "2": "two"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			path := testutil.WriteTempFileNamed(t, testCase.fileName, testCase.content)

			values, err := loader.FromFile(path)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedValue, values)
		})
	}
}

func TestFromFileNormalizesNumbersToACommonForm(t *testing.T) {
	t.Parallel()

	load := func(tb testing.TB, fileName, content string) map[string]any {
		tb.Helper()

		values, err := loader.FromFile(testutil.WriteTempFileNamed(tb, fileName, content))
		require.NoError(tb, err)

		return values
	}

	t.Run("identical numbers from json and yaml compare equal", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"timeout": 50, "ratio": 1.5}`)
		fromYAML := load(t, yamlConfigName, "timeout: 50\nratio: 1.5")

		assert.Equal(t, fromJSON, fromYAML)
	})

	t.Run("an integer and a float are different values", func(t *testing.T) {
		t.Parallel()

		asInteger := load(t, yamlConfigName, "timeout: 1")
		asFloat := load(t, yamlConfigName, "timeout: 1.0")

		assert.NotEqual(t, asInteger, asFloat)
	})

	t.Run("a json float parses as a float distinct from a json integer", func(t *testing.T) {
		t.Parallel()

		asFloat := load(t, jsonConfigName, `{"timeout": 1.0}`)
		asInteger := load(t, jsonConfigName, `{"timeout": 1}`)

		assert.Equal(t, map[string]any{timeoutKey: loader.FloatNumber(1.0)}, asFloat)
		assert.NotEqual(t, asInteger, asFloat)
	})

	t.Run("large integers keep their exact value", func(t *testing.T) {
		t.Parallel()

		values := load(t, jsonConfigName, `{"id": 9007199254740993}`)

		assert.Equal(t, map[string]any{"id": loader.IntNumber(9007199254740993)}, values)
	})
}

func TestFromFileRejectsUnparsableContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		fileName string
		content  string
	}{
		{name: "malformed json", fileName: jsonConfigName, content: `{`},
		{name: "json array is not an object", fileName: jsonConfigName, content: `[1, 2, 3]`},
		{name: "empty json input", fileName: jsonConfigName, content: ``},
		{name: "yaml scalar is not a mapping", fileName: yamlConfigName, content: `just a string`},
		{name: "malformed yaml", fileName: yamlConfigName, content: `key: "unterminated`},
		{name: "unsupported extension", fileName: "config.txt", content: `host: hexlet.io`},
		{name: "no extension", fileName: "config", content: `{}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			path := testutil.WriteTempFileNamed(t, testCase.fileName, testCase.content)

			values, err := loader.FromFile(path)

			require.ErrorIs(t, err, loader.ErrParse)
			assert.Nil(t, values)
		})
	}
}

func TestFromFileReportsMissingFile(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")

	values, err := loader.FromFile(missingPath)

	require.ErrorIs(t, err, loader.ErrRead)
	require.ErrorIs(t, err, fs.ErrNotExist)
	assert.Nil(t, values)
}
