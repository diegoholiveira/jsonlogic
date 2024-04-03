package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelperIsTrue(t *testing.T) {
	assert.False(t, isTrue(nil))
}

func TestToSliceOfNumbers(t *testing.T) {
	json_parsed := []interface{}{"0.0", "-10.0"}
	input := interface{}(json_parsed)
	expected := []float64{0.0, -10.0}
	assert.Equal(t, expected, toSliceOfNumbers(input))
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect bool
	}{
		{
			name:   "int max",
			input:  9223372036854775807,
			expect: true,
		},
		{
			name:   "int min",
			input:  -9223372036854775808,
			expect: true,
		},
		{
			name:   "float",
			input:  1.121,
			expect: true,
		},
		{
			name:   "string",
			input:  "123",
			expect: false,
		},
		{
			name:   "boolean",
			input:  true,
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, isNumber(test.input))
		})
	}
}

func TestToNumber(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect float64
	}{
		{
			name:   "int max value",
			input:  9223372036854775807,
			expect: 9223372036854775807,
		},
		{
			name:   "int min value",
			input:  -9223372036854775808,
			expect: -9223372036854775808,
		},
		{
			name:   "int min value",
			input:  -9223372036854775808,
			expect: -9223372036854775808,
		},
		{
			name:   "int as a string",
			input:  "123",
			expect: 123,
		},
		{
			name:   "float",
			input:  1.121,
			expect: 1.121,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, toNumber(test.input))
		})
	}
}

func TestToString(t *testing.T) {
	var tests = []struct {
		name   string
		input  interface{}
		expect string
	}{
		{
			name:   "nil value",
			input:  nil,
			expect: "",
		},
		{
			name:   "string value",
			input:  "value",
			expect: "value",
		},
		{
			name:   "int value",
			input:  9223372036854775807,
			expect: "9223372036854775807",
		},
		{
			name:   "float value",
			input:  1.21,
			expect: "1.21",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, toString(test.input))
		})
	}
}
