package parser

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Number struct {
	isInt   bool
	display string
}

func IntNumber(value int64) Number {
	return Number{isInt: true, display: strconv.FormatInt(value, 10)}
}

func FloatNumber(value float64) Number {
	return Number{isInt: false, display: formatFloat(value)}
}

func UintNumber(value uint64) Number {
	return Number{isInt: true, display: strconv.FormatUint(value, 10)}
}

func (number Number) String() string {
	return number.display
}

func (number Number) MarshalJSON() ([]byte, error) {
	return []byte(number.display), nil
}

func normalizeScalar(value any) any {
	switch typed := value.(type) {
	case json.Number:
		return numberFromToken(typed)
	case int:
		return IntNumber(int64(typed))
	case int64:
		return IntNumber(typed)
	case uint64:
		return UintNumber(typed)
	case float64:
		return FloatNumber(typed)
	default:
		return value
	}
}

func numberFromToken(token json.Number) Number {
	text := string(token)

	if strings.ContainsAny(text, ".eE") {
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return Number{isInt: false, display: text}
		}

		return FloatNumber(value)
	}

	signedValue, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		return IntNumber(signedValue)
	}

	unsignedValue, err := strconv.ParseUint(text, 10, 64)
	if err == nil {
		return UintNumber(unsignedValue)
	}

	return Number{isInt: true, display: text}
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
