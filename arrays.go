package jsonlogic

import (
	"fmt"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
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

	logic := parsed[1]

	for _, value := range subject.([]any) {
		v := parseValues(logic, value)

		result = append(result, v)
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
