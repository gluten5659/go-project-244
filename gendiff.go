package code

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parser"
	"fmt"
)

func GenDiff(firstFilePath, secondFilePath, format string) (string, error) {
	firstFile, err := parser.ParseFile(firstFilePath)
	if err != nil {
		return "", fmt.Errorf("load first file: %w", err)
	}

	secondFile, err := parser.ParseFile(secondFilePath)
	if err != nil {
		return "", fmt.Errorf("load second file: %w", err)
	}

	formatter, err := formatters.New(format)
	if err != nil {
		return "", fmt.Errorf("format diff: %w", err)
	}

	formatted, err := formatter.Format(diff.Compare(firstFile, secondFile))
	if err != nil {
		return "", fmt.Errorf("format diff: %w", err)
	}

	return formatted, nil
}
