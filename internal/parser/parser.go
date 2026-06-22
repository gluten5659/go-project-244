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
		err    error
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

	return parsedContent, nil
}
