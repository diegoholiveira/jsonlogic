package jsonlogic

import "github.com/qoala-platform/jsonlogic/v3/internal/typing"

func _and(values, data any) any {
	values = getValuesWithoutParsing(values, data)

	var v float64

	isBoolExpression := true

	for _, value := range values.([]any) {
		value = parseValues(value, data)
		if typing.IsSlice(value) {
			return value
		}

		if typing.IsBool(value) && !value.(bool) {
			return false
		}

		if typing.IsString(value) && typing.ToString(value) == "" {
			return value
		}

		if !typing.IsNumber(value) {
			continue
		}

		isBoolExpression = false

		_value := typing.ToNumber(value)

		if _value > v {
			v = _value
		}
	}

	if isBoolExpression {
		return true
	}

	return v
}

func _or(values, data any) any {
	values = getValuesWithoutParsing(values, data)

	for _, value := range values.([]any) {
		parsed := parseValues(value, data)
		if typing.IsTrue(parsed) {
			return parsed
		}
	}

	return false
}

func conditional(values, data any) any {
	if typing.IsPrimitive(values) {
		return values
	}

	parsed := values.([]any)

	length := len(parsed)

	if length == 0 {
		return nil
	}

	for i := 0; i < length-1; i = i + 2 {
		v := parsed[i]
		if typing.IsMap(v) {
			v = getVar(parsed[i], data)
		}

		if typing.IsTrue(v) {
			return parseValues(parsed[i+1], data)
		}
	}

	if length%2 == 1 {
		return parsed[length-1]
	}

	return nil
}

func unary(operator string, value any) any {
	if operator == "+" || operator == "*" || operator == "/" {
		return typing.ToNumber(value)
	}

	if operator == "-" {
		return -1 * typing.ToNumber(value)
	}

	if operator == "!!" {
		return !unary("!", value).(bool)
	}

	if operator == "abs" {
		return abs(value)
	}

	if operator == "floor" {
		return floor(typing.ToNumber(value))
	}

	if operator == "ceil" {
		return ceil(typing.ToNumber(value))
	}

	b := typing.IsTrue(value)

	if operator == "!" {
		return !b
	}

	return b
}
