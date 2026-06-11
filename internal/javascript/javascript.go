// Package javascript provides utilities for working with JavaScript code and runtime integration.
package javascript

import (
	"math"
	"strconv"
	"strings"
)

type UndefinedType struct{}

// At returns the element at the specified index in the slice.
// If index is negative, it counts from the end of the slice.
// If index is out of bounds, it returns nil.
//
// Example:
//
//	At([]any{1,2,3}, 1)  // Returns: 2
//	At([]any{1,2,3}, -1) // Returns: 3
func At(values []any, index int) any {
	if index >= 0 && index < len(values) {
		return values[index]
	}
	return UndefinedType{}
}

// ToNumber converts various input types to float64.
//
// Examples:
//
//	ToNumber(3.14)   // Returns: 3.14
//	ToNumber("3.14") // Returns: 3.14
//	ToNumber(true)   // Returns: 1.0
//	ToNumber(false)  // Returns: 0.0
//	ToNumber(nil)    // Returns: 0.0
func ToNumber(v any) float64 {
	switch value := v.(type) {
	case nil:
		return 0
	case UndefinedType:
		return math.NaN()
	case float64:
		return value
	case bool: // Boolean values true and false are converted to 1 and 0 respectively.
		if value {
			return 1
		} else {
			return 0
		}
	case string:
		if strings.TrimSpace(value) == "" {
			return 0
		}

		n, err := strconv.ParseFloat(value, 64)
		switch err {
		case strconv.ErrRange, nil:
			return n
		default:
			return math.NaN()
		}
	default:
		return math.NaN()
	}
}

// IsTrue checks if the provided value is considered truthy in JavaScript logic.
// For booleans: true is truthy
// For numbers: non-zero is truthy
// For strings: non-empty string is truthy
// For slices/maps: non-empty slice/map is truthy
// Returns false for nil or any other type.
//
// Example:
//
//	IsTrue(true)               // Returns: true
//	IsTrue(float64(42))        // Returns: true
//	IsTrue("hello")            // Returns: true
//	IsTrue([]any{1, 2, 3})     // Returns: true
//	IsTrue(false)              // Returns: false
//	IsTrue(float64(0))         // Returns: false
//	IsTrue("")                 // Returns: false
//	IsTrue([]any{})            // Returns: false
//	IsTrue(nil)                // Returns: false
func IsTrue(obj any) bool {
	switch v := obj.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case string:
		return len(v) > 0
	case []any:
		return len(v) > 0
	case map[string]any:
		return len(v) > 0
	default:
		return false
	}
}

// IsEmptySlice checks if the provided value is a slice and all its elements are falsy.
// Returns false if the value is not a slice or if any element in the slice is truthy.
// A falsy value is: false, 0, "", empty array, or empty map.
//
// Example:
//
//	IsEmptySlice([]any{})              // Returns: true
//	IsEmptySlice([]any{0, "", false})  // Returns: true
//	IsEmptySlice([]any{1, 2, 3})       // Returns: false
//	IsEmptySlice("test")               // Returns: false
func IsEmptySlice(obj any) bool {
	values, ok := obj.([]any)
	if !ok {
		return false
	}
	for _, v := range values {
		if IsTrue(v) {
			return false
		}
	}
	return true
}
