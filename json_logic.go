package json_logic

import (
	"errors"
	//"log"
	"reflect"
)

func is(obj interface{}, kind reflect.Kind) bool {
	return reflect.TypeOf(obj).Kind() == kind
}

func isBool(obj interface{}) bool {
	return is(obj, reflect.Bool)
}

func isString(obj interface{}) bool {
	return is(obj, reflect.String)
}

func isInt(obj interface{}) bool {
	return is(obj, reflect.Int)
}

func isPrimitive(obj interface{}) bool {
	return isBool(obj) ||
		isString(obj) ||
		isInt(obj)
}

func isMap(obj interface{}) bool {
	return is(obj, reflect.Map)
}

func isArray(obj interface{}) bool {
	return is(obj, reflect.Array)
}

func isSlice(obj interface{}) bool {
	return is(obj, reflect.Slice)
}

func equals(a, b interface{}) bool {
	switch v := a.(type) {
	case float64:
		w := b.(float64)
		return v == w
	case string:
		w := b.(string)
		return v == w
	}

	return false
}

func operation(operator string, values, data interface{}) interface{} {
	if operator == "var" {
		return getVar(values, data)
	}

	parsed := values.([]interface{})

	if operator == "and" {
		return interface{}(parsed[0].(bool) && parsed[1].(bool))
	}

	equals(parsed[0], parsed[1])

	return interface{}(reflect.DeepEqual(parsed[0], parsed[1]))
}

func getVar(value, data interface{}) interface{} {
	if data == nil {
		return value
	}

	var parsed string

	if isSlice(value) {
		parsed = value.([]interface{})[0].(string)
	} else {
		parsed = value.(string)
	}

	if !isMap(data) {
		return nil
	}

	parsed_value := data.(map[string]interface{})[parsed]
	switch v := parsed_value.(type) {
	case int:
		return interface{}(float64(v))
	default:
		return v
	}
}

func parseValues(values, data interface{}) interface{} {
	if isPrimitive(values) {
		return values
	}

	parsed := make([]interface{}, 0)

	for _, value := range values.([]interface{}) {
		if isMap(value) {
			parsed = append(parsed, apply(value, data))
		} else {
			parsed = append(parsed, value)
		}
	}

	return parsed
}

func apply(rules, data interface{}) interface{} {
	for operator, values := range rules.(map[string]interface{}) {
		return operation(operator, parseValues(values, data), data)
	}

	return false
}

func GenericApply(rules, data interface{}) (interface{}, error) {
	if isBool(rules) {
		return rules, nil
	}

	if !isMap(rules) {
		return false, errors.New("The root element needs to be an object")
	}

	return apply(rules, data), nil
}

func BoolApply(rules, data interface{}) (bool, error) {
	value, err := GenericApply(rules, data)
	return value.(bool), err
}

func IntApply(rules, data interface{}) (float64, error) {
	value, err := GenericApply(rules, data)
	return value.(float64), err
}
