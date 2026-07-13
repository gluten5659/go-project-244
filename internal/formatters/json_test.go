package formatters_test

import (
	"code/internal/compare"
	"code/internal/formatters"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodes          []compare.Node
		expectedOutput string
	}{
		{
			name:  "empty diff is wrapped in a diff object",
			nodes: nil,
			expectedOutput: `{
  "diff": []
}`,
		},
		{
			name: "each change kind maps to a typed node",
			nodes: []compare.Node{
				{Kind: compare.Deleted, Key: "gone", Value: 5},
				{Kind: compare.Updated, Key: "x", OldValue: 1, NewValue: 2},
				{Kind: compare.Added, Key: "y", Value: true},
				{Kind: compare.Added, Key: "nothing", Value: nil},
				{Kind: compare.Unchanged, Key: "z", Value: "keep"},
			},
			expectedOutput: `{
  "diff": [
    {
      "key": "gone",
      "type": "removed",
      "value": 5
    },
    {
      "key": "x",
      "newValue": 2,
      "oldValue": 1,
      "type": "updated"
    },
    {
      "key": "y",
      "type": "added",
      "value": true
    },
    {
      "key": "nothing",
      "type": "added",
      "value": null
    },
    {
      "key": "z",
      "type": "unchanged",
      "value": "keep"
    }
  ]
}`,
		},
		{
			name: "nested objects carry children while whole values are inlined",
			nodes: []compare.Node{
				{Kind: compare.Nested, Key: "parent", Children: []compare.Node{
					{Kind: compare.Added, Key: "leaf", Value: 1},
				}},
				{Kind: compare.Added, Key: "obj", Value: map[string]any{"inner": 2}},
			},
			expectedOutput: `{
  "diff": [
    {
      "children": [
        {
          "key": "leaf",
          "type": "added",
          "value": 1
        }
      ],
      "key": "parent",
      "type": "nested"
    },
    {
      "key": "obj",
      "type": "added",
      "value": {
        "inner": 2
      }
    }
  ]
}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			formatted, err := formatters.Format(testCase.nodes, formatters.JSON)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedOutput, formatted)
		})
	}
}

func TestJSONFormatReportsMarshallingFailure(t *testing.T) {
	t.Parallel()

	formatted, err := formatters.Format(
		[]compare.Node{{Kind: compare.Added, Key: "x", Value: math.NaN()}},
		formatters.JSON,
	)

	require.Error(t, err)
	assert.Empty(t, formatted)
}
