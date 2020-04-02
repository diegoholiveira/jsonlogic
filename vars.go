package jsonlogic

import (
	"strings"
)

func solveVars(values, data interface{}) interface{} {
	if isMap(values) {
		logic := map[string]interface{}{}

		for key, value := range values.(map[string]interface{}) {
			if key == "var" {
				if value == "" || strings.HasPrefix(value.(string), ".") {
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
			_value = data.([]interface{})[int(toNumber(part))]
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
