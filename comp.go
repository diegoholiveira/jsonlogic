package jsonlogic

import (
	"reflect"

	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

// less reference javascript implementation
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#description
func less(a, b any) bool {
	// If both values are strings, they are compared as strings,
	// based on the values of the Unicode code points they contain.
	if typing.IsString(a) && typing.IsString(b) {
		return typing.ToString(b) > typing.ToString(a)
	}

	// Otherwise the values are compared as numeric values.
	return javascript.ToNumber(b) > javascript.ToNumber(a)
}

func hardEquals(a, b any) bool {
	ra := reflect.ValueOf(a).Kind()
	rb := reflect.ValueOf(b).Kind()

	if ra != rb {
		return false
	}

	return equals(a, b)
}

func equals(a, b any) bool {
	// comparison to a nil value is falsy
	if a == nil || b == nil {
		// if a and b is nil, return true, else return falsy
		return a == b
	}

	if typing.IsString(a) && typing.IsString(b) {
		return a == b
	}

	return javascript.ToNumber(a) == javascript.ToNumber(b)
}

func between(operator string, values []any, data any) any {
	a := parseValues(values[0], data)
	b := parseValues(values[1], data)
	c := parseValues(values[2], data)

	if operator == "<" {
		return less(a, b) && less(b, c)
	}

	if operator == "<=" {
		return (less(a, b) || equals(a, b)) && (less(b, c) || equals(b, c))
	}

	if operator == ">=" {
		return (less(c, b) || equals(c, b)) && (less(b, a) || equals(b, a))
	}

	return less(c, b) && less(b, a)
}

func _inRange(value any, values any) bool {
	v := values.([]any)

	i := v[0]
	j := v[1]

	if typing.IsNumber(value) {
		return typing.ToNumber(value) >= typing.ToNumber(i) && typing.ToNumber(j) >= typing.ToNumber(value)
	}

	return typing.ToString(value) >= typing.ToString(i) && typing.ToString(j) >= typing.ToString(value)
}
