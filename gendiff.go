package code

import (
	"code/internal/compare"
	"code/internal/formatters"
	"code/internal/loader"
	"fmt"
)

func GenDiff(firstFilePath, secondFilePath, format string) (string, error) {
	firstFile, err := loader.FromFile(firstFilePath)
	if err != nil {
		return "", fmt.Errorf("load first file: %w", err)
	}

	secondFile, err := loader.FromFile(secondFilePath)
	if err != nil {
		return "", fmt.Errorf("load second file: %w", err)
	}

	formatted, err := formatters.Format(compare.Compare(firstFile, secondFile), format)
	if err != nil {
		return "", fmt.Errorf("format diff: %w", err)
	}

	return formatted, nil
}
