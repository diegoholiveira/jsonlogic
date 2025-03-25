package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/jsonlogic/v3"
)

func TestHardEqualsWithNonSliceValues(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"===": 42
	}`)

	var expected json.RawMessage = json.RawMessage("false")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestHardEqualsWithSingleValueInSlice(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"===": [42]
	}`)

	var expected json.RawMessage = json.RawMessage("false")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestHardEqualsWithNilInParams(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"===": [null, 42]
	}`)

	var expected json.RawMessage = json.RawMessage("false")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))

	rule = json.RawMessage(`{
		"===": [null, null]
	}`)

	expected = json.RawMessage("true")

	output, err = jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestHardEqualsWithDifferentTypes(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"===": ["42", 42]
	}`)

	var expected json.RawMessage = json.RawMessage("false")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))

	rule = json.RawMessage(`{
		"===": ["42", "43"]
	}`)

	expected = json.RawMessage("false")

	output, err = jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))

	rule = json.RawMessage(`{
		"===": ["42", "42"]
	}`)

	expected = json.RawMessage("true")

	output, err = jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}
