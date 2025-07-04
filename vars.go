package jsonlogic

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/barkimedes/go-deepcopy"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func solveVars(values, data any) any {
	if typing.IsMap(values) {
		logic := map[string]any{}

		for key, value := range values.(map[string]any) {
			if key == "var" {
				if typing.IsString(value) && (value == "" || strings.HasPrefix(value.(string), ".")) {
					logic["var"] = value
					continue
				}

				val := getVar(value, data)
				if val != nil {
					return val
				}

				logic["var"] = value
			} else {
				logic[key] = solveVars(value, data)
			}
		}

		return any(logic)
	}

	if typing.IsSlice(values) {
		inputSlice := values.([]any)
		logic := make([]any, 0, len(inputSlice))

		for _, value := range inputSlice {
			logic = append(logic, solveVars(value, data))
		}

		return logic
	}

	return values
}

func getVar(values, data any) any {
	values = parseValues(values, data)
	if values == nil {
		if !typing.IsPrimitive(data) {
			return nil
		}
		return data
	}

	if typing.IsString(values) && typing.ToString(values) == "" {
		return data
	}

	if typing.IsNumber(values) {
		values = typing.ToString(values)
	}

	var _default any

	if typing.IsSlice(values) { // syntax sugar
		v := values.([]any)

		if len(v) == 0 {
			return data
		}

		if len(v) == 2 {
			_default = v[1]
		}

		values = v[0].(string)
	}

	if data == nil {
		return _default
	}

	parts := strings.Split(values.(string), ".")

	var _value any = data

	for _, part := range parts {
		if part == "" {
			continue
		}

		if typing.IsMap(_value) {
			_value = _value.(map[string]any)[part]
		} else if typing.IsSlice(_value) {
			pos := int(typing.ToNumber(part))
			container := _value.([]any)
			_value = container[pos]
		} else {
			return _default
		}

		if _value == nil {
			return _default
		}
	}

	return _value
}

func solveVarsBackToJsonLogic(rule, data any) (json.RawMessage, error) {
	ruleMap := rule.(map[string]any)
	result := make(map[string]any)

	for operator, values := range ruleMap {
		result[operator] = solveVars(values, data)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// we need to use Unquote due to unicode characters (example \u003e= need to be >=)
	// used for prettier json.RawMessage
	resultEscaped, err := strconv.Unquote(strings.Replace(strconv.Quote(string(resultJson)), `\\u`, `\u`, -1))

	if err != nil {
		return nil, err
	}

	return []byte(resultEscaped), nil
}

func setProperty(values, data any) any {
	values = parseValues(values, data).([]any)

	_value := values.([]any)

	object := _value[0]

	if !typing.IsMap(object) {
		return object
	}

	property := _value[1].(string)
	modified, err := deepcopy.Anything(object)
	if err != nil {
		panic(err)
	}

	_modified := modified.(map[string]any)
	_modified[property] = parseValues(_value[2], data)

	return any(_modified)
}
