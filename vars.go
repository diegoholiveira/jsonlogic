package jsonlogic

import (
	"encoding/json"
	"strconv"
	"strings"

)

func solveVars(values, data any) any {
	if m, ok := values.(map[string]any); ok {
		if len(m) == 0 {
			return m
		}
		logic := map[string]any{}

		for key, value := range m {
			if key == "var" {
				if s, ok := value.(string); ok && (s == "" || strings.HasPrefix(s, ".")) {
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

	if s, ok := values.([]any); ok {
		if len(s) == 0 {
			return s
		}
		logic := make([]any, 0, len(s))

		for _, value := range s {
			logic = append(logic, solveVars(value, data))
		}

		return logic
	}

	return values
}

func getVar(values, data any) any {
	values = parseValues(values, data)
	if values == nil {
		if !isPrimitive(data) {
			return nil
		}
		return data
	}

	if s, ok := values.(string); ok && s == "" {
		return data
	}

	if _, ok := values.(float64); ok {
		values = toString(values)
	}

	var _default any

	if v, ok := values.([]any); ok { // syntax sugar
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

		if mm, ok := _value.(map[string]any); ok {
			_value = mm[part]
		} else if sv, ok := _value.([]any); ok {
			pos := int(toNumber(part))
			if pos < 0 || pos >= len(sv) {
				return _default
			}
			_value = sv[pos]
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
	resultEscaped, err := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(resultJson)), `\\u`, `\u`))
	if err != nil {
		return nil, err
	}

	return []byte(resultEscaped), nil
}

// deepCopyMap returns a deep copy of a value produced by encoding/json:
// map[string]any, []any, or a primitive. It only handles the types that
// json.Unmarshal produces, which is all we need here.
func deepCopyMap(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v2 := range val {
			out[k] = deepCopyMap(v2)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, v2 := range val {
			out[i] = deepCopyMap(v2)
		}
		return out
	default:
		return val
	}
}

func setProperty(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok {
		return nil
	}

	if len(parsed) < 3 {
		if len(parsed) == 0 {
			return nil
		}
		return parsed[0]
	}

	object := parsed[0]

	if _, ok := object.(map[string]any); !ok {
		return object
	}

	property := parsed[1].(string)
	_modified := deepCopyMap(object).(map[string]any)
	_modified[property] = parseValues(parsed[2], data)

	return _modified
}
