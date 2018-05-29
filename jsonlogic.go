package jsonlogic

import (
	"errors"
	"fmt"
	//"log"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func less(a, b interface{}) bool {
	switch v := a.(type) {
	case float64:
		w := toFloat(b)
		return w > v
	case string:
		w := toString(b)
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
		w := toFloat(b)
		return v == w
	case string:
		w := toString(b)
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
	if operator == "+" || operator == "*" {
		return toFloat(value)
	}

	if operator == "-" {
		return -1 * toFloat(value)
	}

	if operator == "!!" {
		return !unary("!", value).(bool)
	}

	var b bool

	if isSlice(value) && reflect.ValueOf(value).Len() > 0 {
		b = true
	}

	if isNumber(value) {
		v := toFloat(value)
		b = v != 0
	}

	if isBool(value) {
		b = value.(bool)
	}

	if isString(value) && len(toString(value)) > 0 {
		b = true
	}

	if operator == "!" {
		return !b
	}

	return b
}

func _and(values []interface{}) interface{} {
	var v float64

	isBoolExpression := true

	for _, value := range values {
		if isSlice(value) {
			return value
		}

		if isBool(value) && !value.(bool) {
			return false
		}

		if isString(value) && toString(value) == "" {
			return value
		}

		if !isNumber(value) {
			continue
		}

		isBoolExpression = false

		_value := toFloat(value)

		if _value > v {
			v = _value
		}
	}

	if isBoolExpression {
		return true
	}

	return v
}

func _or(values []interface{}) interface{} {
	for _, value := range values {
		if isBool(value) && value.(bool) {
			return true
		}

		if isString(value) {
			if len(toString(value)) > 0 {
				return value
			}

			continue
		}

		if isNumber(value) && value.(float64) > 0 {
			return value
		}
	}

	return false
}

func _in(value interface{}, values interface{}) bool {
	if isString(values) {
		return strings.Contains(values.(string), value.(string))
	}

	for _, v := range values.([]interface{}) {
		if v == value {
			return true
		}
	}

	return false
}

func mod(a interface{}, b interface{}) interface{} {
	_a := toFloat(a)
	_b := toFloat(b)

	return math.Mod(_a, _b)
}

func concat(values interface{}) interface{} {
	if isString(values) {
		return values
	}

	var s strings.Builder
	for _, text := range values.([]interface{}) {
		if isNumber(text) {
			s.WriteString(strconv.FormatFloat(text.(float64), 'f', -1, 64))
		}

		if isString(text) {
			s.WriteString(text.(string))
		}
	}

	return strings.TrimSpace(s.String())
}

func max(values interface{}) interface{} {
	bigger := math.SmallestNonzeroFloat64

	for _, n := range values.([]interface{}) {
		_n := toFloat(n)
		if _n > bigger {
			bigger = _n
		}
	}

	return bigger
}

func min(values interface{}) interface{} {
	smallest := math.MaxFloat64

	for _, n := range values.([]interface{}) {
		_n := toFloat(n)
		if smallest > _n {
			smallest = _n
		}
	}

	return smallest
}

func sum(values interface{}) interface{} {
	sum := float64(0)

	for _, n := range values.([]interface{}) {
		sum += toFloat(n)
	}

	return sum
}

func minus(values interface{}) interface{} {
	var sum float64

	for _, n := range values.([]interface{}) {
		if sum == 0 {
			sum = toFloat(n)

			continue
		}

		sum -= toFloat(n)
	}

	return sum
}

func mult(values interface{}) interface{} {
	sum := float64(1)

	for _, n := range values.([]interface{}) {
		sum *= toFloat(n)
	}

	return sum
}

func div(values interface{}) interface{} {
	var sum float64

	for _, n := range values.([]interface{}) {
		if sum == 0 {
			sum = toFloat(n)

			continue
		}

		sum /= toFloat(n)
	}

	return sum
}

func substr(values interface{}) interface{} {
	rp := reflect.ValueOf(values)
	parsed := values.([]interface{})

	runes := []rune(toString(parsed[0]))

	from := int(toFloat(parsed[1]))
	length := len(runes)

	if from < 0 {
		from = length + from
	}

	if rp.Len() == 3 {
		length = int(toFloat(parsed[2]))
	}

	var to int
	if length < 0 {
		length = len(runes) + length
		to = length
	} else {
		to = from + length
	}

	if from < 0 {
		from, to = to, len(runes)-from
	}

	if to > len(runes) {
		to = len(runes)
	}

	if from > len(runes) {
		from = len(runes)
	}

	return string(runes[from:to])
}

func merge(values interface{}) interface{} {
	result := make([]interface{}, 0)

	if isPrimitive(values) {
		return append(result, values)
	}

	if isSlice(values) {
		for _, value := range values.([]interface{}) {
			_values := merge(value).([]interface{})

			result = append(result, _values...)
		}
	}

	return result
}

func conditional(values interface{}) interface{} {
	if isPrimitive(values) {
		return values
	}

	rp := reflect.ValueOf(values)

	length := rp.Len()

	if length == 0 {
		return nil
	}

	parsed := values.([]interface{})

	for i := 0; i < length-1; i = i + 2 {
		if isTrue(parsed[i]) {
			return parsed[i+1]
		}
	}

	if length%2 == 1 {
		return parsed[length-1]
	}

	return nil
}

func getVar(value, data interface{}) interface{} {
	if value == nil {
		return data
	}

	if isString(value) && toString(value) == "" {
		return data
	}

	if isNumber(value) {
		value = toString(value)
	}

	var _default interface{}

	if isSlice(value) { // syntax sugar
		length := reflect.ValueOf(value).Len()

		if length == 0 {
			return data
		}

		if length == 2 {
			_default = value.([]interface{})[1]
		}

		value = value.([]interface{})[0].(string)
	}

	if data == nil {
		return _default
	}

	parts := strings.Split(value.(string), ".")

	var _value interface{}

	for _, part := range parts {
		if toFloat(part) > 0 {
			_value = data.([]interface{})[int(toFloat(part))]
		} else {
			_value = data.(map[string]interface{})[part]
		}

		if _value == nil {
			return _default
		}

		if isPrimitive(_value) {
			continue
		}

		data = _value
	}

	if _value == nil {
		return _default
	}

	return _value
}

func missing(values interface{}) interface{} {
	return nil
}

func operation(operator string, values, data interface{}) interface{} {
	if operator == "missing" {
		return missing(values)
	}

	if operator == "var" {
		return getVar(values, data)
	}

	if operator == "cat" {
		return concat(values)
	}

	if operator == "substr" {
		return substr(values)
	}

	if operator == "merge" {
		return merge(values)
	}

	if operator == "if" {
		return conditional(values)
	}

	if isPrimitive(values) {
		return unary(operator, values)
	}

	if operator == "max" {
		return max(values)
	}

	if operator == "min" {
		return min(values)
	}

	rp := reflect.ValueOf(values)
	parsed := values.([]interface{})

	if rp.Len() == 1 {
		return unary(operator, parsed[0])
	}

	if operator == "+" {
		return sum(values)
	}

	if operator == "-" {
		return minus(values)
	}

	if operator == "*" {
		return mult(values)
	}

	if operator == "/" {
		return div(values)
	}

	if operator == "and" {
		return _and(parsed)
	}

	if operator == "or" {
		return _or(parsed)
	}

	if operator == "?:" {
		if parsed[0].(bool) {
			return parsed[1]
		}

		return parsed[2]
	}

	if operator == "in" {
		return _in(parsed[0], parsed[1])
	}

	if operator == "%" {
		return mod(parsed[0], parsed[1])
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

func parseValues(values, data interface{}) interface{} {
	if values == nil || isPrimitive(values) {
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

		value.Set(reflect.ValueOf(_result))
	}
}

// Apply executes the rules passed with the data as context
// and generates an result of any kind (boolean, map, string and others)
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
