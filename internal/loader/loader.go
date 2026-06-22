package loader

import (
	"code/internal/files"
	"code/internal/parser"
	"fmt"
)

func FromFile(path string) (map[string]any, error) {
	fileType, content, err := files.Read(path)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	values, err := parser.Parse(fileType, content)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	return values, nil
}
