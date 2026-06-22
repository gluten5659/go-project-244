package parser

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var ErrParse = errors.New("parse config")

func Parse(fileType string, content []byte) (map[string]any, error) {
	var (
		config map[string]any
		err    error
	)

	switch fileType {
	case "yaml", "yml":
		err = yaml.Unmarshal(content, &config)
	case "json":
		err = json.Unmarshal(content, &config)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParse, err)
	}

	return config, nil
}
