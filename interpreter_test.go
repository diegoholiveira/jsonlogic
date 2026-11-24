package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

const hotPathRule = `{"and":[
	{"==":[{"var":"user.env"},"prod"]},
	{"!==":[{"var":"key"},"blocked"]},
	{"in":[{"var":"id"},["290636","290316",null]]}
]}`

func prepareHotPathRule(tb testing.TB) (*jsonlogic.Rule, any) {
	tb.Helper()
	var rule any
	require.NoError(tb, json.Unmarshal([]byte(hotPathRule), &rule))
	prepared, err := jsonlogic.Prepare(rule)
	require.NoError(tb, err)
	return prepared, map[string]any{
		"user": map[string]any{"env": "prod"}, "key": "req-1", "id": "290316",
	}
}

func TestHotPathDoesNotAllocate(t *testing.T) {
	rule, data := prepareHotPathRule(t)

	got, err := rule.Apply(data)
	require.NoError(t, err)
	require.Equal(t, true, got)

	allocs := testing.AllocsPerRun(100, func() { _, _ = rule.Apply(data) })
	assert.Zero(t, allocs)
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

func BenchmarkHotPath(b *testing.B) {
	rule, data := prepareHotPathRule(b)
	b.Run("single", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = rule.Apply(data)
		}
	})
	b.Run("parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = rule.Apply(data)
			}
		})
	})
}
