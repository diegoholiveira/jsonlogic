package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
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

func TestArithmeticCoercion(t *testing.T) {
	// These cases cover the affected json-logic/.github arithmetic records at
	// 22aac2a67bf07b48f952d6b0f61dfe70947d1a23 without vendoring the suite data.
	tests := []struct {
		name    string
		rule    string
		want    string
		wantErr string
	}{
		{name: "divide case 6 booleans", rule: `{ "/": [false, true] }`, want: `0`},
		{name: "divide case 7 multiple types", rule: `{ "/": ["8", true, 2] }`, want: `4`},
		{name: "minus case 5 booleans", rule: `{ "-": [true, false, true] }`, want: `0`},
		{name: "minus case 6 multiple types", rule: `{ "-": [1, "2", 3, "4", "", true, false, null] }`, want: `-9`},
		{name: "minus case 20 invalid string", rule: `{ "-": ["Hey", 1] }`, wantErr: "NaN"},
		{name: "modulo case 8 booleans", rule: `{ "%": [false, true] }`, want: `0`},
		{name: "modulo case 27 invalid string", rule: `{ "%": ["Hey", 1] }`, wantErr: "NaN"},
		{name: "multiply case 5 booleans", rule: `{ "*": [true, false, true] }`, want: `0`},
		{name: "multiply case 6 multiple types", rule: `{ "*": [1, "2", 3, "4", true] }`, want: `24`},
		{name: "multiply case 8 null", rule: `{ "*": ["1", null] }`, want: `0`},
		{name: "multiply case 25 invalid string", rule: `{ "*": ["Hey", 1] }`, wantErr: "NaN"},
		{name: "plus case 5 booleans", rule: `{ "+": [true, false, true] }`, want: `2`},
		{name: "plus case 6 multiple types", rule: `{ "+": [1, "2", 3, "4", "", true, false, null] }`, want: `11`},
		{name: "plus case 19 true", rule: `{ "+": true }`, want: `1`},
		{name: "plus case 20 false", rule: `{ "+": false }`, want: `0`},
		{name: "plus case 22 null", rule: `{ "+": null }`, want: `0`},
		{name: "plus case 24 invalid string", rule: `{ "+": ["Hey", 1] }`, wantErr: "NaN"},
		{name: "plus case 28 invalid scalar string", rule: `{ "+": "Hello" }`, wantErr: "NaN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := jsonlogic.ApplyRaw(json.RawMessage(tt.rule), json.RawMessage("null"))
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			if !assert.NoError(t, err) {
				return
			}
			assert.JSONEq(t, tt.want, string(output))
		})
	}
}
