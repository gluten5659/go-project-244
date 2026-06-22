package parser_test

import (
	"code/internal/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	exampleHost = "hexlet.io"
	hostKey     = "host"
	timeoutKey  = "timeout"
	settingsKey = "settings"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		fileType       string
		content        string
		expectedConfig map[string]any
		expectedError  error
	}{
		{
			name:           "json flat object",
			fileType:       parser.TypeJSON,
			content:        `{"host": "hexlet.io", "timeout": 50}`,
			expectedConfig: map[string]any{hostKey: exampleHost, timeoutKey: float64(50)},
		},
		{
			name:           "json empty object",
			fileType:       parser.TypeJSON,
			content:        `{}`,
			expectedConfig: map[string]any{},
		},
		{
			name:           "json nested object parses into a nested map with float64 numbers",
			fileType:       parser.TypeJSON,
			content:        `{"settings": {"timeout": 50}}`,
			expectedConfig: map[string]any{settingsKey: map[string]any{timeoutKey: float64(50)}},
		},
		{
			name:          "malformed json",
			fileType:      parser.TypeJSON,
			content:       `{`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "json array is not an object",
			fileType:      parser.TypeJSON,
			content:       `[1, 2, 3]`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "empty json input",
			fileType:      parser.TypeJSON,
			content:       ``,
			expectedError: parser.ErrParse,
		},
		{
			name:           "yaml flat object parses numbers into int",
			fileType:       parser.TypeYAML,
			content:        "host: hexlet.io\ntimeout: 50",
			expectedConfig: map[string]any{hostKey: exampleHost, timeoutKey: 50},
		},
		{
			name:           "yml extension parses like yaml",
			fileType:       parser.TypeYML,
			content:        "host: hexlet.io\ntimeout: 50",
			expectedConfig: map[string]any{hostKey: exampleHost, timeoutKey: 50},
		},
		{
			name:           "yaml nested object parses into a nested map",
			fileType:       parser.TypeYAML,
			content:        "settings:\n  timeout: 50",
			expectedConfig: map[string]any{settingsKey: map[string]any{timeoutKey: 50}},
		},
		{
			name:          "yaml scalar is not a mapping",
			fileType:      parser.TypeYAML,
			content:       `just a string`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "malformed yaml",
			fileType:      parser.TypeYAML,
			content:       `key: "unterminated`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "unsupported file type",
			fileType:      "txt",
			content:       `host: hexlet.io`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "empty file type",
			fileType:      "",
			content:       `{}`,
			expectedError: parser.ErrParse,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			config, err := parser.Parse(testCase.fileType, []byte(testCase.content))

			if testCase.expectedError != nil {
				require.ErrorIs(t, err, testCase.expectedError)
				assert.Nil(t, config)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedConfig, config)
		})
	}
}
