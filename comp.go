package jsonlogic

import "github.com/diegoholiveira/jsonlogic/v3/internal/javascript"

// parsePair evaluates the operands of a list with exactly two items without
// allocating a slice for them. ok is false for any other shape.
func parsePair(values, data any) (a, b any, ok bool) {
	pair, isSlice := values.([]any)
	if !isSlice || len(pair) != 2 {
		return nil, nil, false
	}
	return parseValues(pair[0], data), parseValues(pair[1], data), true
}

func hardEquals(values, data any) any {
	a, b, ok := parsePair(values, data)
	if !ok {
		parsed, isSlice := parseValues(values, data).([]any)
		if !isSlice || len(parsed) < 2 {
			return false
		}
		a, b = parsed[0], parsed[1]
	}

	if a == nil || b == nil {
		return a == b
	}

	if !sameType(a, b) {
		return false
	}

	return equals(a, b)
}

func isLessThan(values, data any) any {
	parsed := parseValues(values, data).([]any)
	if len(parsed) < 2 {
		return false
	}

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
	if len(parsed) < 2 {
		return false
	}

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
	if len(parsed) < 2 {
		return false
	}
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
	if len(parsed) < 2 {
		return false
	}

	a := parsed[0]
	b := parsed[1]

	if len(parsed) == 3 {
		c := parsed[2]

		return (less(c, b) || equals(c, b)) && (less(b, a) || equals(b, a))
	}

	return less(b, a) || equals(b, a)
}

func isEqual(values, data any) any {
	a, b, ok := parsePair(values, data)
	if !ok {
		parsed := parseValues(values, data).([]any)
		if len(parsed) < 2 {
			return false
		}
		a, b = parsed[0], parsed[1]
	}

	return equals(a, b)
}

// less reference javascript implementation
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#description
func less(a, b any) bool {
	// If both values are strings, they are compared as strings,
	// based on the values of the Unicode code points they contain.
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if aOk && bOk {
		return bStr > aStr
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

	_, aOk := a.(string)
	_, bOk := b.(string)
	if aOk && bOk {
		return a == b
	}

	return javascript.ToNumber(a) == javascript.ToNumber(b)
}
