package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessCompareStrings(t *testing.T) {
	assert.True(t, less("a", "b"))
}

func TestLessCompareStringAndNil(t *testing.T) {
	assert.True(t, less("a", interface{}))
}

func TestEqualsWithBooleans(t *testing.T) {
	assert.True(t, equals(true, true))
}
