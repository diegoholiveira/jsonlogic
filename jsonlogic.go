package jsonlogic

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"strings"
)

func less(a, b interface{}) bool {
	if isNumber(a) {
		return toNumber(b) > toNumber(a)
	}

	return toString(b) > toString(a)
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
	if isNumber(a) {
		return toNumber(a) == toNumber(b)
	}

	return toString(a) == toString(b)
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
		return toNumber(value)
	}

	if operator == "-" {
		return -1 * toNumber(value)
	}

	if operator == "!!" {
		return !unary("!", value).(bool)
	}

	b := isTrue(value)

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

		_value := toNumber(value)

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

		if isMap(value) {
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
	_a := toNumber(a)
	_b := toNumber(b)

	return math.Mod(_a, _b)
}

func concat(values interface{}) interface{} {
	if isString(values) {
		return values
	}

	var s bytes.Buffer
	for _, text := range values.([]interface{}) {
		s.WriteString(toString(text))
	}

	return strings.TrimSpace(s.String())
}

func max(values interface{}) interface{} {
	bigger := math.SmallestNonzeroFloat64

	for _, n := range values.([]interface{}) {
		_n := toNumber(n)
		if _n > bigger {
			bigger = _n
		}
	}

	return bigger
}

func min(values interface{}) interface{} {
	smallest := math.MaxFloat64

	for _, n := range values.([]interface{}) {
		_n := toNumber(n)
		if smallest > _n {
			smallest = _n
		}
	}

	return smallest
}

func sum(values interface{}) interface{} {
	sum := float64(0)

	for _, n := range values.([]interface{}) {
		sum += toNumber(n)
	}

	return sum
}

func minus(values interface{}) interface{} {
	var sum float64

	for _, n := range values.([]interface{}) {
		if sum == 0 {
			sum = toNumber(n)

			continue
		}

		sum -= toNumber(n)
	}

	return sum
}

func mult(values interface{}) interface{} {
	sum := float64(1)

	for _, n := range values.([]interface{}) {
		sum *= toNumber(n)
	}

	return sum
}

func div(values interface{}) interface{} {
	var sum float64

	for _, n := range values.([]interface{}) {
		if sum == 0 {
			sum = toNumber(n)

			continue
		}

		sum /= toNumber(n)
	}

	return sum
}

func substr(values interface{}) interface{} {
	rp := reflect.ValueOf(values)
	parsed := values.([]interface{})

	runes := []rune(toString(parsed[0]))

	from := int(toNumber(parsed[1]))
	length := len(runes)

	if from < 0 {
		from = length + from
	}

	if rp.Len() == 3 {
		length = int(toNumber(parsed[2]))
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

func conditional(values, data interface{}) interface{} {
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
		v := parsed[i]
		if isMap(v) {
			v = getVar(parsed[i], data)
		}

		if isTrue(v) {
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
		if part == "" {
			continue
		}

		if toNumber(part) > 0 {
			_value = data.([]interface{})[int(toNumber(part))]
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

func missing(values, data interface{}) interface{} {
	if isString(values) {
		values = []interface{}{values}
	}

	missing := make([]interface{}, 0)

	for _, _var := range values.([]interface{}) {
		_value := getVar(_var, data)

		if _value == nil {
			missing = append(missing, _var)
		}
	}

	return missing
}

func missingSome(values, data interface{}) interface{} {
	parsed := values.([]interface{})
	number := int(toNumber(parsed[0]))
	vars := parsed[1]

	missing := make([]interface{}, 0)
	found := make([]interface{}, 0)

	for _, _var := range vars.([]interface{}) {
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

	return make([]interface{}, 0)
}

func filter(values, data interface{}) interface{} {
	parsed := values.([]interface{})
	subject := apply(parsed[0], data)

	result := make([]interface{}, 0)
	for _, value := range subject.([]interface{}) {
		v := parseValues(parsed[1], value)

		if isBool(v) && v.(bool) {
			result = append(result, value)
		}

		if isNumber(v) && toNumber(v) != 0 {
			result = append(result, value)
		}
	}

	return result
}

func _map(values, data interface{}) interface{} {
	parsed := values.([]interface{})
	subject := apply(parsed[0], data)

	result := make([]interface{}, 0)

	if subject == nil {
		return result
	}

	for _, value := range subject.([]interface{}) {
		v := parseValues(parsed[1], value)
		if v == nil {
			continue
		}

		if isNumber(v) && toNumber(v) != 0 {
			result = append(result, toNumber(v))
		}
	}

	return result
}

func reduce(values, data interface{}) interface{} {
	parsed := values.([]interface{})
	subject := apply(parsed[0], data)

	if subject == nil {
		return float64(0)
	}

	context := map[string]interface{}{
		"current":     float64(0),
		"accumulator": toNumber(parsed[2]),
	}

	for _, value := range subject.([]interface{}) {
		context["current"] = value

		v := apply(parsed[1], context)

		if v == nil {
			continue
		}

		context["accumulator"] = toNumber(v)
	}

	return context["accumulator"]
}

func all(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	subject := apply(parsed[0], data)

	if !isTrue(subject) {
		return false
	}

	conditions := parsed[1]

	for _, value := range subject.([]interface{}) {
		v := apply(conditions, value)

		if !isTrue(v) {
			return false
		}
	}

	return true
}

func none(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	subject := apply(parsed[0], data)

	if !isTrue(subject) {
		return true
	}

	conditions := parsed[1]

	for _, value := range subject.([]interface{}) {
		v := apply(conditions, value)

		if isTrue(v) {
			return false
		}
	}

	return true
}

func some(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	subject := apply(parsed[0], data)

	if !isTrue(subject) {
		return false
	}

	conditions := parsed[1]

	for _, value := range subject.([]interface{}) {
		v := apply(conditions, value)

		if isTrue(v) {
			return true
		}
	}

	return false
}

func operation(operator string, values, data interface{}) interface{} {
	if operator == "missing" {
		return missing(values, data)
	}

	if operator == "missing_some" {
		return missingSome(values, data)
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
		return conditional(values, data)
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

	if isMap(values) {
		return apply(values, data)
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
		if operator == "filter" {
			return filter(values, data)
		}

		if operator == "map" {
			return _map(values, data)
		}

		if operator == "reduce" {
			return reduce(values, data)
		}

		if operator == "all" {
			return all(values, data)
		}

		if operator == "none" {
			return none(values, data)
		}

		if operator == "some" {
			return some(values, data)
		}

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
// and generates an result of any kind (bool, map, string and others)
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
