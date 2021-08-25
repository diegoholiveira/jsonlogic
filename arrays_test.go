package jsonlogic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterParseTheSubjectFromFirstPosition(t *testing.T) {
	var parsed interface{}

	err := json.Unmarshal([]byte(`[
		[1,2,3,4,5],
		{"%":[{"var":""},2]}
	]`), &parsed)
	if err != nil {
		panic(err)
	}

	result := filter(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`[1,3,5]`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}

func TestFilterParseTheSubjectFromNullValue(t *testing.T) {
	var parsed interface{}

	err := json.Unmarshal([]byte(`[
		null,
		{"%":[{"var":""},2]}
	]`), &parsed)
	if err != nil {
		panic(err)
	}

	result := filter(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`[]`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}

func TestReduceSkipNullValues(t *testing.T) {
	var parsed interface{}

	err := json.Unmarshal([]byte(`[
		[1,2,null,4,5],
		{"+":[{"var":"current"}, {"var":"accumulator"}]},
		0
	]`), &parsed)
	if err != nil {
		panic(err)
	}

	result := reduce(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`12`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}

func TestReduceBoolValues(t *testing.T) {
	var parsed interface{}

	err := json.Unmarshal([]byte(`[
		[true,false,true,null],
		{"or":[{"var":"current"}, {"var":"accumulator"}]},
		false
	]`), &parsed)
	if err != nil {
		panic(err)
	}

	result := reduce(parsed, nil)

	var expected interface{}

	json.Unmarshal([]byte(`true`), &expected) // nolint:errcheck

	assert.Equal(t, expected, result)
}
