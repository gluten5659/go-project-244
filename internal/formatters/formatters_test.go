package formatters_test

import (
	"code/internal/diff"
	"code/internal/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatRejectsUnsupportedFormat(t *testing.T) {
	t.Parallel()

	formatted, err := formatters.Format(
		[]diff.Node{{Kind: diff.Added, Key: "x", Value: 1}},
		"bogus",
	)

	require.ErrorIs(t, err, formatters.ErrUnsupportedFormat)
	assert.Empty(t, formatted)
}

func TestSupportedNames(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		[]string{formatters.JSON, formatters.Plain, formatters.Stylish},
		formatters.SupportedNames(),
	)
}
