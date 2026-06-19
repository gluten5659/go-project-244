package parser_test

import (
	"code/internal/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		content        string
		expectedConfig map[string]any
		expectedError  error
	}{
		{
			name:           "flat object",
			content:        `{"host": "hexlet.io", "timeout": 50}`,
			expectedConfig: map[string]any{"host": "hexlet.io", "timeout": float64(50)},
		},
		{
			name:           "empty object",
			content:        `{}`,
			expectedConfig: map[string]any{},
		},
		{
			name:           "nested object parses into a nested map with float64 numbers",
			content:        `{"settings": {"timeout": 50}}`,
			expectedConfig: map[string]any{"settings": map[string]any{"timeout": float64(50)}},
		},
		{
			name:          "malformed json",
			content:       `{`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "json array is not an object",
			content:       `[1, 2, 3]`,
			expectedError: parser.ErrParse,
		},
		{
			name:          "empty input",
			content:       ``,
			expectedError: parser.ErrParse,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			config, err := parser.Parse([]byte(testCase.content))

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
