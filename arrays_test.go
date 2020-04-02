package jsonlogic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterParseTheSubjectFromFirstPosition(t *testing.T) {
	var parsed interface{}

	json.Unmarshal([]byte(`[
		[1,2,3,4,5],
		{"%":[{"var":""},2]}
	]`), &parsed) // nolint:errcheck

	result := filter(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`[1,3,5]`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}

func TestFilterParseTheSubjectFromNullValue(t *testing.T) {
	var parsed interface{}

	json.Unmarshal([]byte(`[
		null,
		{"%":[{"var":""},2]}
	]`), &parsed) // nolint:errcheck

	result := filter(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`[]`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}

func TestReduceSkipNullValues(t *testing.T) {
	var parsed interface{}

	json.Unmarshal([]byte(`[
		[1,2,null,4,5],
		{"+":[{"var":"current"}, {"var":"accumulator"}]},
		0
	]`), &parsed) // nolint:errcheck

	result := reduce(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`12`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}
