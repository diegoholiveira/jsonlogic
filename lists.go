package jsonlogic

import (
	"fmt"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

type ErrReduceDataType struct {
	dataType string
}

func (e ErrReduceDataType) Error() string {
	return fmt.Sprintf("The type \"%s\" is not supported", e.dataType)
}

func extractSubject(parsed []any, data any) any {
	var subject any

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	return subject
}

func filter(values, data any) any {
	parsed := values.([]any)

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

		if typing.IsTrue(v) {
			result = append(result, value)
		}
	}

	return result
}

func _map(values, data any) any {
	parsed := values.([]any)

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

	var (
		accumulator any
		valueType   string
	)

	{
		initialValue := parsed[2]
		if typing.IsMap(initialValue) {
			initialValue = apply(initialValue, data)
		}

		if typing.IsBool(initialValue) {
			accumulator = typing.IsTrue(initialValue)
			valueType = "bool"
		} else if typing.IsNumber(initialValue) {
			accumulator = typing.ToNumber(initialValue)
			valueType = "number"
		} else if typing.IsString(initialValue) {
			accumulator = typing.ToString(initialValue)
			valueType = "string"
		} else {
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
			context["accumulator"] = typing.IsTrue(v)
		case "number":
			context["accumulator"] = typing.ToNumber(v)
		case "string":
			context["accumulator"] = typing.ToString(v)
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

	if typing.IsString(b) {
		return strings.Contains(b.(string), a.(string))
	}

	if !typing.IsSlice(b) {
		return false
	}

	for _, element := range b.([]any) {
		if typing.IsSlice(element) {
			if _inRange(a, element.([]any)) {
				return true
			}

			continue
		}

		if typing.IsNumber(a) {
			if typing.ToNumber(element) == a {
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
	if typing.IsPrimitive(values) {
		return []any{values}
	}

	inputSlice := values.([]any)
	sliceLen := len(inputSlice)
	if sliceLen == 0 {
		return inputSlice
	}

	totalCapacity := 0
	for _, value := range inputSlice {
		if typing.IsSlice(value) {
			totalCapacity += len(value.([]any))
		} else {
			totalCapacity++
		}
	}

	result := make([]any, 0, totalCapacity)

	for _, value := range inputSlice {
		if !typing.IsSlice(value) {
			result = append(result, value)
			continue
		}

		result = append(result, value.([]any)...)
	}

	return result
}

func missing(values, data any) any {
	values = parseValues(values, data)
	if typing.IsString(values) {
		values = []any{values}
	}

	missing := make([]any, 0)

	for _, _var := range values.([]any) {
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
	number := int(typing.ToNumber(parsed[0]))
	vars := parsed[1]

	missing := make([]any, 0)
	found := make([]any, 0)

	for _, _var := range vars.([]any) {
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
	if !typing.IsTrue(subject) {
		return false
	}

	for _, value := range subject.([]any) {
		conditions := solveVars(parsed[1], value)
		v := apply(conditions, value)

		if !typing.IsTrue(v) {
			return false
		}
	}

	return true
}

func none(values, data any) any {
	parsed := values.([]any)

	subject := extractSubject(parsed, data)

	if !typing.IsTrue(subject) {
		return true
	}

	conditions := solveVars(parsed[1], data)

	for _, value := range subject.([]any) {
		v := apply(conditions, value)

		if typing.IsTrue(v) {
			return false
		}
	}

	return true
}

func some(values, data any) any {
	parsed := values.([]any)
	subject := extractSubject(parsed, data)

	if !typing.IsTrue(subject) {
		return false
	}

	for _, value := range subject.([]any) {
		v := apply(
			solveVars(
				solveVars(parsed[1], data),
				value,
			),
			value,
		)

		if typing.IsTrue(v) {
			return true
		}
	}

	return false
}

func _inRange(value any, values []any) bool {
	i := values[0]
	j := values[1]

	return typing.ToNumber(value) >= typing.ToNumber(i) && typing.ToNumber(j) >= typing.ToNumber(value)
}
