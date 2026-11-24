package jsonlogic

import (
	"fmt"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
)

// ErrReduceDataType represents an error when an unsupported data type is used in reduce operations.
// It contains the data type name that caused the error.
type ErrReduceDataType struct {
	dataType string
}

func (e ErrReduceDataType) Error() string {
	return fmt.Sprintf("The type \"%s\" is not supported", e.dataType)
}

func extractSubject(parsed []any, data any) any {
	switch v := parsed[0].(type) {
	case literal:
		// bound from data
		if s, ok := v.value.([]any); ok {
			return s
		}
	case []any:
		return copyContainers(v)
	case map[string]any:
		return apply(v, data)
	}

	return nil
}

// copyContainers returns a list that belongs to the rule without sharing its objects
// and lists, which may end up in the result.
func copyContainers(list []any) []any {
	for _, item := range list {
		if item != nil && !isPrimitive(item) {
			return deepCopyAny(list).([]any)
		}
	}
	return list
}

func filter(values, data any) any {
	parsed := values.([]any)
	if len(parsed) < 2 {
		return []any{}
	}

	subject := extractSubject(parsed, data)
	if subject == nil {
		return []any{}
	}

	subjectSlice := subject.([]any)
	subjectLen := len(subjectSlice)

	// Pre-allocate result with capacity that's reasonable for filtering
	// Assuming at least half might pass the filter (heuristic)
	result := make([]any, 0, subjectLen/2)

	logic := bindVars(parsed[1], data)

	for _, value := range subjectSlice {
		v := parseValues(logic, value)

		if javascript.IsTrue(v) {
			result = append(result, value)
		}
	}

	return result
}

func _map(values, data any) any {
	parsed := values.([]any)
	if len(parsed) < 2 {
		return []any{}
	}

	subject := extractSubject(parsed, data)
	if subject == nil {
		return []any{}
	}

	subjectSlice := subject.([]any)
	subjectLen := len(subjectSlice)

	result := make([]any, 0, subjectLen)

	logic := parsed[1]

	for _, value := range subjectSlice {
		v := parseValues(logic, value)
		result = append(result, v)
	}

	return result
}

func reduce(values, data any) any {
	parsed := values.([]any)
	if len(parsed) < 3 {
		return float64(0)
	}

	var (
		accumulator any
		valueType   string
	)

	{
		initialValue := parseValues(parsed[2], data)

		switch v := initialValue.(type) {
		case bool:
			accumulator = javascript.IsTrue(v)
			valueType = "bool"
		case float64:
			accumulator = v
			valueType = "number"
		case string:
			accumulator = v
			valueType = "string"
		default:
			panic(ErrReduceDataType{
				dataType: fmt.Sprintf("%T", unwrapLiteral(parsed[2])),
			})
		}
	}

	context := map[string]any{
		"current":     float64(0),
		"accumulator": accumulator,
		"valueType":   valueType,
	}

	subject := extractSubject(parsed, data)
	if subject == nil {
		return float64(0)
	}

	for _, value := range subject.([]any) {
		if value == nil {
			continue
		}

		context["current"] = value

		v := apply(parsed[1], context)

		switch context["valueType"] {
		case "bool":
			context["accumulator"] = javascript.IsTrue(v)
		case "number":
			context["accumulator"] = toNumber(v)
		case "string":
			context["accumulator"] = toString(v)
		}
	}

	return context["accumulator"]
}

// parseInOperands evaluates the operands of "in". A list of primitives evaluates to
// itself, so it is used as is instead of being copied. b may be a *stringSet.
func parseInOperands(values, data any) (a, b any) {
	if pair, ok := values.([]any); ok && len(pair) == 2 {
		a = parseValues(pair[0], data)
		switch list := pair[1].(type) {
		case *stringSet:
			return a, list
		case []any:
			if allPrimitive(list) {
				return a, pair[1]
			}
		}
		return a, parseValues(pair[1], data)
	}

	parsed := parseValues(values, data).([]any)
	if len(parsed) > 1 {
		b = parsed[1]
	}
	return parsed[0], b
}

