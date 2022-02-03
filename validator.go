package jsonlogic

import (
	"encoding/json"
	"io"
)

// IsValid reads a JSON Logic rule from io.Reader and validates it
func IsValid(rule io.Reader) bool {
	var _rule interface{}

	decoderRule := json.NewDecoder(rule)
	err := decoderRule.Decode(&_rule)
	if err != nil {
		return false
	}

	return ValidateJsonLogic(_rule)
}

func ValidateJsonLogic(rules interface{}) bool {
	if isVar(rules) {
		return true
	}

	if isMap(rules) {
		for operator, value := range rules.(map[string]interface{}) {
			if !isOperator(operator) {
				return false
			}

			return ValidateJsonLogic(value)
		}
	}

	if isSlice(rules) {
		for _, value := range rules.([]interface{}) {
			if isSlice(value) || isMap(value) {
				if ValidateJsonLogic(value) {
					continue
				}

				return false
			}

			if isVar(value) || isPrimitive(value) {
				continue
			}
		}

		return true
	}

	return isPrimitive(rules)
}

func isOperator(op string) bool {
	operators := []string{
		"==",
		"===",
		"!=",
		"!==",
		">",
		">=",
		"<",
		"<=",
		"!",
		"or",
		"and",
		"?:",
		"in",
		"in_sorted",
		"cat",
		"%",
		"abs",
		"max",
		"min",
		"+",
		"-",
		"*",
		"/",
		"substr",
		"merge",
		"if",
		"!!",
		"missing",
		"missing_some",
		"some",
		"filter",
		"map",
		"reduce",
		"all",
		"none",
		"set",
		"var",
	}

	for customOperator := range customOperators {
		operators = append(operators, customOperator)
	}

	for _, operator := range operators {
		if operator == op {
			return true
		}
	}

	return false
}

func isVar(value interface{}) bool {
	if !isMap(value) {
		return false
	}

	_var, ok := value.(map[string]interface{})["var"]
	if !ok {
		return false
	}

	return isString(_var) || isNumber(_var) || _var == nil
}
