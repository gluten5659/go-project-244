package files_test

import (
	"code/internal/files"
	"code/internal/parser"
	"code/internal/testutil"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadExistingFile(t *testing.T) {
	t.Parallel()

	fileContent := `{"key": "value"}`
	path := testutil.WriteTempFile(t, fileContent)

	fileType, content, err := files.Read(path)

	require.NoError(t, err)
	assert.Equal(t, parser.TypeJSON, fileType)
	assert.Equal(t, fileContent, string(content))
}

func TestReadMissingFile(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")

	_, content, err := files.Read(missingPath)

	require.ErrorIs(t, err, files.ErrRead)
	require.ErrorIs(t, err, fs.ErrNotExist)
	assert.Nil(t, content)
}

func TestReadDeterminesFileType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		fileName         string
		expectedFileType string
	}{
		{name: "json extension", fileName: "config.json", expectedFileType: parser.TypeJSON},
		{name: "yaml extension", fileName: "config.yaml", expectedFileType: parser.TypeYAML},
		{name: "yml extension", fileName: "config.yml", expectedFileType: parser.TypeYML},
		{name: "no extension", fileName: "config", expectedFileType: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			path := testutil.WriteTempFileNamed(t, testCase.fileName, "")

			fileType, _, err := files.Read(path)

			require.NoError(t, err)
			assert.Equal(t, testCase.expectedFileType, fileType)
		})
	}
}
