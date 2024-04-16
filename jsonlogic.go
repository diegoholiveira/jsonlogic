package jsonlogic

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/barkimedes/go-deepcopy"
)

type ErrInvalidOperator struct {
	operator string
}

func (e ErrInvalidOperator) Error() string {
	return fmt.Sprintf("The operator \"%s\" is not supported", e.operator)
}

func between(operator string, values []interface{}, data interface{}) interface{} {
	a := parseValues(values[0], data)
	b := parseValues(values[1], data)
	c := parseValues(values[2], data)

	if operator == "<" {
		return less(a, b) && less(b, c)
	}

	if operator == "<=" {
		return (less(a, b) || equals(a, b)) && (less(b, c) || equals(b, c))
	}

	if operator == ">=" {
		return (less(c, b) || equals(c, b)) && (less(b, a) || equals(b, a))
	}

	return less(c, b) && less(b, a)
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

	i := sort.Search(len(valuesSlice), findElement)
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
	if value == nil || values == nil {
		return false
	}

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

func max(values interface{}) interface{} {
	converted := values.([]interface{})
	size := len(converted)
	if size == 0 {
		return nil
	}

	bigger := toNumber(converted[0])

	for i := 1; i < size; i++ {
		_n := toNumber(converted[i])
		if _n > bigger {
			bigger = _n
		}
	}

	return bigger
}

func min(values interface{}) interface{} {
	converted := values.([]interface{})
	size := len(converted)
	if size == 0 {
		return nil
	}

	smallest := toNumber(converted[0])

	for i := 1; i < size; i++ {
		_n := toNumber(converted[i])
		if smallest > _n {
			smallest = _n
		}
	}

	return smallest
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

	parsed := values.([]interface{})

	length := len(parsed)

	if length == 0 {
		return nil
	}

	for i := 0; i < length-1; i = i + 2 {
		v := parsed[i]
		if isMap(v) {
			v = getVar(parsed[i], data)
		}

		if isTrue(v) {
			return parseValues(parsed[i+1], data)
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
	modified, err := deepcopy.Anything(object)
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

	for _, value := range subject.([]interface{}) {
		conditions := solveVars(parsed[1], value)
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

	for _, value := range subject.([]interface{}) {
		v := apply(
			solveVars(
				solveVars(parsed[1], data),
				value,
			),
			value,
		)

		if isTrue(v) {
			return true
		}
	}

	return false
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
	ruleMap := rules.(map[string]interface{})

	// A map with more than 1 key counts as a primitive
	// end recursion
	if len(ruleMap) > 1 {
		return ruleMap
	}

	for operator, values := range ruleMap {
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
	if data == nil {
		data = strings.NewReader("{}")
	}

	var _rule interface{}
	var _data interface{}

	decoder := json.NewDecoder(rule)
	err := decoder.Decode(&_rule)
	if err != nil {
		return err
	}

	decoder = json.NewDecoder(data)
	err = decoder.Decode(&_data)
	if err != nil {
		return err
	}

	output, err := ApplyInterface(_rule, _data)
	if err != nil {
		return err
	}

	return json.NewEncoder(result).Encode(output)
}

func GetJsonLogicWithSolvedVars(rule, data json.RawMessage) ([]byte, error) {
	if data == nil {
		data = json.RawMessage("{}")
	}

	// parse rule and data from json.RawMessage to interface
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

	return solveVarsBackToJsonLogic(_rule, _data)
}

// ApplyRaw receives a rule and data as json.RawMessage and returns the result
// of the rule applied to the data.
func ApplyRaw(rule, data json.RawMessage) (json.RawMessage, error) {
	if data == nil {
		data = json.RawMessage("{}")
	}

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

	result, err := ApplyInterface(_rule, _data)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&result)
}

// ApplyInterface receives a rule and data as interface{} and returns the result
// of the rule applied to the data.
//
// Deprecated: Use Apply instead because ApplyInterface will be private in the next version.
func ApplyInterface(rule, data interface{}) (output interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			// fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			err = e.(error)
		}
	}()

	if isMap(rule) {
		return apply(rule, data), err
	}

	if isSlice(rule) {
		var parsed []interface{}

		for _, value := range rule.([]interface{}) {
			parsed = append(parsed, parseValues(value, data))
		}

		return interface{}(parsed), nil
	}

	return rule, err
}
