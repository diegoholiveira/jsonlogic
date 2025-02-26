package javascript

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAt(t *testing.T) {
	tests := []struct {
		name     string
		values   []any
		index    int
		expected any
	}{
		{
			name:     "valid index",
			values:   []any{1, "test", true},
			index:    1,
			expected: "test",
		},
		{
			name:     "index out of bounds (positive)",
			values:   []any{1, 2, 3},
			index:    5,
			expected: UndefinedType{},
		},
		{
			name:     "index out of bounds (negative)",
			values:   []any{1, 2, 3},
			index:    -1,
			expected: UndefinedType{},
		},
		{
			name:     "empty array",
			values:   []any{},
			index:    0,
			expected: UndefinedType{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := At(tt.values, tt.index)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		isNaN    bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:  "undefined input",
			input: UndefinedType{},
			isNaN: true,
		},
		{
			name:     "int input",
			input:    42,
			expected: 42,
		},
		{
			name:     "float64 input",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "true boolean input",
			input:    true,
			expected: 1,
		},
		{
			name:     "false boolean input",
			input:    false,
			expected: 0,
		},
		{
			name:     "valid numeric string",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "whitespace string",
			input:    "   ",
			expected: 0,
		},
		{
			name:  "invalid numeric string",
			input: "not a number",
			isNaN: true,
		},
		{
			name:  "complex type (map)",
			input: map[string]int{"a": 1, "b": 2},
			isNaN: true,
		},
		{
			name:  "complex type (struct)",
			input: struct{ Name string }{"test"},
			isNaN: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNumber(tt.input)

			if tt.isNaN {
				assert.True(t, math.IsNaN(result), "Expected NaN result for %v", tt.input)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
