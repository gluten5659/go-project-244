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

func TestCommandRendersDiff(t *testing.T) {
	t.Parallel()

	output := strings.Builder{}
	command := newTestCommand(t, &output)

	arguments := []string{commandName, "testdata/file1.json", "testdata/file2.json"}

	err := command.Run(t.Context(), arguments)

	require.NoError(t, err)

	expected := "{\n" +
		"  - follow: false\n" +
		"    host: hexlet.io\n" +
		"  - proxy: 123.234.53.22\n" +
		"  - timeout: 50\n" +
		"  + timeout: 20\n" +
		"  + verbose: true\n" +
		"}\n"

	assert.Equal(t, expected, output.String())
}

func TestCommandRendersNestedDiff(t *testing.T) {
	t.Parallel()

	output := strings.Builder{}
	command := newTestCommand(t, &output)

	arguments := []string{commandName, "testdata/nested1.json", "testdata/nested2.json"}

	err := command.Run(t.Context(), arguments)

	require.NoError(t, err)

	expected := "{\n" +
		"    common: {\n" +
		"      + follow: false\n" +
		"        setting1: Value 1\n" +
		"      - setting2: 200\n" +
		"      - setting3: true\n" +
		"      + setting3: null\n" +
		"      + setting4: blah blah\n" +
		"      + setting5: {\n" +
		"            key5: value5\n" +
		"        }\n" +
		"        setting6: {\n" +
		"            doge: {\n" +
		"              - wow: \n" +
		"              + wow: so much\n" +
		"            }\n" +
		"            key: value\n" +
		"          + ops: vops\n" +
		"        }\n" +
		"    }\n" +
		"    group1: {\n" +
		"      - baz: bas\n" +
		"      + baz: bars\n" +
		"        foo: bar\n" +
		"      - nest: {\n" +
		"            key: value\n" +
		"        }\n" +
		"      + nest: str\n" +
		"    }\n" +
		"  - group2: {\n" +
		"        abc: 12345\n" +
		"        deep: {\n" +
		"            id: 45\n" +
		"        }\n" +
		"    }\n" +
		"  + group3: {\n" +
		"        deep: {\n" +
		"            id: {\n" +
		"                number: 45\n" +
		"            }\n" +
		"        }\n" +
		"        fee: 100500\n" +
		"    }\n" +
		"}\n"

	assert.Equal(t, expected, output.String())
}

func TestCommandExitCodes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		buildArguments   func(t testing.TB) []string
		expectedExitCode int
	}{
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
