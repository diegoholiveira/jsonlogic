package jsonlogic

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

// Apply reads a rule and data from `io.Reader`, applies the rule to the data
// and writes the result to the provided writer. It returns an error if rule
// processing or data handling fails.
//
// Parameters:
//   - rule: io.Reader representing the transformation rule to be applied
//   - data: io.Reader containing the input data to transform
//   - result: io.Writer containing the transformed data
//
// Returns:
//   - err: error if the transformation fails or if type assertions are invalid
func Apply(rule, data io.Reader, result io.Writer) error {
	if data == nil {
		data = strings.NewReader("{}")
	}

	var _rule any
	var _data any

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

// ApplyRaw applies a validation rule to a JSON data input, both provided as raw JSON messages.
// It processes the input data according to the provided rule and returns the transformed result.
//
// Parameters:
//   - rule: json.RawMessage representing the transformation rule to be applied
//   - data: json.RawMessage containing the input data to transform
//
// Returns:
//   - output: json.RawMessage containing the transformed data
//   - err: error if the transformation fails or if type assertions are invalid
func ApplyRaw(rule, data json.RawMessage) (json.RawMessage, error) {
	if data == nil {
		data = json.RawMessage("{}")
	}

	var _rule any
	var _data any

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

// ApplyInterface applies a transformation rule to input data using interface type assertions.
// It processes the input data according to the provided rule and returns the transformed result.
//
// Parameters:
//   - rule: interface{} representing the transformation rule to be applied
//   - data: interface{} containing the input data to transform
//
// Returns:
//   - output: interface{} containing the transformed data
//   - err: error if the transformation fails or if type assertions are invalid
func ApplyInterface(rule, data any) (output any, err error) {
	defer func() {
		if e := recover(); e != nil {
			// fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			err = e.(error)
		}
	}()

	if typing.IsMap(rule) {
		return apply(rule, data), err
	}

	if typing.IsSlice(rule) {
		inputSlice := rule.([]any)
		parsed := make([]any, 0, len(inputSlice))

		for _, value := range inputSlice {
			parsed = append(parsed, parseValues(value, data))
		}

		return any(parsed), nil
	}

	return rule, err
}

func GetJsonLogicWithSolvedVars(rule, data json.RawMessage) ([]byte, error) {
	if data == nil {
		data = json.RawMessage("{}")
	}

	// parse rule and data from json.RawMessage to interface
	var _rule any
	var _data any

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

func parseValues(values, data any) any {
	if values == nil || typing.IsPrimitive(values) {
		return values
	}

	if typing.IsMap(values) {
		return apply(values, data)
	}

	inputSlice := values.([]any)
	length := len(inputSlice)
	if length == 0 {
		return inputSlice
	}

	parsed := make([]any, 0, length)

	for _, value := range inputSlice {
		if typing.IsMap(value) {
			parsed = append(parsed, apply(value, data))
		} else {
			parsed = append(parsed, parseValues(value, data))
		}
	}

	return parsed
}

func apply(rules, data any) any {
	ruleMap := rules.(map[string]any)

	// A map with more than 1 key counts as a primitive so it's time to end recursion
	if len(ruleMap) > 1 {
		return ruleMap
	}

	for operator, values := range ruleMap {
		return operation(operator, values, data)
	}

	return make(map[string]any)
}
