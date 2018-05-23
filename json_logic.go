package json_logic

import (
	"errors"
	// "log"
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

func operation(operator string, parsed []interface{}) bool {
	return reflect.DeepEqual(parsed[0], parsed[1])
}

func apply(values interface{}) interface{} {
	return values
}

func parseValues(values interface{}) []interface{} {
	parsed := make([]interface{}, 0)
	for _, value := range values.([]interface{}) {
		parsed = append(parsed, apply(value))
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

	for operator, values := range rules.(map[string]interface{}) {
		parsed := parseValues(values)
		return operation(operator, parsed), nil
	}

	return false, nil
}
