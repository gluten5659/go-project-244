package formatters_test

import (
	"code/internal/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRejectsUnsupportedFormat(t *testing.T) {
	t.Parallel()

	formatter, err := formatters.New("bogus")

	require.ErrorIs(t, err, formatters.ErrUnsupportedFormat)
	assert.Nil(t, formatter)
}

func TestListSupportedNames(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		[]string{formatters.JSON, formatters.Plain, formatters.Stylish},
		formatters.ListSupportedNames(),
	)
}

func TestNewBuildsEverySupportedName(t *testing.T) {
	t.Parallel()

	for _, name := range formatters.ListSupportedNames() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			formatter, err := formatters.New(name)

			require.NoError(t, err)
			assert.NotNil(t, formatter)
		})
	}
}
