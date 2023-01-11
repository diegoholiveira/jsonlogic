package jsonlogic

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#string_to_string_comparison
func TestLessCompareStringToString(t *testing.T) {
	assert.True(t, less("a", "b"))
	assert.False(t, less("a", "a"))
	assert.False(t, less("a", "3"))
}

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#string_to_number_comparison
func TestLessCompareStringToNumber(t *testing.T) {
	assert.False(t, less("5", 3))
	assert.False(t, less("3", 3))
	assert.True(t, less("3", 5))

	assert.False(t, less("hello", 5))
	assert.False(t, less(5, "hello"))

	//console.log("5" < 3n);         // false
	//console.log("3" < 5n);         // true
}

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#number_to_number_comparison
func TestLessCompareNumberToNumber(t *testing.T) {
	assert.False(t, less(5, 3))
	assert.False(t, less(3, 3))
	assert.True(t, less(3, 5))
}

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than#comparing_boolean_null_undefined_nan
func TestLessCompareBooleanNullUndefinedNan(t *testing.T) {
	assert.False(t, less(true, false))
	assert.True(t, less(false, true))

	assert.True(t, less(0, true))
	assert.False(t, less(true, 1))

	assert.False(t, less(nil, 0))
	assert.True(t, less(nil, 1))

	assert.False(t, less(undefinedType{}, 3))
	assert.False(t, less(3, undefinedType{}))

	assert.False(t, less(3, math.NaN()))
	assert.False(t, less(math.NaN(), 3))
}

func TestLessCompareStringAndNil(t *testing.T) {
	assert.False(t, less("a", nil))
}

func TestLessCompareNumberAndNil(t *testing.T) {
	assert.False(t, less(12, nil))
}

func TestEqualsWithBooleans(t *testing.T) {
	assert.True(t, equals(true, true))
}

func TestEqualsWithNil(t *testing.T) {
	assert.True(t, equals(nil, nil))
	assert.False(t, equals(nil, ""))
}