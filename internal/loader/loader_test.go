package loader_test

import (
	"code/internal/files"
	"code/internal/loader"
	"code/internal/parser"
	"code/internal/testutil"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromFileParsesContent(t *testing.T) {
	t.Parallel()

	path := testutil.WriteTempFile(t, `{"host": "hexlet.io", "timeout": 50}`)

	values, err := loader.FromFile(path)

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"host": "hexlet.io", "timeout": float64(50)}, values)
}

func TestFromFileMissingFile(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")

	values, err := loader.FromFile(missingPath)

	require.ErrorIs(t, err, files.ErrRead)
	require.ErrorIs(t, err, fs.ErrNotExist)
	assert.Nil(t, values)
}

func TestFromFileMalformedContent(t *testing.T) {
	t.Parallel()

	path := testutil.WriteTempFile(t, `{`)

	values, err := loader.FromFile(path)

	require.ErrorIs(t, err, parser.ErrParse)
	assert.Nil(t, values)
}

func TestFromFileParsesYAMLByExtension(t *testing.T) {
	t.Parallel()

	path := testutil.WriteTempFileNamed(t, "config.yaml", "host: hexlet.io\ntimeout: 50")

	values, err := loader.FromFile(path)

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"host": "hexlet.io", "timeout": 50}, values)
}
