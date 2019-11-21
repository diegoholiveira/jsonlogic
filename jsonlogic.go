package jsonlogic

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/copystructure"
)

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
	if operator == "+" || operator == "*" || operator == "/" {
		return toNumber(value)
	}

	if operator == "-" {
		return -1 * toNumber(value)
	}

	if operator == "!!" {
		return !unary("!", value).(bool)
	}

	if operator == "abs" {
		return abs(value)
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
		if isTrue(value) {
			return value
		}
	}

	return false
}

func _inRange(value interface{}, values interface{}) bool {
	v := values.([]interface{})

	i := v[0]
	j := v[1]

	if isNumber(value) {
		return toNumber(value) >= toNumber(i) && toNumber(j) >= toNumber(value)
	}

	return toString(value) >= toString(i) && toString(j) >= toString(value)
}

// Expect values to be in alphabetical ascending order
func _inSorted(value interface{}, values interface{}) bool {
	valuesSlice := values.([]interface{})

	findElement := func(i int) bool {
		element := valuesSlice[i]

		if isSlice(valuesSlice[i]) {
			sliceElement := valuesSlice[i].([]interface{})
			start := sliceElement[0]
			end := sliceElement[1]

			return (toString(start) <= toString(value) && toString(end) >= toString(value)) || toString(end) > toString(value)
		}

		return toString(element) >= toString(value)
	}

	i := sort.Search(len(valuesSlice),  findElement)
	if i >= len(valuesSlice) {
		return false
	}

	if isSlice(valuesSlice[i]) {
		sliceElement := valuesSlice[i].([]interface{})
		start := sliceElement[0]
		end := sliceElement[1]

		return toString(start) <= toString(value) && toString(end) >= toString(value)
	}

	return toString(valuesSlice[i]) == toString(value)
}

func _in(value interface{}, values interface{}) bool {
	if isString(values) {
		return strings.Contains(values.(string), value.(string))
	}

	for _, element := range values.([]interface{}) {
		if isSlice(element) {
			if _inRange(value, element) {
				return true
			}

			continue
		}

		if isNumber(value) {
			if toNumber(element) == value {
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

func mod(a interface{}, b interface{}) interface{} {
	_a := toNumber(a)
	_b := toNumber(b)

	return math.Mod(_a, _b)
}

func abs(a interface{}) interface{} {
	_a := toNumber(a)

	return math.Abs(_a)
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

func merge(values interface{}, level int8) interface{} {
	result := make([]interface{}, 0)

	if isPrimitive(values) || level > 1 {
		return append(result, values)
	}

	if isSlice(values) {
		for _, value := range values.([]interface{}) {
			_values := merge(value, level+1).([]interface{})

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

func setProperty(value, data interface{}) interface{} {
	_value := value.([]interface{})

	object := _value[0]

	if !isMap(object) {
		return object
	}

	property := _value[1].(string)
	modified, err := copystructure.Copy(object)
	if err != nil {
		panic(err)
	}

	_modified := modified.(map[string]interface{})
	_modified[property] = parseValues(_value[2], data)

	return interface{}(_modified)
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

func all(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	var subject interface{}

	if isMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if isSlice(parsed[0]) {
		subject = parsed[0]
	}

	if !isTrue(subject) {
		return false
	}

	conditions := solveVars(parsed[1], data)

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

	var subject interface{}

	if isMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if isSlice(parsed[0]) {
		subject = parsed[0]
	}

	if !isTrue(subject) {
		return true
	}

	conditions := solveVars(parsed[1], data)

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

	var subject interface{}

	if isMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if isSlice(parsed[0]) {
		subject = parsed[0]
	}

	if !isTrue(subject) {
		return false
	}

	conditions := solveVars(parsed[1], data)

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

	if operator == "set" {
		return setProperty(values, data)
	}

	if operator == "cat" {
		return concat(values)
	}

	if operator == "substr" {
		return substr(values)
	}

	if operator == "merge" {
		return merge(values, 0)
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

	if operator == "in_sorted" {
		return _inSorted(parsed[0], parsed[1])
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

	// an empty-map rule should return an empty-map
	return make(map[string]interface{})
}

// Apply read the rule and it's data from io.Reader, executes it
// and write back a JSON into an io.Writer result
func Apply(rule, data io.Reader, result io.Writer) error {
	var _rule interface{}
	var _data interface{}

	decoderRule := json.NewDecoder(rule)
	err := decoderRule.Decode(&_rule)
	if err != nil {
		return err
	}

	decoderData := json.NewDecoder(data)
	err = decoderData.Decode(&_data)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(result)

	if isMap(_rule) {
		encoder.Encode(apply(_rule, _data))
	} else {
		encoder.Encode(_rule)
	}

	return nil
}

func ApplyRaw(rule, data json.RawMessage) (json.RawMessage, error) {
	var _rule interface{}
	var _data interface{}

	err := json.Unmarshal(rule, &_rule)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &_data)
	if err != nil {
		return nil, err
	}

	var result interface{}

	if isMap(_rule) {
		result = apply(_rule, _data)
	} else {
		result = _rule
	}

	var output json.RawMessage

	output, err = json.Marshal(&result)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func ApplyInterface(rule, data interface{}) (interface{}, error) {
	var result interface{}

	if isMap(rule) {
		result = apply(rule, data)
	} else {
		result = rule
	}

	return result, nil
}
