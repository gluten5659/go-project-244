package cliapp_test

import (
	"code/internal/cliapp"
	"code/internal/testutil"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

const (
	commandName    = "gendiff"
	formatFlag     = "--format"
	nestedFirst    = "testdata/nested1.json"
	nestedSecond   = "testdata/nested2.json"
	exitUsage      = 64
	exitDataErr    = 65
	exitNoInput    = 66
	exitIOErr      = 74
	exitPermission = 77
	unreadableMode = 0o000
)

func newTestCommand(tb testing.TB, output *strings.Builder) *cli.Command {
	tb.Helper()

	command := cliapp.NewCommand()
	command.Writer = output
	command.ExitErrHandler = func(context.Context, *cli.Command, error) {}

	return command
}

func readGolden(tb testing.TB, path string) string {
	tb.Helper()

	content, err := os.ReadFile(path)
	require.NoError(tb, err)

	return string(content)
}

func TestCommandRendersGolden(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		arguments  []string
		goldenFile string
	}{
		{
			name:       "flat json inputs render in stylish",
			arguments:  []string{commandName, "testdata/file1.json", "testdata/file2.json"},
			goldenFile: "testdata/expected/flat.stylish",
		},
		{
			name:       "flat yaml inputs render the same stylish output",
			arguments:  []string{commandName, "testdata/file1.yaml", "testdata/file2.yaml"},
			goldenFile: "testdata/expected/flat.stylish",
		},
		{
			name:       "nested inputs render in stylish",
			arguments:  []string{commandName, nestedFirst, nestedSecond},
			goldenFile: "testdata/expected/nested.stylish",
		},
		{
			name:       "nested inputs render in plain",
			arguments:  []string{commandName, formatFlag, "plain", nestedFirst, nestedSecond},
			goldenFile: "testdata/expected/nested.plain",
		},
		{
			name:       "nested inputs render in json",
			arguments:  []string{commandName, formatFlag, "json", nestedFirst, nestedSecond},
			goldenFile: "testdata/expected/nested.json",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := strings.Builder{}
			command := newTestCommand(t, &output)

			err := command.Run(t.Context(), testCase.arguments)

			require.NoError(t, err)
			assert.Equal(t, readGolden(t, testCase.goldenFile), output.String())
		})
	}
}

func TestCommandExitCodes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		buildArguments   func(t testing.TB) []string
		expectedExitCode int
	}{
		{
			name: "unsupported format yields a usage error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				readable := testutil.WriteTempFile(tb, `{}`)

				return []string{commandName, formatFlag, "bogus", readable, readable}
			},
			expectedExitCode: exitUsage,
		},
		{
			name: "no file paths yield a usage error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				return []string{commandName}
			},
			expectedExitCode: exitUsage,
		},
		{
			name: "a single file path yields a usage error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				readable := testutil.WriteTempFile(tb, `{}`)

				return []string{commandName, readable}
			},
			expectedExitCode: exitUsage,
		},
		{
			name: "extra file paths yield a usage error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				readable := testutil.WriteTempFile(tb, `{}`)

				return []string{commandName, readable, readable, readable}
			},
			expectedExitCode: exitUsage,
		},
		{
			name: "malformed json yields a data error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				malformed := testutil.WriteTempFile(tb, `{`)
				readable := testutil.WriteTempFile(tb, `{}`)

				return []string{commandName, malformed, readable}
			},
			expectedExitCode: exitDataErr,
		},
		{
			name: "missing file yields a no-input error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				missing := filepath.Join(tb.TempDir(), "missing.json")
				readable := testutil.WriteTempFile(tb, `{}`)

				return []string{commandName, missing, readable}
			},
			expectedExitCode: exitNoInput,
		},
		{
			name: "directory instead of a file yields an io error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				readable := testutil.WriteTempFile(tb, `{}`)

				return []string{commandName, tb.TempDir(), readable}
			},
			expectedExitCode: exitIOErr,
		},
		{
			name: "unknown flag yields a usage error",
			buildArguments: func(tb testing.TB) []string {
				tb.Helper()

				return []string{commandName, "--unknown"}
			},
			expectedExitCode: exitUsage,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := strings.Builder{}
			command := newTestCommand(t, &output)

			err := command.Run(t.Context(), testCase.buildArguments(t))

			var exitError cli.ExitCoder

			require.ErrorAs(t, err, &exitError)
			assert.Equal(t, testCase.expectedExitCode, exitError.ExitCode())
		})
	}
}

func TestCommandReportsPermissionError(t *testing.T) {
	t.Parallel()

	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to the root user")
	}

	unreadable := testutil.WriteTempFile(t, `{}`)
	require.NoError(t, os.Chmod(unreadable, unreadableMode))

	readable := testutil.WriteTempFile(t, `{}`)

	output := strings.Builder{}
	command := newTestCommand(t, &output)

	err := command.Run(t.Context(), []string{commandName, unreadable, readable})

	var exitError cli.ExitCoder

	require.ErrorAs(t, err, &exitError)
	assert.Equal(t, exitPermission, exitError.ExitCode())
}
