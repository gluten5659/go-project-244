package files

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrRead = errors.New("read file")

func Read(path string) (string, []byte, error) {
	fileType := getFileType(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return fileType, nil, fmt.Errorf("%w: %w", ErrRead, err)
	}

	return fileType, content, nil
}

func getFileType(path string) string {
	lastDotIndex := strings.LastIndex(path, `.`)
	if lastDotIndex == -1 {
		return ""
	}

	return path[lastDotIndex+1:]
}
