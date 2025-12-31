package jsonlogic

import (
	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func _and(values, data any) any {
	args := values.([]any)
	var last any
	for _, value := range args {
		parsed := parseValues(value, data)
		last = parsed
		if !typing.IsTrue(parsed) {
			return parsed
		}
	}

	return last
}

func _or(values, data any) any {
	args := values.([]any)
	var last any
	for _, value := range args {
		parsed := parseValues(value, data)
		last = parsed
		if typing.IsTrue(parsed) {
			return parsed
		}
	}

	return last
}

func evaluateClause(clause any, data any) any {
	parsed := parseValues(clause, data)

	if typing.IsMap(parsed) {
		return apply(parsed, data)
	}

	return parsed
}

func conditional(values, data any) any {
	values = values.([]any)

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
	// If the slice is not empty, there is an argument to negate
	if typing.IsSlice(values) && len(values.([]any)) > 0 {
		return !typing.IsTrue(values.([]any)[0])
	}
	return !typing.IsTrue(values)
}
