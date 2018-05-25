package jsonlogic

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func is(obj interface{}, kind reflect.Kind) bool {
	return obj != nil && reflect.TypeOf(obj).Kind() == kind
}

func isBool(obj interface{}) bool {
	return is(obj, reflect.Bool)
}

func isString(obj interface{}) bool {
	return is(obj, reflect.String)
}

func isInt(obj interface{}) bool {
	return is(obj, reflect.Int) || is(obj, reflect.Float64)
}

func isPrimitive(obj interface{}) bool {
	return isBool(obj) || isString(obj) || isInt(obj)
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

func less(a, b interface{}) bool {
	switch v := a.(type) {
	case float64:
		var w float64

		if isString(b) {
			w, _ = strconv.ParseFloat(b.(string), 64)
		} else {
			w = b.(float64)
		}

		return w > v
	case string:
		var w string

		if isInt(b) {
			w = strconv.FormatFloat(b.(float64), 'f', -1, 64)
		} else {
			w = b.(string)
		}

		return w > v
	}

	return false
}

func hardEquals(a, b interface{}) bool {
	ra := reflect.ValueOf(a).Kind()
	rb := reflect.ValueOf(b).Kind()

	if ra != rb {
		return false
	}

	return equals(a, b)
}

func equals(a, b interface{}) bool {
	switch v := a.(type) {
	case float64:
		var w float64

		if isString(b) {
			w, _ = strconv.ParseFloat(b.(string), 64)
		} else {
			w = b.(float64)
		}

		return v == w
	case string:
		w := b.(string)
		return v == w
	}

	return false
}

func between(operator string, values []interface{}, data interface{}) interface{} {
	a := values[0]
	b := values[1]
	c := values[2]

	if operator == "<" {
		return less(a, b) && less(b, c)
	}

	if operator == "<=" {
		return (less(a, b) || equals(a, b)) && (less(b, c) || equals(b, c))
	}

	return false
}

func unary(operator string, value interface{}) interface{} {
	var b bool
	if isBool(value) {
		b = value.(bool)
	}
	if isInt(value) {
		b = value.(float64) > 0
	}

	if operator == "!" {
		return !b
	}

	return b
}

func _and(values []interface{}) interface{} {
	r := interface{}(true)
	v := interface{}(float64(0))

	for _, value := range values {
		if isBool(value) {
			r = interface{}(r.(bool) && value.(bool))

			continue
		}

		if value.(float64) > v.(float64) {
			v = interface{}(value)
		}
	}

	if r.(bool) && v.(float64) > 0 {
		return v
	}

	return r
}

func _or(values []interface{}) interface{} {
	r := false
	for _, value := range values {
		if isBool(value) {
			r = r || value.(bool)
			continue
		}

		if isInt(value) && value.(float64) > 0 {
			return value
		}

		r = false
	}

	return r
}

func operation(operator string, values, data interface{}) interface{} {
	if operator == "var" {
		return getVar(values, data)
	}

	if isPrimitive(values) {
		return unary(operator, values)
	}

	rp := reflect.ValueOf(values)

	parsed := values.([]interface{})

	if rp.Len() == 1 {
		return unary(operator, parsed[0].(bool))
	}

	if operator == "and" {
		return _and(parsed)
	}

	if operator == "or" {
		return _or(parsed)
	}

	if rp.Len() == 3 {
		return between(operator, parsed, data)
	}

	if operator == "<" {
		return less(parsed[0], parsed[1])
	}

	if operator == ">" {
		return less(parsed[1], parsed[0])
	}

	if operator == "<=" {
		return less(parsed[0], parsed[1]) || equals(parsed[0], parsed[1])
	}

	if operator == ">=" {
		return less(parsed[1], parsed[0]) || equals(parsed[0], parsed[1])
	}

	if operator == "===" {
		return hardEquals(parsed[0], parsed[1])
	}

	if operator == "!=" {
		return !equals(parsed[0], parsed[1])
	}

	if operator == "!==" {
		return !hardEquals(parsed[0], parsed[1])
	}

	return equals(parsed[0], parsed[1])
}

func getVar(value, data interface{}) interface{} {
	if data == nil {
		return value
	}

	var parsed string

	if isSlice(value) {
		parsed = value.([]interface{})[0].(string)
	} else if isInt(value) {
		index := int(value.(float64))
		return data.([]interface{})[index]
	} else {
		parsed = value.(string)
	}

	parts := strings.Split(parsed, ".")

	_data := data

	for _, part := range parts {
		_data = getVarFromData(part, _data, value)
	}

	return _data
}

func getVarFromData(value string, data, originalValue interface{}) interface{} {
	if !isMap(data) {
		return nil
	}

	parsedValue := data.(map[string]interface{})[value]
	if parsedValue == nil && isSlice(originalValue) {
		parsedValue = originalValue.([]interface{})[1]
	}

	if parsedValue == nil {
		return nil
	}

	switch v := parsedValue.(type) {
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

func convertToResult(result interface{}, _result interface{}) {
	value := reflect.ValueOf(result).Elem()
	target := reflect.TypeOf(result).Elem()

	switch target.Kind() {
	case reflect.Float64:
		value.SetFloat(_result.(float64))
	case reflect.String:
		value.SetString(_result.(string))
	case reflect.Bool:
		value.SetBool(_result.(bool))
	default:
		if _result == nil {
			return
		}

		log.Println("Using default conversion")

		value.Set(reflect.ValueOf(_result))
	}
}

func Apply(rules, data interface{}, result interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Result must be a pointer")
	}

	if !rv.Elem().CanSet() {
		return errors.New("Result must be addressable")
	}

	if isMap(rules) {
		convertToResult(result, apply(rules, data))

		return nil
	}

	convertToResult(result, rules)

	return nil
}

func BoolApply(rules, data interface{}) (bool, error) {
	var value bool
	err := Apply(rules, data, &value)
	return value, err
}

func IntApply(rules, data interface{}) (float64, error) {
	var value float64
	err := Apply(rules, data, &value)
	return value, err
}

func StringApply(rules, data interface{}) (string, error) {
	var value string
	err := Apply(rules, data, &value)
	return value, err
}
