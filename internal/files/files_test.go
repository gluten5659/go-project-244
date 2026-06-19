package files_test

import (
	"code/internal/files"
	"code/internal/testutil"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadExistingFile(t *testing.T) {
	t.Parallel()

	content := `{"key": "value"}`
	path := testutil.WriteTempFile(t, content)

	read, err := files.Read(path)

	require.NoError(t, err)
	assert.Equal(t, content, string(read))
}

func TestReadMissingFile(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")

	read, err := files.Read(missingPath)

	require.ErrorIs(t, err, files.ErrRead)
	require.ErrorIs(t, err, fs.ErrNotExist)
	assert.Nil(t, read)
}
