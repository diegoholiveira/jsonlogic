package jsonlogic

import (
	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

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

func evaluateClause(clause any, data any) any {
	parsed := parseValues(clause, data)

	if typing.IsMap(parsed) {
		return apply(parsed, data)
	}

	return parsed
}

func conditional(values, data any) any {
	values = getValuesWithoutParsing(values, data)

	clauses := values.([]any)

	length := len(clauses)

	if length == 0 {
		return nil
	}

	// Evaluate each if/then pair
	for i := 0; i < length-1; i = i + 2 {
		condition := parseValues(clauses[i], data)

		// If the condition is true, evaluate and return the then clause
		if typing.IsTrue(condition) {
			return evaluateClause(clauses[i+1], data)
		}
	}

	// If no matches and there is an odd number of clauses, evaluate and return the else clause
	if length%2 == 1 {
		return evaluateClause(clauses[length-1], data)
	}

	return nil
}

func negative(values, data any) any {
	values = parseValues(values, data)
	if typing.IsSlice(values) {
		return !typing.IsTrue(values.([]any)[0])
	}
	return !typing.IsTrue(values)
}
