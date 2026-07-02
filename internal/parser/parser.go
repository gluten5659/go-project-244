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

type unmarshaler func([]byte, any) error

func unmarshalers() map[string]unmarshaler {
	return map[string]unmarshaler{
		TypeJSON: json.Unmarshal,
		TypeYAML: yaml.Unmarshal,
		TypeYML:  yaml.Unmarshal,
	}
}

func Parse(fileType string, content []byte) (map[string]any, error) {
	unmarshal, isSupported := unmarshalers()[fileType]
	if !isSupported {
		return nil, fmt.Errorf("%w: unsupported file type %q", ErrParse, fileType)
	}

	var parsedContent map[string]any

	err := unmarshal(content, &parsedContent)
	if err != nil {
		if isSequence(unmarshal, content) {
			return map[string]any{}, nil
		}

		return nil, fmt.Errorf("%w: %w", ErrParse, err)
	}

	return parsedContent, nil
}

func isSequence(unmarshal unmarshaler, content []byte) bool {
	var sequence []any

	return unmarshal(content, &sequence) == nil
}
