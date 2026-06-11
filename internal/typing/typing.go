// Package typing provides type checking and conversion utilities for JSON data types.
package typing

import (
	"reflect"
	"strconv"
)

func is(obj any, kind reflect.Kind) bool {
	return obj != nil && reflect.TypeOf(obj).Kind() == kind
}

// IsBool checks if the provided value is a boolean type.
// Returns false if the value is nil.
//
// Example:
//
//	IsBool(true)   // Returns: true
//	IsBool(false)  // Returns: true
//	IsBool("true") // Returns: false
//	IsBool(nil)    // Returns: false
func IsBool(obj any) bool {
	return is(obj, reflect.Bool)
}

// IsString checks if the provided value is a string type.
// Returns false if the value is nil.
//
// Example:
//
//	IsString("test")  // Returns: true
//	IsString("")      // Returns: true
//	IsString(42)      // Returns: false
//	IsString(nil)     // Returns: false
func IsString(obj any) bool {
	return is(obj, reflect.String)
}

// IsNumber checks if the provided value is a numeric type (float64).
// Returns false for any other type including nil.
//
// Example:
//
//	IsNumber(3.14)     // Returns: true
//	IsNumber("42")     // Returns: false
//	IsNumber(nil)      // Returns: false
func IsNumber(obj any) bool {
	_, ok := obj.(float64)
	return ok
}

// IsPrimitive checks if the provided value is a primitive type (boolean, string, or number).
// Returns false if the value is nil or any other type.
//
// Example:
//
//	IsPrimitive(42)      // Returns: true
//	IsPrimitive("test")  // Returns: true
//	IsPrimitive(true)    // Returns: true
//	IsPrimitive([])      // Returns: false
//	IsPrimitive(nil)     // Returns: false
func IsPrimitive(obj any) bool {
	return IsBool(obj) || IsString(obj) || IsNumber(obj)
}

// IsMap checks if the provided value is a map type.
// Returns false if the value is nil.
//
// Example:
//
//	IsMap(map[string]int{"a": 1})  // Returns: true
//	IsMap(map[string]any{})        // Returns: true
//	IsMap([]int{1, 2, 3})          // Returns: false
//	IsMap(nil)                     // Returns: false
func IsMap(obj any) bool {
	return is(obj, reflect.Map)
}

// IsSlice checks if the provided value is a slice type.
// Returns false if the value is nil.
//
// Example:
//
//	IsSlice([]int{1, 2, 3})  // Returns: true
//	IsSlice([]any{})         // Returns: true
//	IsSlice("test")          // Returns: false
//	IsSlice(nil)             // Returns: false
func IsSlice(obj any) bool {
	return is(obj, reflect.Slice)
}

// ToNumber converts the provided value to a float64.
// If the value is a string, it attempts to parse it as a float64.
// If the value is already a float64, it returns it as is.
//
// Example:
//
//	ToNumber(3.14)               // Returns: 3.14
//	ToNumber("42")               // Returns: 42.0
//	ToNumber("3.14")             // Returns: 3.14
//	ToNumber("invalid")          // Returns: 0.0
func ToNumber(value any) float64 {
	if IsString(value) {
		w, _ := strconv.ParseFloat(value.(string), 64)
		return w
	}
	return value.(float64)
}

// ToString converts the provided value to a string.
// For numbers: converts to string representation
// For nil: returns an empty string
// For other types: performs a direct type assertion to string
//
// Example:
//
//	ToString(3.14)      // Returns: "3.14"
//	ToString("test")    // Returns: "test"
//	ToString(nil)       // Returns: ""
func ToString(value any) string {
	if IsNumber(value) {
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	}
	if value == nil {
		return ""
	}
	return value.(string)
}
