// Package jsonlogic provides a Go implementation of JSONLogic rules engine.
// JSONLogic is a way to write rules that involve logic (boolean and mathematical operations),
// consistently in JSON. It's designed to be a lightweight, portable way to share logic
// between front-end and back-end systems.
//
// The package supports all standard JSONLogic operators and allows for custom operator registration.
// Rules can be applied to data using various input/output formats including io.Reader/Writer,
// json.RawMessage, and native Go interfaces.
//
// Basic usage:
//
//	rule := strings.NewReader(`{"==":[{"var":"name"}, "John"]}`)
//	data := strings.NewReader(`{"name":"John"}`)
//	var result strings.Builder
//
//	err := jsonlogic.Apply(rule, data, &result)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// result.String() will be "true"
//
// For more examples and documentation, see: https://jsonlogic.com
package jsonlogic

import (
	"encoding/json"
	"fmt"
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

	output, err := applyInterfaceUnguarded(_rule, _data)
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

	result, err := applyInterfaceUnguarded(_rule, _data)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&result)
}

// ApplyInterface applies a transformation rule to input data using interface type assertions.
// Both rule and data must contain only JSON-compatible types: bool, float64, string, nil,
// map[string]any, and []any. Passing other numeric types (int, int32, float32, etc.) will
// return an error. Use Apply or ApplyRaw if you are working with raw JSON input.
//
// Parameters:
//   - rule: interface{} representing the transformation rule to be applied
//   - data: interface{} containing the input data to transform
//
// Returns:
//   - output: interface{} containing the transformed data
//   - err: error if unsupported types are detected or if the transformation fails
func ApplyInterface(rule, data any) (any, error) {
	if err := scanForUnsupportedTypes(rule); err != nil {
		return nil, err
	}
	if err := scanForUnsupportedTypes(data); err != nil {
		return nil, err
	}
	return applyInterfaceUnguarded(rule, data)
}

func applyInterfaceUnguarded(rule, data any) (output any, err error) {
	defer func() {
		if e := recover(); e != nil {
			// fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			err = e.(error)
		}
	}()

	if m, ok := rule.(map[string]any); ok {
		return apply(m, data), err
	}

	if s, ok := rule.([]any); ok {
		parsed := make([]any, 0, len(s))

		for _, value := range s {
			parsed = append(parsed, parseValues(value, data))
		}

		return any(parsed), nil
	}

	return rule, err
}

func scanForUnsupportedTypes(v any) error {
	switch val := v.(type) {
	case nil, bool, float64, string:
		return nil
	case map[string]any:
		for _, mv := range val {
			if err := scanForUnsupportedTypes(mv); err != nil {
				return err
			}
		}
		return nil
	case []any:
		for _, sv := range val {
			if err := scanForUnsupportedTypes(sv); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("jsonlogic: unsupported type %T; only JSON-compatible types are allowed (bool, float64, string, nil, map[string]any, []any)", val)
	}
}

// GetJsonLogicWithSolvedVars processes a JSON Logic rule by resolving variables with actual data values.
// It returns the rule with variables substituted but maintains the JSON Logic structure.
//
// Parameters:
//   - rule: json.RawMessage containing the JSON Logic rule
//   - data: json.RawMessage containing the data context for variable resolution
//
// Returns:
//   - []byte: the processed rule with resolved variables as JSON bytes
//   - error: error if unmarshaling or processing fails
//
// This is useful for debugging or when you need to see the rule with variables resolved.
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

	if m, ok := values.(map[string]any); ok {
		return apply(m, data)
	}

	inputSlice := values.([]any)
	length := len(inputSlice)
	if length == 0 {
		return inputSlice
	}

	parsed := make([]any, 0, length)

	for _, value := range inputSlice {
		if m, ok := value.(map[string]any); ok {
			parsed = append(parsed, apply(m, data))
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
