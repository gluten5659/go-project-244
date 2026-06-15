package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrParse = errors.New("parse config")

func Parse(content []byte) (map[string]any, error) {
	var config map[string]any

	err := json.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParse, err)
	}

	return config, nil
}
