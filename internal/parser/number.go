package parser

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	minInt64AsFloat = float64(math.MinInt64)
	maxInt64AsFloat = -minInt64AsFloat
)

func normalizeScalar(value any) (any, error) {
	switch typed := value.(type) {
	case json.Number:
		return parseNumberToken(typed)
	case int:
		return int64(typed), nil
	case int64:
		return typed, nil
	case uint64:
		return typed, nil
	case float64:
		return normalizeFloat(typed)
	default:
		return value, nil
	}
}

func parseNumberToken(token json.Number) (any, error) {
	text := string(token)

	if strings.ContainsAny(text, ".eE") {
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: number %q: %w", ErrParse, text, err)
		}

		return normalizeFloat(value)
	}

	signedValue, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		return signedValue, nil
	}

	unsignedValue, err := strconv.ParseUint(text, 10, 64)
	if err == nil {
		return unsignedValue, nil
	}

	return nil, fmt.Errorf("%w: integer %q is out of range", ErrParse, text)
}

func normalizeFloat(value float64) (any, error) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return nil, fmt.Errorf("%w: number %v is not finite", ErrParse, value)
	}

	if fitsInt64(value) {
		return int64(value), nil
	}

	return value, nil
}

func fitsInt64(value float64) bool {
	return value == math.Trunc(value) && value >= minInt64AsFloat && value < maxInt64AsFloat
}
