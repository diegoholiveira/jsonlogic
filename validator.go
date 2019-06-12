package jsonlogic

// IsValid verifies if the JSON Logic is valid
func IsValid(rules interface{}) bool {
	if isPrimitive(rules) || isSlice(rules) || rules == nil {
		return true
	}

	return validateJsonLogic(rules)
}

func validateJsonLogic(rules interface{}) bool {
	if isVar(rules) {
		return true
	}

	if isMap(rules) {
		for operator, value := range rules.(map[string]interface{}) {
			if !isOperator(operator) {
				return false
			}

			return validateJsonLogic(value)
		}

		return false
	}

	if isSlice(rules) {
		for _, value := range rules.([]interface{}) {
			if isSlice(value) || isMap(value) {
				if validateJsonLogic(value) {
					continue
				}

				return false
			}

			if isVar(value) || isPrimitive(value) {
				continue
			}

			return false
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

	if isSlice(_var) {
		return validateJsonLogic(_var)
	}

	return isString(_var) || isNumber(_var) || _var == nil
}
