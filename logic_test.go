package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestConditionalReturnsValuesFromData(t *testing.T) {
	data := json.RawMessage(`{
		"obj": {"cat": ["a", "b"]},
		"arr": [{"cat": ["a", "b"]}],
		"nested": {"if": [true, "a", "b"]}
	}`)

	cases := []struct {
		name     string
		rule     string
		expected string
	}{
		{"if then object", `{"if": [true, {"var": "obj"}]}`, `{"cat": ["a", "b"]}`},
		{"if else object", `{"if": [false, 1, {"var": "obj"}]}`, `{"cat": ["a", "b"]}`},
		{"if else-if object", `{"if": [false, 1, true, {"var": "obj"}, 2]}`, `{"cat": ["a", "b"]}`},
		{"if then array", `{"if": [true, {"var": "arr"}]}`, `[{"cat": ["a", "b"]}]`},
		{"if else nested", `{"if": [false, 1, {"var": "nested"}]}`, `{"if": [true, "a", "b"]}`},
		{"ternary then object", `{"?:": [true, {"var": "obj"}, 1]}`, `{"cat": ["a", "b"]}`},
		{"ternary else object", `{"?:": [false, 1, {"var": "obj"}]}`, `{"cat": ["a", "b"]}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := jsonlogic.ApplyRaw(json.RawMessage(tc.rule), data)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(output))
		})
	}
}

func TestConditionalEvaluatesRuleClauses(t *testing.T) {
	cases := []struct {
		name     string
		rule     string
		expected string
	}{
		{"if then", `{"if": [true, {"cat": ["a", "b"]}, "c"]}`, `"ab"`},
		{"if else", `{"if": [false, "c", {"cat": ["a", "b"]}]}`, `"ab"`},
		{"ternary", `{"?:": [true, {"cat": ["a", "b"]}, "c"]}`, `"ab"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := jsonlogic.ApplyRaw(json.RawMessage(tc.rule), nil)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(output))
		})
	}
}
