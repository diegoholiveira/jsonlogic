package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestCommonRulesDoNotAllocate(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{"env": "prod"}, "key": "req-1", "id": "290316",
	}
	lists := map[string]string{
		"indexed list": `["290636","290316"]`,
		"scanned list": `["290636","290316",null]`,
	}
	for name, list := range lists {
		t.Run(name, func(t *testing.T) {
			var rule any
			require.NoError(t, json.Unmarshal([]byte(`{"and":[
				{"==":[{"var":"user.env"},"prod"]},
				{"!==":[{"var":"key"},"blocked"]},
				{"in":[{"var":"id"},`+list+`]}
			]}`), &rule))
			prepared, err := jsonlogic.Prepare(rule)
			require.NoError(t, err)

			got, err := prepared.Apply(data)
			require.NoError(t, err)
			require.Equal(t, true, got)

			allocs := testing.AllocsPerRun(100, func() { _, _ = prepared.Apply(data) })
			assert.Zero(t, allocs)
		})
	}
}

func TestEqualityAndInOperandShapes(t *testing.T) {
	data := `{"x":"1","n":1,"list":["a","1"],"pair":["1",["1","2"]],"":"empty"}`
	tests := []struct{ rule, want string }{
		{`{"==":[{"var":"x"},{"var":"n"}]}`, `true`},
		{`{"===":[{"var":"x"},{"var":"n"}]}`, `false`},
		{`{"===":[{"var":"gone"},null]}`, `true`},
		{`{"==":["1","1","ignored"]}`, `true`},
		{`{"===":["1","1","ignored"]}`, `true`},
		{`{"==":["1"]}`, `false`},
		{`{"===":[]}`, `false`},
		// operands computed by a single rule
		{`{"==":{"var":"pair"}}`, `false`},
		{`{"in":{"var":"pair"}}`, `true`},

		{`{"in":[{"var":"x"},["a","1"]]}`, `true`},
		{`{"in":[{"var":"n"},["a","1"]]}`, `true`},
		{`{"in":[{"var":"gone"},["a",null]]}`, `true`},
		{`{"in":[{"var":"gone"},["a"]]}`, `false`},
		{`{"in":[{"var":"x"},[]]}`, `false`},
		{`{"in":[{"var":"x"},{"var":"list"}]}`, `true`},
		{`{"in":["x"]}`, `false`},
		{`{"in":["ell","hello"]}`, `true`},
		// lists that still need evaluating
		{`{"in":["1",["a",{"var":"x"}]]}`, `true`},
		{`{"in":[2,[[1,3]]]}`, `true`},
		{`{"in":[{"var":"n"},[[{"var":"n"},3]]]}`, `true`},

		// dotted paths with empty segments
		{`{"var":"list.1"}`, `"1"`},
		{`{"var":".list..1."}`, `"1"`},
		{`{"var":"."}`, data},
		{`{"var":["list.9","default"]}`, `"default"`},
	}
	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			got, err := jsonlogic.ApplyRaw(json.RawMessage(tt.rule), json.RawMessage(data))
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

// Operands that do not change the result are still evaluated.
func TestExtraOperandsAreEvaluated(t *testing.T) {
	for _, op := range []string{"==", "!=", "===", "!==", "in"} {
		rule := map[string]any{op: []any{"1", "1", map[string]any{"unknown": []any{}}}}
		_, err := jsonlogic.ApplyInterface(rule, nil)
		assert.Error(t, err, op)
	}
}
