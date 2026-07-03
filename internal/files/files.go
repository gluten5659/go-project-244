package files

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrRead = errors.New("read file")

func Read(path string) (string, []byte, error) {
	fileType := extension(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return fileType, nil, fmt.Errorf("%w: %w", ErrRead, errors.Unwrap(err))
	}

	return fileType, content, nil
}

func extension(path string) string {
	lastDotIndex := strings.LastIndex(path, `.`)
	if lastDotIndex == -1 {
		return ""
	}

	return path[lastDotIndex+1:]
}
