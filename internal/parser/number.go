package parser

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Number struct {
	display string
}

func IntNumber(value int64) Number {
	return Number{display: strconv.FormatInt(value, 10)}
}

func FloatNumber(value float64) Number {
	return Number{display: formatFloat(value)}
}

func UintNumber(value uint64) Number {
	return Number{display: strconv.FormatUint(value, 10)}
}

func (number Number) String() string {
	return number.display
}

func (number Number) MarshalJSON() ([]byte, error) {
	return []byte(number.display), nil
}

func normalizeScalar(value any) (any, error) {
	switch typed := value.(type) {
	case json.Number:
		return numberFromToken(typed)
	case int:
		return IntNumber(int64(typed)), nil
	case int64:
		return IntNumber(typed), nil
	case uint64:
		return UintNumber(typed), nil
	case float64:
		return finiteFloat(typed)
	default:
		return value, nil
	}
}

func finiteFloat(value float64) (Number, error) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return Number{}, fmt.Errorf("%w: number %v is not finite", ErrParse, value)
	}

	return FloatNumber(value), nil
}

func numberFromToken(token json.Number) (Number, error) {
	text := string(token)

	if strings.ContainsAny(text, ".eE") {
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return Number{}, fmt.Errorf("%w: number %q: %w", ErrParse, text, err)
		}

		return finiteFloat(value)
	}

	signedValue, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		return IntNumber(signedValue), nil
	}

	unsignedValue, err := strconv.ParseUint(text, 10, 64)
	if err == nil {
		return UintNumber(unsignedValue), nil
	}

	return Number{}, fmt.Errorf("%w: integer %q is out of range", ErrParse, text)
}

func formatFloat(value float64) string {
	if value == 0 {
		return "0.0"
	}

	text := strconv.FormatFloat(value, 'g', -1, 64)
	if isPlainInteger(text) {
		text += ".0"
	}

	return text
}

func isPlainInteger(text string) bool {
	digits := strings.TrimPrefix(text, "-")
	if digits == "" {
		return false
	}

	for _, symbol := range digits {
		if symbol < '0' || symbol > '9' {
			return false
		}
	}

	return true
}
