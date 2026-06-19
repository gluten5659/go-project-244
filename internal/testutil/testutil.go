package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const configFileMode = 0o600

func WriteTempFile(tb testing.TB, content string) string {
	tb.Helper()

	path := filepath.Join(tb.TempDir(), "config.json")

	require.NoError(tb, os.WriteFile(path, []byte(content), configFileMode))

	return path
}
