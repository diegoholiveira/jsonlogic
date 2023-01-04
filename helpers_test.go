package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelperIsTrue(t *testing.T) {
	assert.False(t, isTrue(nil))
}

func TestToString(t *testing.T) {
	assert.Equal(t, "", toString(nil))
}

func TestToSliceOfNumbers(t *testing.T) {
	json_parsed := []interface{}{"0.0", "-10.0"}
	input := interface{}(json_parsed)
	expected := []float64{0.0, -10.0}
	assert.Equal(t, expected, toSliceOfNumbers(input))
}
