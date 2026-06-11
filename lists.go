package jsonlogic

import (
	"fmt"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
)

// ErrReduceDataType represents an error when an unsupported data type is used in reduce operations.
// It contains the data type name that caused the error.
type ErrReduceDataType struct {
	dataType string
}

func (e ErrReduceDataType) Error() string {
	return fmt.Sprintf("The type \"%s\" is not supported", e.dataType)
}

func extractSubject(parsed []any, data any) any {
	var subject any

	if s, ok := parsed[0].([]any); ok {
		subject = s
	} else if m, ok := parsed[0].(map[string]any); ok {
		subject = apply(m, data)
	}

	return subject
}

func filter(values, data any) any {
	parsed := values.([]any)
	if len(parsed) < 2 {
		return []any{}
	}

	subject := extractSubject(parsed, data)
	if subject == nil {
		return []any{}
	}

	subjectSlice := subject.([]any)
	subjectLen := len(subjectSlice)

	// Pre-allocate result with capacity that's reasonable for filtering
	// Assuming at least half might pass the filter (heuristic)
	result := make([]any, 0, subjectLen/2)

	logic := solveVars(parsed[1], data)

	for _, value := range subjectSlice {
		v := parseValues(logic, value)

		if javascript.IsTrue(v) {
			result = append(result, value)
		}
	}

	return result
}

func _map(values, data any) any {
	parsed := values.([]any)
	if len(parsed) < 2 {
		return []any{}
	}

	subject := extractSubject(parsed, data)
	if subject == nil {
		return []any{}
	}

	subjectSlice := subject.([]any)
	subjectLen := len(subjectSlice)

	result := make([]any, 0, subjectLen)

	logic := parsed[1]

	for _, value := range subjectSlice {
		v := parseValues(logic, value)
		result = append(result, v)
	}

	return result
}

func reduce(values, data any) any {
	parsed := values.([]any)
	if len(parsed) < 3 {
		return float64(0)
	}

	var (
		accumulator any
		valueType   string
	)

	{
		initialValue := parsed[2]
		if m, ok := initialValue.(map[string]any); ok {
			initialValue = apply(m, data)
		}

		switch v := initialValue.(type) {
		case bool:
			accumulator = javascript.IsTrue(v)
			valueType = "bool"
		case float64:
			accumulator = v
			valueType = "number"
		case string:
			accumulator = v
			valueType = "string"
		default:
			panic(ErrReduceDataType{
				dataType: fmt.Sprintf("%T", parsed[2]),
			})
		}
	}

	context := map[string]any{
		"current":     float64(0),
		"accumulator": accumulator,
		"valueType":   valueType,
	}

	subject := extractSubject(parsed, data)
	if subject == nil {
		return float64(0)
	}

	for _, value := range subject.([]any) {
		if value == nil {
			continue
		}

		context["current"] = value

		v := apply(parsed[1], context)

		switch context["valueType"] {
		case "bool":
			context["accumulator"] = javascript.IsTrue(v)
		case "number":
			context["accumulator"] = toNumber(v)
		case "string":
			context["accumulator"] = toString(v)
		}
	}

	return context["accumulator"]
}

func _in(values, data any) any {
	values = parseValues(values, data)

	parsed := values.([]any)

	a := parsed[0]
	var b any
	if len(parsed) > 1 {
		b = parsed[1]
	}

	if bs, ok := b.(string); ok {
		return strings.Contains(bs, a.(string))
	}

	bSlice, ok := b.([]any)
	if !ok {
		return false
	}

	for _, element := range bSlice {
		if es, ok := element.([]any); ok {
			if _inRange(a, es) {
				return true
			}

			continue
		}

		if _, ok := a.(float64); ok {
			if toNumber(element) == a {
				return true
			}

			continue
		}

		if element == a {
			return true
		}
	}

	return false
}

func merge(values, data any) any {
	values = parseValues(values, data)
	if isPrimitive(values) {
		return []any{values}
	}

	inputSlice := values.([]any)
	sliceLen := len(inputSlice)
	if sliceLen == 0 {
		return inputSlice
	}

	totalCapacity := 0
	for _, value := range inputSlice {
		if sv, ok := value.([]any); ok {
			totalCapacity += len(sv)
		} else {
			totalCapacity++
		}
	}

	result := make([]any, 0, totalCapacity)

	for _, value := range inputSlice {
		sv, ok := value.([]any)
		if !ok {
			result = append(result, value)
			continue
		}

		result = append(result, sv...)
	}

	return result
}

func missing(values, data any) any {
	values = parseValues(values, data)
	if _, ok := values.(string); ok {
		values = []any{values}
	}

	s := values.([]any)
	if len(s) == 0 {
		return []any{}
	}

	missing := make([]any, 0)

	for _, _var := range s {
		_value := getVar(_var, data)

		if _value == nil {
			missing = append(missing, _var)
		}
	}

	return missing
}

func missingSome(values, data any) any {
	values = parseValues(values, data)
	parsed := values.([]any)
	number := int(toNumber(parsed[0]))

	vars, ok := parsed[1].([]any)
	if !ok || len(vars) == 0 {
		return []any{}
	}

	missing := make([]any, 0)
	found := make([]any, 0)

	for _, _var := range vars {
		_value := getVar(_var, data)

		if _value == nil {
			missing = append(missing, _var)
		} else {
			found = append(found, _var)
		}
	}

	if number > len(found) {
		return missing
	}

	return make([]any, 0)
}

func all(values, data any) any {
	parsed := values.([]any)

	subject := extractSubject(parsed, data)
	if !javascript.IsTrue(subject) {
		return false
	}

	for _, value := range subject.([]any) {
		conditions := solveVars(parsed[1], value)
		v := apply(conditions, value)

		if !javascript.IsTrue(v) {
			return false
		}
	}

	return true
}

func none(values, data any) any {
	parsed := values.([]any)

	subject := extractSubject(parsed, data)

	if !javascript.IsTrue(subject) {
		return true
	}

	conditions := solveVars(parsed[1], data)

	for _, value := range subject.([]any) {
		v := apply(conditions, value)

		if javascript.IsTrue(v) {
			return false
		}
	}

	return true
}

func some(values, data any) any {
	parsed := values.([]any)
	subject := extractSubject(parsed, data)

	if !javascript.IsTrue(subject) {
		return false
	}

	logic := solveVars(parsed[1], data)
	for _, value := range subject.([]any) {
		conditions := solveVars(logic, value)
		v := apply(conditions, value)

		if javascript.IsTrue(v) {
			return true
		}
	}

	return false
}

func _inRange(value any, values []any) bool {
	i := values[0]
	j := values[1]

	return toNumber(value) >= toNumber(i) && toNumber(j) >= toNumber(value)
}
