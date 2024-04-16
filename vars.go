package jsonlogic

import (
	"encoding/json"
	"strconv"
	"strings"
)

func solveVars(values, data interface{}) interface{} {
	if isMap(values) {
		logic := map[string]interface{}{}

		for key, value := range values.(map[string]interface{}) {
			if key == "var" {
				if isString(value) && (value == "" || strings.HasPrefix(value.(string), ".")) {
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

		return interface{}(logic)
	}

	if isSlice(values) {
		logic := []interface{}{}

		for _, value := range values.([]interface{}) {
			logic = append(logic, solveVars(value, data))
		}

		return logic
	}

	return values
}

func getVar(value, data interface{}) interface{} {
	if value == nil {
		if !isPrimitive(data) {
			return nil
		}
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
		v := value.([]interface{})

		if len(v) == 0 {
			return data
		}

		if len(v) == 2 {
			_default = v[1]
		}

		value = v[0].(string)
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

		if isMap(data) {
			_value = data.(map[string]interface{})[part]
		}

		if isSlice(data) {
			pos := int(toNumber(part))
			container := data.([]interface{})
			if pos >= len(container) {
				return _default
			}
			_value = container[pos]
		}

		if _value == nil {
			return _default
		}

		if isPrimitive(_value) {
			continue
		}

		data = _value
	}

	return _value
}

func solveVarsBackToJsonLogic(rule, data interface{}) (json.RawMessage, error) {
	ruleMap := rule.(map[string]interface{})
	result := make(map[string]interface{})

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
