package jsonlogic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinusWithOneNumber(t *testing.T) {
	json_parsed := []interface{}{"-10.0"}
	input := interface{}(json_parsed)
	expected := -10.0
	assert.Equal(t, expected, minus(input))
}

func TestDivWithOneNumber(t *testing.T) {
	json_parsed := []interface{}{"2.0"}
	input := interface{}(json_parsed)
	expected := 2.0
	assert.Equal(t, expected, minus(input))
}

func TestSubOperation(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"-": [
			0,
			10
		]
	}`)

	var expected json.RawMessage = json.RawMessage("-10")

	output, err := ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}
