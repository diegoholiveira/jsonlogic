package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/jsonlogic/v3"
)

func TestSubOperation(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"-": [
			0,
			10
		]
	}`)

	var expected json.RawMessage = json.RawMessage("-10")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestAbsOperationWithScalar(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"abs": -42
	}`)

	var expected json.RawMessage = json.RawMessage("42")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestAbsOperationWithArray(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"abs": [-42]
	}`)

	var expected json.RawMessage = json.RawMessage("42")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestSumOperationWithEmptyArray(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"+": []
	}`)

	var expected json.RawMessage = json.RawMessage("0")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}
