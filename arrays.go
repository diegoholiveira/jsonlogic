package jsonlogic

import (
	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

// containsAll checks if all elements in the second array exist in the first array.
// Returns true if every element of the required array is found in the search array.
//
// Example:
//
//	{"contains_all": [["a", "b", "c"], ["a", "b"]]} // true
//	{"contains_all": [["a", "b"], ["a", "b", "c"]]} // false
func containsAll(values, data any) any {
	parsed, ok := values.([]any)
	if !ok || len(parsed) != 2 {
		return false
	}

	searchArray := toAnySlice(parsed[0])
	if searchArray == nil {
		return false
	}

	requiredArray := toAnySlice(parsed[1])
	if requiredArray == nil {
		return false
	}

	// Empty required array means all are "contained"
	if len(requiredArray) == 0 {
		return true
	}

	for _, required := range requiredArray {
		if !containsElement(searchArray, required) {
			return false
		}
	}

	return true
}

// containsAny checks if any element in the second array exists in the first array.
// Returns true if at least one element of the check array is found in the search array.
//
// Example:
//
//	{"contains_any": [["a", "b", "c"], ["x", "b"]]} // true
//	{"contains_any": [["a", "b", "c"], ["x", "y"]]} // false
func containsAny(values, data any) any {
	parsed, ok := values.([]any)
	if !ok || len(parsed) != 2 {
		return false
	}

	searchArray := toAnySlice(parsed[0])
	if searchArray == nil {
		return false
	}

	checkArray := toAnySlice(parsed[1])
	if checkArray == nil {
		return false
	}

	for _, check := range checkArray {
		if containsElement(searchArray, check) {
			return true
		}
	}

	return false
}

// containsNone checks if no elements in the second array exist in the first array.
// Returns true if none of the elements of the check array are found in the search array.
//
// Example:
//
//	{"contains_none": [["a", "b", "c"], ["x", "y"]]} // true
//	{"contains_none": [["a", "b", "c"], ["x", "b"]]} // false
func containsNone(values, data any) any {
	parsed, ok := values.([]any)
	if !ok || len(parsed) != 2 {
		return true
	}

	searchArray := toAnySlice(parsed[0])
	if searchArray == nil {
		return true
	}

	checkArray := toAnySlice(parsed[1])
	if checkArray == nil {
		return true
	}

	for _, check := range checkArray {
		if containsElement(searchArray, check) {
			return false
		}
	}

	return true
}

// toAnySlice converts an interface{} to []any if possible.
func toAnySlice(value any) []any {
	if value == nil {
		return nil
	}

	if slice, ok := value.([]any); ok {
		return slice
	}

	return nil
}

// containsElement checks if an element exists in a slice using proper comparison.
func containsElement(slice []any, element any) bool {
	for _, item := range slice {
		if isEqualValue(item, element) {
			return true
		}
	}
	return false
}

// isEqualValue compares two values with type coercion for numbers.
func isEqualValue(a, b any) bool {
	// Direct equality check
	if a == b {
		return true
	}

	// Handle number comparison with type coercion
	if typing.IsNumber(a) && typing.IsNumber(b) {
		return typing.ToNumber(a) == typing.ToNumber(b)
	}

	// Handle string comparison
	if typing.IsString(a) && typing.IsString(b) {
		return a.(string) == b.(string)
	}

	return false
}