func allPrimitive(list []any) bool {
	for _, item := range list {
		if item != nil && !isPrimitive(item) {
			return false
		}
	}
	return true
}

func _in(values, data any) any {
	a, b := parseInOperands(values, data)

	if set, ok := b.(*stringSet); ok {
		if s, ok := a.(string); ok {
			return set.contains(s)
		}
		b = set.list
	}

	if bs, ok := b.(string); ok {
		return strings.Contains(bs, a.(string))
	}

	bSlice, ok := b.([]any)
	if !ok {
		return false
	}

	for _, element := range bSlice {
		if es, ok := element.([]any); ok {
			if _inRange(a, es) {
				return true
			}

			continue
		}

		if _, ok := a.(float64); ok {
			if toNumber(element) == a {
				return true
			}

			continue
		}

		if element == a {
			return true
		}
	}

	return false
}

func merge(values, data any) any {
	values = parseValues(values, data)
	if isPrimitive(values) {
		return []any{values}
	}

	inputSlice := values.([]any)
	sliceLen := len(inputSlice)
	if sliceLen == 0 {
		return inputSlice
	}

	totalCapacity := 0
	for _, value := range inputSlice {
		if sv, ok := value.([]any); ok {
			totalCapacity += len(sv)
		} else {
			totalCapacity++
		}
	}

	result := make([]any, 0, totalCapacity)

	for _, value := range inputSlice {
		sv, ok := value.([]any)
		if !ok {
			result = append(result, value)
			continue
		}

		result = append(result, sv...)
	}

	return result
}

func missing(values, data any) any {
	values = parseValues(values, data)
	if _, ok := values.(string); ok {
		values = []any{values}
	}

	s := values.([]any)
	if len(s) == 0 {
		return []any{}
	}

	missing := make([]any, 0, len(s))

	for _, _var := range s {
		if isMissing(lookupVar(_var, data)) {
			missing = append(missing, _var)
		}
	}

	return missing
}

// isMissing reports whether a value counts as missing: like json-logic-js,
// both null and the empty string do.
func isMissing(value any) bool {
	return value == nil || value == ""
}

func missingSome(values, data any) any {
	parsed := parseValues(values, data).([]any)
	number := int(toNumber(parsed[0]))

	vars, ok := parsed[1].([]any)
	if !ok || len(vars) == 0 {
		return []any{}
	}

	missing := make([]any, 0, len(vars))
	foundCount := 0

	for _, _var := range vars {
		if isMissing(lookupVar(_var, data)) {
			missing = append(missing, _var)
		} else {
			foundCount++
		}
	}

	if number > foundCount {
		return missing
	}

	return []any{}
}

func all(values, data any) any {
	parsed := values.([]any)

	subject := extractSubject(parsed, data)
	if !javascript.IsTrue(subject) {
		return false
	}

	conditions := bindVars(parsed[1], data)

	for _, value := range subject.([]any) {
		v := apply(conditions, value)

		if !javascript.IsTrue(v) {
			return false
		}
	}

	return true
}

func none(values, data any) any {
	parsed := values.([]any)

	subject := extractSubject(parsed, data)

	if !javascript.IsTrue(subject) {
		return true
	}

	conditions := bindVars(parsed[1], data)

	for _, value := range subject.([]any) {
		v := apply(conditions, value)

		if javascript.IsTrue(v) {
			return false
		}
	}

	return true
}

func some(values, data any) any {
	parsed := values.([]any)
	subject := extractSubject(parsed, data)

	if !javascript.IsTrue(subject) {
		return false
	}

	conditions := bindVars(parsed[1], data)
	for _, value := range subject.([]any) {
		v := apply(conditions, value)

		if javascript.IsTrue(v) {
			return true
		}
	}

	return false
}

func _inRange(value any, values []any) bool {
	i := values[0]
	j := values[1]
	v := toNumber(value)

	return v >= toNumber(i) && toNumber(j) >= v
}
