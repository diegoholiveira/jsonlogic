package jsonlogic

import (
	"math"
	"reflect"
	"strconv"
)

// at simulate undefined in javascript
func at(values []interface{}, index int) interface{} {
	if index >= 0 && index < len(values) {
		return values[index]
	}
	return undefinedType{}
}

type undefinedType struct{}

func toNumberForLess(v interface{}) float64 {
	switch value := v.(type) {
	case nil:
		return 0
	case undefinedType:
		return math.NaN()
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(value).Convert(reflect.TypeOf(float64(0))).Float()
	case bool: // Boolean values true and false are converted to 1 and 0 respectively.
		if value {
			return 1
		} else {
			return 0
		}
	case string:
		n, err := strconv.ParseFloat(value, 64)
		switch err {
		case strconv.ErrRange, nil:
			return n
		default:
			return math.NaN()
		}
	default:
		return math.NaN()
	}
}

// less reference javascript implementation
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#description
func less(a, b interface{}) bool {
	// If both values are strings, they are compared as strings,
	// based on the values of the Unicode code points they contain.
	if isString(a) && isString(b) {
		return toString(b) > toString(a)
	}

	// Otherwise the values are compared as numeric values.
	return toNumberForLess(b) > toNumberForLess(a)
}

func hardEquals(a, b interface{}) bool {
	ra := reflect.ValueOf(a).Kind()
	rb := reflect.ValueOf(b).Kind()

	if ra != rb {
		return false
	}

	return equals(a, b)
}

func equals(a, b interface{}) bool {
	// comparison to a nil value is falsy
	if a == nil || b == nil {
		// if a and b is nil, return true, else return falsy
		return a == b
	}

	if isNumber(a) {
		return toNumber(a) == toNumber(b)
	}

	if isBool(a) {
		if !isBool(b) {
			return false
		}
		return isTrue(a) == isTrue(b)
	}

	return toString(a) == toString(b)
}
