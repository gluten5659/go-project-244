package parser_test

import (
	"code/internal/parser"
	"code/internal/testutil"
	"io/fs"
	"math"
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

func TestParseFileParsesContent(t *testing.T) {
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
			expectedValue: map[string]any{hostKey: hostValue, timeoutKey: int64(50)},
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
			expectedValue: map[string]any{settingsKey: map[string]any{timeoutKey: int64(50)}},
		},
		{
			name:          "yaml flat object normalizes numbers the same as json",
			fileName:      yamlConfigName,
			content:       "host: hexlet.io\ntimeout: 50",
			expectedValue: map[string]any{hostKey: hostValue, timeoutKey: int64(50)},
		},
		{
			name:          "yml extension parses like yaml",
			fileName:      "config.yml",
			content:       "host: hexlet.io\ntimeout: 50",
			expectedValue: map[string]any{hostKey: hostValue, timeoutKey: int64(50)},
		},
		{
			name:          "yaml nested object",
			fileName:      yamlConfigName,
			content:       "settings:\n  timeout: 50",
			expectedValue: map[string]any{settingsKey: map[string]any{timeoutKey: int64(50)}},
		},
		{
			name:          "yaml non-string keys normalize into string-keyed maps",
			fileName:      yamlConfigName,
			content:       "settings:\n  1: one\n  2: two",
			expectedValue: map[string]any{settingsKey: map[string]any{"1": "one", "2": "two"}},
		},
		{
			name:          "array values normalize element by element",
			fileName:      jsonConfigName,
			content:       `{"ports": [80, 443]}`,
			expectedValue: map[string]any{"ports": []any{int64(80), int64(443)}},
		},
		{
			name:          "a whole float normalizes to an integer",
			fileName:      jsonConfigName,
			content:       `{"ratio": 0.0}`,
			expectedValue: map[string]any{"ratio": int64(0)},
		},
		{
			name:          "a fractional value stays a float",
			fileName:      jsonConfigName,
			content:       `{"ratio": 1.5}`,
			expectedValue: map[string]any{"ratio": 1.5},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			path := testutil.WriteTempFileNamed(t, testCase.fileName, testCase.content)

			values, err := parser.ParseFile(path)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedValue, values)
		})
	}
}

