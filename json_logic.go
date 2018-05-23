package json_logic

import (
	"errors"
	//"log"
	"reflect"
)

func isBool(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Bool
}

func isMap(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Map
}

func isArray(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Array
}

func isSlice(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Slice
}

func operation(operator string, parsed []interface{}) bool {
	if operator == "and" {
		return parsed[0].(bool) && parsed[1].(bool)
	}

	return reflect.DeepEqual(parsed[0], parsed[1])
}

func apply(rules, data interface{}) bool {
	for operator, values := range rules.(map[string]interface{}) {
		parsed := parseValues(values, data)
		return operation(operator, parsed.([]interface{}))
	}

	return false
}

func parseValues(values, data interface{}) interface{} {
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

func Apply(rules, data interface{}) (bool, error) {
	if isBool(rules) {
		return rules.(bool), nil
	}

	if !isMap(rules) {
		return false, errors.New("The root element needs to be an object")
	}

	return apply(rules, data), nil
}
