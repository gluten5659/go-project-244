package parser_test

import (
	"code/internal/parser"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumberString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		number   parser.Number
		expected string
	}{
		"integer":                  {parser.IntNumber(50), "50"},
		"negative integer":         {parser.IntNumber(-7), "-7"},
		"integer beyond int64":     {parser.UintNumber(10000000000000000000), "10000000000000000000"},
		"fractional float":         {parser.FloatNumber(1.5), "1.5"},
		"whole float keeps marker": {parser.FloatNumber(1.0), "1.0"},
		"zero float":               {parser.FloatNumber(0), "0.0"},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.number.String())
		})
	}
}

func TestNumberMarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		number   parser.Number
		expected string
	}{
		"integer serializes as a bare json number": {parser.IntNumber(50), "50"},
		"whole float keeps its decimal marker":     {parser.FloatNumber(1.0), "1.0"},
		"fractional float":                         {parser.FloatNumber(1.5), "1.5"},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			encoded, err := json.Marshal(testCase.number)

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, string(encoded))
		})
	}
}
