package jsonlogic

import (
	"encoding/json"
	"strconv"
	"strings"
)

// literal wraps a value bound from data by bindVars; evaluating it returns the
// value unchanged.
type literal struct {
	value any
}

// unwrapLiteral returns the value held by a literal, or v unchanged.
func unwrapLiteral(v any) any {
	if l, ok := v.(literal); ok {
		return l.value
	}
	return v
}

// solveVars substitutes resolved vars into the rule for display. Use bindVars
// for logic that will be evaluated.
func solveVars(values, data any) any {
	return resolveVars(values, data, false)
}

// bindVars substitutes outer-scope vars into logic that will be evaluated
// against another scope, wrapping non-primitive values in literal.
func bindVars(values, data any) any {
	return resolveVars(values, data, true)
}

func resolveVars(values, data any, bind bool) any {
	if m, ok := values.(map[string]any); ok {
		if len(m) == 0 {
			return m
		}
		// apply treats multi-key maps as values.
		if bind && len(m) > 1 {
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
					if bind && !isPrimitive(val) {
						return literal{value: val}
					}
					return val
				}

				logic["var"] = value
			} else {
				logic[key] = resolveVars(value, data, bind)
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
			logic = append(logic, resolveVars(value, data, bind))
		}

		return logic
	}

	return values
}

func getVar(values, data any) any {
	return lookupVar(parseValues(values, data), data)
}

// lookupVar resolves an already-evaluated var path against data.
func lookupVar(values, data any) any {
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

	_value := data

	for _, part := range parts {
		if part == "" {
			continue
		}

		if mm, ok := _value.(map[string]any); ok {
			_value = mm[part]
		} else if sv, ok := _value.([]any); ok {
			f, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return _default
			}
			pos := int(f)
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

// deepCopyMap returns a recursively deep-copied map[string]any.
func deepCopyMap(v map[string]any) map[string]any {
	out := make(map[string]any, len(v))
	for k, v2 := range v {
		out[k] = deepCopyAny(v2)
	}
	return out
}

// deepCopyAny is the recursive helper for deepCopyMap.
// It handles map[string]any, []any, and primitive values, the only types produced by encoding/json.
func deepCopyAny(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return deepCopyMap(val)
	case []any:
		out := make([]any, len(val))
		for i, v2 := range val {
			out[i] = deepCopyAny(v2)
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

	object, ok := parsed[0].(map[string]any)
	if !ok {
		return parsed[0]
	}

	property, ok := parsed[1].(string)
	if !ok {
		return parsed[0]
	}

	modified := deepCopyMap(object)
	modified[property] = parsed[2]

	return modified
}