func TestParseFileNormalizesNumbersToACommonForm(t *testing.T) {
	t.Parallel()

	load := func(tb testing.TB, fileName, content string) map[string]any {
		tb.Helper()

		values, err := parser.ParseFile(testutil.WriteTempFileNamed(tb, fileName, content))
		require.NoError(tb, err)

		return values
	}

	t.Run("identical numbers from json and yaml compare equal", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"timeout": 50, "ratio": 1.5}`)
		fromYAML := load(t, yamlConfigName, "timeout: 50\nratio: 1.5")

		assert.Equal(t, fromJSON, fromYAML)
	})

	t.Run("a whole float and the same integer are one value", func(t *testing.T) {
		t.Parallel()

		asInteger := load(t, yamlConfigName, "timeout: 1")
		asFloat := load(t, yamlConfigName, "timeout: 1.0")

		assert.Equal(t, map[string]any{timeoutKey: int64(1)}, asInteger)
		assert.Equal(t, asInteger, asFloat)
	})

	t.Run("a json whole float equals the same json integer", func(t *testing.T) {
		t.Parallel()

		asFloat := load(t, jsonConfigName, `{"timeout": 1.0}`)
		asInteger := load(t, jsonConfigName, `{"timeout": 1}`)

		assert.Equal(t, map[string]any{timeoutKey: int64(1)}, asFloat)
		assert.Equal(t, asInteger, asFloat)
	})

	t.Run("a whole float from yaml equals an integer from json", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"timeout": 1}`)
		fromYAML := load(t, yamlConfigName, "timeout: 1.0")

		assert.Equal(t, fromJSON, fromYAML)
	})

	t.Run("the rule holds inside nested objects", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"settings": {"timeout": 1.0}}`)
		fromYAML := load(t, yamlConfigName, "settings:\n  timeout: 1")

		assert.Equal(t, map[string]any{settingsKey: map[string]any{timeoutKey: int64(1)}}, fromJSON)
		assert.Equal(t, fromJSON, fromYAML)
	})

	t.Run("the rule holds inside arrays", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"ports": [1.0, 2, 2.5]}`)
		fromYAML := load(t, yamlConfigName, "ports:\n  - 1\n  - 2.0\n  - 2.5")

		assert.Equal(t, map[string]any{"ports": []any{int64(1), int64(2), 2.5}}, fromJSON)
		assert.Equal(t, fromJSON, fromYAML)
	})

	t.Run("the rule holds inside arrays of objects", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"hosts": [{"timeout": 1.0}]}`)
		fromYAML := load(t, yamlConfigName, "hosts:\n  - timeout: 1")

		assert.Equal(t, map[string]any{"hosts": []any{
			map[string]any{timeoutKey: int64(1)},
		}}, fromJSON)
		assert.Equal(t, fromJSON, fromYAML)
	})

	t.Run("a whole float at the int64 floor becomes an integer", func(t *testing.T) {
		t.Parallel()

		values := load(t, jsonConfigName, `{"floor": -9223372036854775808.0}`)

		assert.Equal(t, map[string]any{"floor": int64(math.MinInt64)}, values)
	})

	t.Run("a whole float above the int64 ceiling stays a float", func(t *testing.T) {
		t.Parallel()

		values := load(t, jsonConfigName, `{"huge": 1e19}`)

		assert.Equal(t, map[string]any{"huge": 1e19}, values)
	})

	t.Run("large integers keep their exact value", func(t *testing.T) {
		t.Parallel()

		values := load(t, jsonConfigName, `{"id": 9007199254740993}`)

		assert.Equal(t, map[string]any{"id": int64(9007199254740993)}, values)
	})

	t.Run("an integer beyond int64 but within uint64 stays exact", func(t *testing.T) {
		t.Parallel()

		fromJSON := load(t, jsonConfigName, `{"big": 10000000000000000000}`)
		fromYAML := load(t, yamlConfigName, "big: 10000000000000000000\n")

		assert.Equal(t, map[string]any{"big": uint64(10000000000000000000)}, fromJSON)
		assert.Equal(t, fromJSON, fromYAML)
	})
}

func TestParseFileRejectsUnparsableContent(t *testing.T) {
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
		{name: "yaml positive infinity", fileName: yamlConfigName, content: "value: .inf"},
		{name: "yaml negative infinity", fileName: yamlConfigName, content: "value: -.inf"},
		{name: "yaml not a number", fileName: yamlConfigName, content: "value: .nan"},
		{name: "json overflowing float", fileName: jsonConfigName, content: `{"value": 1e400}`},
		{
			name:     "json integer beyond uint64",
			fileName: jsonConfigName,
			content:  `{"value": 123456789012345678901234567890}`,
		},
		{name: "non-finite inside an array", fileName: jsonConfigName, content: `{"list": [1e400]}`},
		{
			name:     "non-finite under a non-string key",
			fileName: yamlConfigName,
			content:  "settings:\n  1: .inf",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			path := testutil.WriteTempFileNamed(t, testCase.fileName, testCase.content)

			values, err := parser.ParseFile(path)

			require.ErrorIs(t, err, parser.ErrParse)
			assert.Nil(t, values)
		})
	}
}

func TestParseFileReportsMissingFile(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")

	values, err := parser.ParseFile(missingPath)

	require.ErrorIs(t, err, parser.ErrRead)
	require.ErrorIs(t, err, fs.ErrNotExist)
	assert.Nil(t, values)
}

func TestParseFileKeepsTheFilesystemErrorInTheChain(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")

	_, err := parser.ParseFile(missingPath)

	var pathError *fs.PathError

	require.ErrorAs(t, err, &pathError)
	assert.Equal(t, "open", pathError.Op)
	assert.Equal(t, missingPath, pathError.Path)
}
