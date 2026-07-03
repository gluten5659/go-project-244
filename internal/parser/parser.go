package parser

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var ErrParse = errors.New("parse config")

const (
	TypeJSON = "json"
	TypeYAML = "yaml"
	TypeYML  = "yml"
)

func Parse(fileType string, content []byte) (map[string]any, error) {
	var (
		parsedContent map[string]any
		err           error
	)

	switch fileType {
	case TypeYAML, TypeYML:
		err = yaml.Unmarshal(content, &parsedContent)
	case TypeJSON:
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
