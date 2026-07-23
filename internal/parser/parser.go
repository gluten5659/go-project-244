package parser

import (
	"bytes"
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

func ParseFile(path string) (map[string]any, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRead, err)
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
		decoder := json.NewDecoder(bytes.NewReader(content))
		decoder.UseNumber()
		err = decoder.Decode(&parsedContent)
	default:
		return nil, fmt.Errorf("%w: unsupported file type %q", ErrParse, fileType)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParse, err)
	}

	return normalizeMap(parsedContent)
}

func normalizeMap(content map[string]any) (map[string]any, error) {
	for key, value := range content {
		normalized, err := normalizeValue(value)
		if err != nil {
			return nil, err
		}

		content[key] = normalized
	}

	return content, nil
}

func normalizeValue(value any) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		return normalizeMap(typed)
	case map[any]any:
		return normalizeStringKeyedMap(typed)
	case []any:
		return normalizeSlice(typed)
	default:
		return normalizeScalar(value)
	}
}

func normalizeStringKeyedMap(content map[any]any) (map[string]any, error) {
	converted := make(map[string]any, len(content))
	for key, value := range content {
		normalized, err := normalizeValue(value)
		if err != nil {
			return nil, err
		}

		converted[fmt.Sprint(key)] = normalized
	}

	return converted, nil
}

func normalizeSlice(content []any) ([]any, error) {
	for index, value := range content {
		normalized, err := normalizeValue(value)
		if err != nil {
			return nil, err
		}

		content[index] = normalized
	}

	return content, nil
}
