package jsonlogic

import (
	"reflect"
	"strconv"
)

func is(obj any, kind reflect.Kind) bool {
	return obj != nil && reflect.TypeOf(obj).Kind() == kind
}

func isBool(obj any) bool {
	return is(obj, reflect.Bool)
}

func isString(obj any) bool {
	return is(obj, reflect.String)
}

func isNumber(obj any) bool {
	switch obj.(type) {
	case int, float64:
		return true
	default:
		return false
	}
}

func isPrimitive(obj any) bool {
	return isBool(obj) || isString(obj) || isNumber(obj)
}

func isMap(obj any) bool {
	return is(obj, reflect.Map)
}

func isSlice(obj any) bool {
	return is(obj, reflect.Slice)
}

func isEmptySlice(obj any) bool {
	if !isSlice(obj) {
		return false
	}

	for _, v := range obj.([]any) {
		if isTrue(v) {
			return false
		}
	}

	return true
}

func isTrue(obj any) bool {
	if isBool(obj) {
		return obj.(bool)
	}

	if isNumber(obj) {
		n := toNumber(obj)
		return n != 0
	}

	if isString(obj) || isSlice(obj) || isMap(obj) {
		length := reflect.ValueOf(obj).Len()
		return length > 0
	}

	return false
}

func toSliceOfNumbers(values any) []float64 {
	_values := values.([]any)

	numbers := make([]float64, len(_values))
	for i, n := range _values {
		numbers[i] = toNumber(n)
	}
	return numbers
}

func toNumber(value any) float64 {
	if isString(value) {
		w, _ := strconv.ParseFloat(value.(string), 64)

		return w
	}

	switch value := value.(type) {
	case int:
		return float64(value)
	default:
		return value.(float64)
	}
}

func toString(value any) string {
	if isNumber(value) {
		switch value := value.(type) {
		case int:
			return strconv.FormatInt(int64(value), 10)
		default:
			return strconv.FormatFloat(value.(float64), 'f', -1, 64)
		}
	}

	if value == nil {
		return ""
	}

	return value.(string)
}
