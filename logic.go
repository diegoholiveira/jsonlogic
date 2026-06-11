package jsonlogic

import (
	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
)

func _and(values, data any) any {
	s := values.([]any)
	if len(s) == 0 {
		return nil
	}
	var last any
	for _, value := range s {
		last = parseValues(value, data)
		if !javascript.IsTrue(last) {
			return last
		}
	}
	return last
}

func _or(values, data any) any {
	s := values.([]any)
	if len(s) == 0 {
		return nil
	}
	var last any
	for _, value := range s {
		last = parseValues(value, data)
		if javascript.IsTrue(last) {
			return last
		}
	}
	return last
}

func evaluateClause(clause any, data any) any {
	parsed := parseValues(clause, data)

	if m, ok := parsed.(map[string]any); ok {
		return apply(m, data)
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
		if javascript.IsTrue(condition) {
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
	if s, ok := values.([]any); ok && len(s) > 0 {
		return !javascript.IsTrue(s[0])
	}
	return !javascript.IsTrue(values)
}
