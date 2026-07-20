package formatters_test

import (
	"code/internal/diff"
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
		nodes          []diff.Node
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
			nodes: []diff.Node{
				{Kind: diff.Deleted, Key: "gone", Value: 5},
				{Kind: diff.Updated, Key: "x", OldValue: 1, NewValue: 2},
				{Kind: diff.Added, Key: "y", Value: true},
				{Kind: diff.Added, Key: "nothing", Value: nil},
				{Kind: diff.Unchanged, Key: "z", Value: "keep"},
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
			nodes: []diff.Node{
				{Kind: diff.Nested, Key: "parent", Children: []diff.Node{
					{Kind: diff.Added, Key: "leaf", Value: 1},
				}},
				{Kind: diff.Added, Key: "obj", Value: map[string]any{"inner": 2}},
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
		[]diff.Node{{Kind: diff.Added, Key: "x", Value: math.NaN()}},
		formatters.JSON,
	)

	require.Error(t, err)
	assert.Empty(t, formatted)
}
