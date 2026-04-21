package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestSetProperty(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"set": [
			{"a": 1, "b": 2},
			"c",
			3
		]
	}`)

	var expected json.RawMessage = json.RawMessage(`{"a":1,"b":2,"c":3}`)

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestSetPropertyWithNonMapInput(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"set": [
			"not_a_map",
			"property",
			"value"
		]
	}`)

	var expected json.RawMessage = json.RawMessage(`"not_a_map"`)

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestGetJsonLogicWithSolvedVarsInvalidRule(t *testing.T) {
	rule := json.RawMessage(`invalid_json`)
	data := json.RawMessage(`{}`)

	_, err := jsonlogic.GetJsonLogicWithSolvedVars(rule, data)
	assert.Error(t, err)
}

func TestGetJsonLogicWithSolvedVarsInvalidData(t *testing.T) {
	rule := json.RawMessage(`{}`)
	data := json.RawMessage(`invalid_json`)

	_, err := jsonlogic.GetJsonLogicWithSolvedVars(rule, data)
	assert.Error(t, err)
}

func TestGetJsonLogicWithSolvedVarsNoData(t *testing.T) {
	rule := json.RawMessage(`{"var": "foo"}`)
	var data json.RawMessage = nil

	output, err := jsonlogic.GetJsonLogicWithSolvedVars(rule, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"var":"foo"}`
	assert.JSONEq(t, expected, string(output))
}

func TestSolveVarsBackToJsonLogicWithUnicodeChars(t *testing.T) {
	rule := json.RawMessage(`{">=":[{"var":"value"},10]}`)
	data := json.RawMessage(`{"value":20}`)

	output, err := jsonlogic.GetJsonLogicWithSolvedVars(rule, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{">=":[20,10]}`
	assert.JSONEq(t, expected, string(output))
}

func TestGetVarWithOutOfBoundsArrayIndex(t *testing.T) {
	rule := json.RawMessage(`{"var": "items.999"}`)
	data := json.RawMessage(`{"items": [1, 2, 3]}`)

	output, err := jsonlogic.ApplyRaw(rule, data)
	assert.NoError(t, err)
	assert.JSONEq(t, `null`, string(output))
}

func TestGetVarWithOutOfBoundsArrayIndexReturnsDefault(t *testing.T) {
	rule := json.RawMessage(`{"var": ["items.999", "fallback"]}`)
	data := json.RawMessage(`{"items": [1, 2, 3]}`)

	output, err := jsonlogic.ApplyRaw(rule, data)
	assert.NoError(t, err)
	assert.JSONEq(t, `"fallback"`, string(output))
}

func TestSetPropertyWithMissingValueArgument(t *testing.T) {
	rule := json.RawMessage(`{"set": [{"a": 1, "b": 2}, "c"]}`)

	output, err := jsonlogic.ApplyRaw(rule, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"a":1,"b":2}`, string(output))
}

func TestSetPropertyWithOnlyObjectArgument(t *testing.T) {
	rule := json.RawMessage(`{"set": [{"a": 1, "b": 2}]}`)

	output, err := jsonlogic.ApplyRaw(rule, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"a":1,"b":2}`, string(output))
}
