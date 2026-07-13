package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrRead  = errors.New("read file")
	ErrParse = errors.New("parse config")
)

const (
	typeJSON = "json"
	typeYAML = "yaml"
	typeYML  = "yml"
)

func FromFile(path string) (map[string]any, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%q: %w: %w", path, ErrRead, errors.Unwrap(err))
	}

	fileType := strings.TrimPrefix(filepath.Ext(path), ".")

	values, err := parse(fileType, content)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", path, err)
	}

	return values, nil
}

func parse(fileType string, content []byte) (map[string]any, error) {
	var (
		parsedContent map[string]any
		err           error
	)

	switch fileType {
	case typeYAML, typeYML:
		err = yaml.Unmarshal(content, &parsedContent)
	case typeJSON:
		err = json.Unmarshal(content, &parsedContent)
	default:
		return nil, fmt.Errorf("%w: unsupported file type %q", ErrParse, fileType)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParse, err)
	}

	return normalizeMap(parsedContent), nil
}

func normalizeMap(content map[string]any) map[string]any {
	for key, value := range content {
		content[key] = normalizeValue(value)
	}

	return content
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return normalizeMap(typed)
	case map[any]any:
		converted := make(map[string]any, len(typed))
		for key, nested := range typed {
			converted[fmt.Sprint(key)] = normalizeValue(nested)
		}

		return converted
	case []any:
		for index, nested := range typed {
			typed[index] = normalizeValue(nested)
		}

		return typed
	default:
		return value
	}
}
