package jsonlogic

import (
	"encoding/json"
	"io"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

// IsValid reads a JSON Logic rule from io.Reader and validates its syntax.
// It checks if the rule conforms to valid JSON Logic format and uses supported operators.
//
// Parameters:
//   - rule: io.Reader containing the JSON Logic rule to validate
//
// Returns:
//   - bool: true if the rule is valid, false otherwise
//
// The function returns false if the JSON cannot be parsed or if the rule contains invalid operators.
func IsValid(rule io.Reader) bool {
	var _rule any

	decoderRule := json.NewDecoder(rule)
	err := decoderRule.Decode(&_rule)
	if err != nil {
		return false
	}

	return ValidateJsonLogic(_rule)
}

// ValidateJsonLogic validates if the given rules conform to JSON Logic format.
// It recursively checks the structure and ensures all operators are supported.
//
// Parameters:
//   - rules: any value representing the JSON Logic rule to validate
//
// Returns:
//   - bool: true if the rules are valid JSON Logic, false otherwise
//
// The function handles primitives, maps (operators), slices (arrays), and variable references.
func ValidateJsonLogic(rules any) bool {
	if isVar(rules) {
		return true
	}

	if typing.IsMap(rules) {
		rulesMap := rules.(map[string]any)

		// A map with more than 1 key counts as a primitive so it's time to end recursion
		if len(rulesMap) > 1 {
			return true
		}

		for operator, value := range rulesMap {
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
	operatorsLock.RLock()
	_, isOperator := operators[op]
	operatorsLock.RUnlock()
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
