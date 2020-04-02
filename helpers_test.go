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
