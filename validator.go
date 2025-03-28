package jsonlogic

import (
	"encoding/json"
	"io"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

// IsValid reads a JSON Logic rule from io.Reader and validates it
func IsValid(rule io.Reader) bool {
	var _rule any

	decoderRule := json.NewDecoder(rule)
	err := decoderRule.Decode(&_rule)
	if err != nil {
		return false
	}

	return ValidateJsonLogic(_rule)
}

func ValidateJsonLogic(rules any) bool {
	if isVar(rules) {
		return true
	}

	if typing.IsMap(rules) {
		for operator, value := range rules.(map[string]any) {
			if !isOperator(operator) {
				return false
			}

			return ValidateJsonLogic(value)
		}
	}

	if typing.IsSlice(rules) {
		for _, value := range rules.([]any) {
			if typing.IsSlice(value) || typing.IsMap(value) {
				if ValidateJsonLogic(value) {
					continue
				}

				return false
			}

			if isVar(value) || typing.IsPrimitive(value) {
				continue
			}
		}

		return true
	}

	return typing.IsPrimitive(rules)
}

func isOperator(op string) bool {
	_, isOperator := operators[op]
	return isOperator
}

func isVar(value any) bool {
	if !typing.IsMap(value) {
		return false
	}

	_var, ok := value.(map[string]any)["var"]
	if !ok {
		return false
	}

	return typing.IsString(_var) || typing.IsNumber(_var) || _var == nil
}
