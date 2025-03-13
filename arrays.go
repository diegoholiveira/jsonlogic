package jsonlogic

import (
	"fmt"
	"strings"

	"github.com/qoala-platform/jsonlogic/v3/internal/typing"
)

type ErrReduceDataType struct {
	dataType string
}

func (e ErrReduceDataType) Error() string {
	return fmt.Sprintf("The type \"%s\" is not supported", e.dataType)
}

func filter(values, data any) any {
	parsed := values.([]any)

	var subject any

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	result := make([]any, 0)

	if subject == nil {
		return result
	}

	logic := solveVars(parsed[1], data)

	for _, value := range subject.([]any) {
		v := parseValues(logic, value)

		if typing.IsTrue(v) {
			result = append(result, value)
		}
	}

	return result
}

func _map(values, data any) any {
	parsed := values.([]any)

	var subject any

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	result := make([]any, 0)

	if subject == nil {
		return result
	}

	logic := solveVars(parsed[1], data)

	for _, value := range subject.([]any) {
		v := parseValues(logic, value)

		if typing.IsTrue(v) || typing.IsNumber(v) || typing.IsBool(v) {
			result = append(result, v)
		}
	}

	return result
}

func reduce(values, data any) any {
	parsed := values.([]any)

	var subject any

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if subject == nil {
		return float64(0)
	}

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

func _in(value any, values any) bool {
	if value == nil || values == nil {
		return false
	}

	if typing.IsString(values) {
		return strings.Contains(values.(string), value.(string))
	}

	if !typing.IsSlice(values) {
		return false
	}

	for _, element := range values.([]any) {
		if typing.IsSlice(element) {
			if _inRange(value, element) {
				return true
			}

			continue
		}

		if typing.IsNumber(value) {
			if typing.ToNumber(element) == value {
				return true
			}

			continue
		}

		if element == value {
			return true
		}
	}

	return false
}

func merge(values any, level int8) any {
	result := make([]any, 0)

	if typing.IsPrimitive(values) || level > 1 {
		return append(result, values)
	}

	if typing.IsSlice(values) {
		for _, value := range values.([]any) {
			_values := merge(value, level+1).([]any)

			result = append(result, _values...)
		}
	}

	return result
}

func missing(values, data any) any {
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

	var subject any

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

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

	var subject any

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

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

	var subject any

	if typing.IsMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if typing.IsSlice(parsed[0]) {
		subject = parsed[0]
	}

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
