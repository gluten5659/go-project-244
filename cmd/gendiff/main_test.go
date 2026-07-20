package main_test

import (
	"code/internal/testutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const exitNoInput = 66

func buildBinary(t *testing.T) string {
	t.Helper()

	binaryPath := filepath.Join(t.TempDir(), "gendiff")

	output, err := exec.CommandContext(t.Context(), "go", "build", "-o", binaryPath, ".").
		CombinedOutput()
	require.NoError(t, err, string(output))

	return binaryPath
}

func TestBinary(t *testing.T) {
	t.Parallel()

	binaryPath := buildBinary(t)

	t.Run("renders a diff of mixed json and yaml inputs", func(t *testing.T) {
		t.Parallel()

		firstPath := testutil.WriteTempFileNamed(t, "first.json", `{"timeout": 50}`)
		secondPath := testutil.WriteTempFileNamed(t, "second.yaml", "timeout: 20\n")

		output, err := exec.CommandContext(t.Context(), binaryPath, firstPath, secondPath).
			CombinedOutput()

		require.NoError(t, err)
		assert.Equal(t, "{\n  - timeout: 50\n  + timeout: 20\n}\n", string(output))
	})

	t.Run("reports a missing file with a sysexits exit code", func(t *testing.T) {
		t.Parallel()

		readablePath := testutil.WriteTempFileNamed(t, "second.json", `{}`)
		missingPath := filepath.Join(t.TempDir(), "missing.json")

		output, err := exec.CommandContext(t.Context(), binaryPath, missingPath, readablePath).
			CombinedOutput()

		var exitError *exec.ExitError

		require.ErrorAs(t, err, &exitError)
		assert.Equal(t, exitNoInput, exitError.ExitCode())
		assert.Contains(t, string(output), "no such file or directory")
	})
}
