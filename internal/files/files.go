package files

import (
	"errors"
	"fmt"
	"os"
)

var ErrRead = errors.New("read file")

func Read(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRead, err)
	}

	return content, nil
}
