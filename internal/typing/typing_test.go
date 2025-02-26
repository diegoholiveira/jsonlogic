package typing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBool(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"true value", true, true},
		{"false value", false, true},
		{"nil value", nil, false},
		{"string value", "true", false},
		{"int value", 1, false},
		{"float value", 1.5, false},
		{"slice value", []any{}, false},
		{"map value", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBool(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"empty string", "", true},
		{"non-empty string", "hello", true},
		{"nil value", nil, false},
		{"boolean value", true, false},
		{"int value", 1, false},
		{"float value", 1.5, false},
		{"slice value", []any{}, false},
		{"map value", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"int zero", 0, true},
		{"int positive", 42, true},
		{"int negative", -10, true},
		{"float zero", 0.0, true},
		{"float positive", 3.14, true},
		{"float negative", -2.5, true},
		{"nil value", nil, false},
		{"boolean value", true, false},
		{"string value", "123", false},
		{"slice value", []any{}, false},
		{"map value", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPrimitive(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"boolean", true, true},
		{"string", "hello", true},
		{"int", 42, true},
		{"float", 3.14, true},
		{"nil value", nil, false},
		{"slice value", []any{}, false},
		{"map value", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrimitive(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsMap(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"empty map", map[string]any{}, true},
		{"non-empty map", map[string]any{"key": "value"}, true},
		{"nil value", nil, false},
		{"boolean value", true, false},
		{"int value", 1, false},
		{"float value", 1.5, false},
		{"string value", "hello", false},
		{"slice value", []any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"empty slice", []any{}, true},
		{"non-empty slice", []any{1, 2, 3}, true},
		{"nil value", nil, false},
		{"boolean value", true, false},
		{"int value", 1, false},
		{"float value", 1.5, false},
		{"string value", "hello", false},
		{"map value", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSlice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsEmptySlice(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"empty slice", []any{}, true},
		{"slice with zeros", []any{0, 0, 0}, true},
		{"slice with empty strings", []any{"", ""}, true},
		{"slice with false values", []any{false, false}, true},
		{"slice with mixed falsy values", []any{0, "", false, []any{}}, true},
		{"non-empty slice with truthy value", []any{0, 1, 0}, false},
		{"non-empty slice with true", []any{false, true}, false},
		{"nil value", nil, false},
		{"boolean value", true, false},
		{"int value", 1, false},
		{"float value", 1.5, false},
		{"string value", "hello", false},
		{"map value", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptySlice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsTrue(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"true boolean", true, true},
		{"false boolean", false, false},
		{"positive number", 42, true},
		{"negative number", -10, true},
		{"zero number", 0, false},
		{"non-empty string", "hello", true},
		{"empty string", "", false},
		{"non-empty slice", []any{1, 2, 3}, true},
		{"empty slice", []any{}, false},
		{"non-empty map", map[string]any{"key": "value"}, true},
		{"empty map", map[string]any{}, false},
		{"nil value", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTrue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
	}{
		{"int zero", 0, 0.0},
		{"int positive", 42, 42.0},
		{"int negative", -10, -10.0},
		{"float zero", 0.0, 0.0},
		{"float positive", 3.14, 3.14},
		{"float negative", -2.5, -2.5},
		{"string number integer", "42", 42.0},
		{"string number float", "3.14", 3.14},
		{"string number negative", "-10", -10.0},
		{"string empty", "", 0.0},
		{"string non-number", "hello", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"int zero", 0, "0"},
		{"int positive", 42, "42"},
		{"int negative", -10, "-10"},
		{"float zero", 0.0, "0"},
		{"float positive", 3.14, "3.14"},
		{"float negative", -2.5, "-2.5"},
		{"string", "hello", "hello"},
		{"empty string", "", ""},
		{"nil value", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}