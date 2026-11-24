package jsonlogic

import (
	"strconv"

	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
)

func toNumber(value any) float64 {
	return javascript.ToNumber(value)
}

func toString(value any) string {
	if n, ok := value.(float64); ok {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	if value == nil {
		return ""
	}
	return value.(string)
}

func isPrimitive(obj any) bool {
	switch obj.(type) {
	case bool, string, float64:
		return true
	}
	return false
}

func sameType(a, b any) bool {
	switch a.(type) {
	case bool:
		_, ok := b.(bool)
		return ok
	case float64:
		_, ok := b.(float64)
		return ok
	case string:
		_, ok := b.(string)
		return ok
	case map[string]any:
		_, ok := b.(map[string]any)
		return ok
	case []any:
		_, ok := b.([]any)
		return ok
	}
	return false
}

// normalizeNumbers converts Go numeric values to float64, the type encoding/json
// uses for numbers, including inside []any and map[string]any. Any other value is
// returned unchanged. Containers are copied only when they hold a value to convert.
func normalizeNumbers(value any) any {
	normalized, _ := convertNumbers(value)
	return normalized
}

// convertNumbers also reports whether anything was converted.
func convertNumbers(value any) (any, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case []any:
		var out []any
		for i, item := range v {
			if converted, ok := convertNumbers(item); ok {
				if out == nil {
					out = append(out, v...)
				}
				out[i] = converted
			}
		}
		if out != nil {
			return out, true
		}
	case map[string]any:
		var out map[string]any
		for key, item := range v {
			if converted, ok := convertNumbers(item); ok {
				if out == nil {
					out = make(map[string]any, len(v))
					for k, original := range v {
						out[k] = original
					}
				}
				out[key] = converted
			}
		}
		if out != nil {
			return out, true
		}
	}
	return value, false
}
