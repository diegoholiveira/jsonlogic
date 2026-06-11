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
