package jsonlogic

import "strconv"

func toNumber(value any) float64 {
	if s, ok := value.(string); ok {
		w, _ := strconv.ParseFloat(s, 64)
		return w
	}
	return value.(float64)
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
