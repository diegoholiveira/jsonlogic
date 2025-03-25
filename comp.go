package jsonlogic

import (
	"reflect"

	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func hardEquals(values, data any) any {
	values = parseValues(values, data)
	if !typing.IsSlice(values) {
		return false
	}

	parsed := values.([]any)
	if len(parsed) < 2 {
		return false
	}

	a, b := parsed[0], parsed[1]

	if a == nil || b == nil {
		return a == b
	}

	ra := reflect.ValueOf(a).Kind()
	rb := reflect.ValueOf(b).Kind()

	if ra != rb {
		return false
	}

	return equals(a, b)
}

func isLessThan(values, data any) any {
	parsed := parseValues(values, data).([]any)

	a := parsed[0]
	b := parsed[1]

	if len(parsed) == 3 {
		c := parsed[2]

		return less(a, b) && less(b, c)
	}

	return less(a, b)
}

func isLessOrEqualThan(values, data any) any {
	parsed := parseValues(values, data).([]any)

	a := parsed[0]
	b := parsed[1]

	if len(parsed) == 3 {
		c := parsed[2]

		return (less(a, b) || equals(a, b)) && (less(b, c) || equals(b, c))
	}

	return less(a, b) || equals(a, b)
}

func isGreaterThan(values, data any) any {
	parsed := parseValues(values, data).([]any)
	a := parsed[0]
	b := parsed[1]

	if len(parsed) == 3 {
		c := parsed[2]

		return less(c, b) && less(b, a)
	}

	return less(b, a)
}

func isGreaterOrEqualThan(values, data any) any {
	parsed := parseValues(values, data).([]any)

	a := parsed[0]
	b := parsed[1]

	if len(parsed) == 3 {
		c := parsed[2]

		return (less(c, b) || equals(c, b)) && (less(b, a) || equals(b, a))
	}

	return less(b, a) || equals(b, a)
}

func isEqual(values, data any) any {
	parsed := parseValues(values, data).([]any)

	a := parsed[0]
	b := parsed[1]

	return equals(a, b)
}

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
