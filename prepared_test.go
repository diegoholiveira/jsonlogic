package jsonlogic

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decode(tb testing.TB, s string) any {
	tb.Helper()
	var v any
	require.NoError(tb, json.Unmarshal([]byte(s), &v))
	return v
}

func TestOptimize(t *testing.T) {
	variable := &node{operator: "var", values: "a"}
	raw := decode(t, `[{"var":"l"},{"in":[{"var":""},["x"]]}]`)

	tests := []struct {
		rule string
		want any
	}{
		{`{"==":[{"var":"a"},1]}`, &node{"==", []any{variable, float64(1)}}},
		{`{"!":{"var":"a"}}`, &node{"!", variable}},
		{`[{"var":"a"},1]`, []any{variable, float64(1)}},
		{`{"in":[{"var":"a"},["x","x"]]}`, &node{"in", []any{variable, &stringSet{
			list: []any{"x", "x"}, items: map[string]struct{}{"x": {}},
		}}}},
		{`{"in":[{"var":"a"},["x",1]]}`, &node{"in", []any{variable, []any{"x", float64(1)}}}},
		{`{"in":[{"var":"a"},[]]}`, &node{"in", []any{variable, []any{}}}},

		// plain data
		{`{"var":"a","in":["x",["x"]]}`, decode(t, `{"var":"a","in":["x",["x"]]}`)},
		{`{}`, map[string]any{}},
		{`{"filter":[{"var":"l"},{"in":[{"var":""},["x"]]}]}`, &node{"filter", raw}},
		{`{"custom":[{"var":"l"},{"in":[{"var":""},["x"]]}]}`, &node{"custom", raw}},
		{`{"and":{"var":"a"}}`, &node{"and", map[string]any{"var": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			assert.Equal(t, tt.want, optimize(decode(t, tt.rule)))
		})
	}
}

func TestPrepareOptimizesTheRule(t *testing.T) {
	prepared, err := Prepare(decode(t, `{"in":[{"var":"a"},["x"]]}`))
	require.NoError(t, err)

	root, ok := prepared.logic.(*node)
	require.True(t, ok)
	assert.IsType(t, &stringSet{}, root.values.([]any)[1])
}

// assertSameAsApplyInterface compares formatted values because results may hold NaN.
func assertSameAsApplyInterface(t *testing.T, rule, data any) bool {
	t.Helper()
	prepared, err := Prepare(rule)
	require.NoError(t, err)

	want, wantErr := ApplyInterface(rule, data)
	got, gotErr := prepared.Apply(data)
	if fmt.Sprint(wantErr) == fmt.Sprint(gotErr) && fmt.Sprintf("%#v", want) == fmt.Sprintf("%#v", got) {
		return true
	}

	r, _ := json.Marshal(rule)
	t.Errorf("rule %s with data %#v:\nApplyInterface: %#v, %v\nPrepared:       %#v, %v", r, data, want, wantErr, got, gotErr)
	return false
}

func TestPreparedRuleMatchesApplyInterfaceForEveryType(t *testing.T) {
	rules := []string{
		`{"in":[{"var":"a"},["x","1",""]]}`,
		`{"in":[{"var":"a"},["x",1,null]]}`,
		`{"in":[{"var":"a"},["x"],"ignored"]}`,
		`{"in":[{"var":"a"},{"var":"b"}]}`,
		`{"==":[{"var":"a"},{"var":"b"}]}`,
		`{"!==":[{"var":"a"},{"var":"b"}]}`,
		`{"<=":[{"var":"a"},{"var":"b"}]}`,
		`{"and":[{"==":[{"var":"a"},"x"]},{"!":[{"in":[{"var":"b"},["1"]]}]}]}`,
		`{"or":[{"var":"a"},{"var":"b"}]}`,
		`{"if":[{"var":"a"},{"cat":[{"var":"b"},"!"]},[{"var":"b"}]]}`,
		`{"+":[{"var":"a"},{"*":[{"var":"b"},2]}]}`,
		`{"merge":[{"var":"a"},[{"var":"b"},{"in":["x",["x"]]}]]}`,
		// values returned or iterated as plain data
		`{"if":[true,{"in":["x",["x"]],"var":"a"}]}`,
		`{"map":[[{"in":["x",["x"]]},1],{"var":""}]}`,
		`{"filter":[{"var":"b"},{"in":[{"var":""},["x","1"]]}]}`,
		`{"and":[{"some":[{"var":"b"},{"in":[{"var":""},["x"]]}]},{"var":"a"}]}`,
		`{"reduce":[{"var":"b"},{"cat":[{"var":"accumulator"},{"var":"current"}]},{"var":"a"}]}`,
		// errors
		`{"and":{"var":"a"}}`,
		`{"if":{"var":"a"}}`,
		`{"==":{"var":"a"}}`,
		`{"and":[true,{"unknown":[{"var":"a"}]}]}`,
	}
	values := []any{
		"x", "1", "", float64(1), float64(0), true, false, nil,
		[]any{}, []any{"x", "1"}, map[string]any{}, map[string]any{"in": []any{"x", []any{"x"}}},
	}

	for _, r := range rules {
		rule := decode(t, r)
		for _, a := range values {
			for _, b := range values {
				assertSameAsApplyInterface(t, rule, map[string]any{"a": a, "b": b})
			}
		}
	}
}

func TestPreparedRuleMatchesApplyInterfaceForRandomRules(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	names := []string{"a", "b", "a.in", "a.1", "", ".a"}
	scalars := []any{"x", "1", "", "a", float64(1), float64(0), true, nil}
	ops := []string{
		"and", "or", "if", "?:", "!", "!!", "==", "!=", "===", "!==", "<", "<=", ">", ">=",
		"in", "in", "in", "+", "-", "*", "/", "%", "min", "max", "abs", "cat", "substr", "merge",
		"missing", "missing_some", "set", "contains_all", "contains_any", "contains_none",
		"filter", "map", "reduce", "all", "some", "none", "unknown",
	}

	var value, rule func(depth int) any
	value = func(depth int) any {
		switch k := r.Intn(9); {
		case k < 5 || depth == 0:
			return scalars[r.Intn(len(scalars))]
		case k < 8:
			list := make([]any, r.Intn(4))
			for i := range list {
				list[i] = scalars[r.Intn(4)]
			}
			return list
		default:
			return map[string]any{"a": value(depth - 1), "in": value(depth - 1)}
		}
	}
	rule = func(depth int) any {
		if depth == 0 || r.Intn(5) == 0 {
			return value(2)
		}
		if r.Intn(4) == 0 {
			name := names[r.Intn(len(names))]
			if r.Intn(3) == 0 {
				return map[string]any{"var": []any{name, rule(depth - 1)}}
			}
			return map[string]any{"var": name}
		}
		op := ops[r.Intn(len(ops))]
		if r.Intn(8) == 0 {
			return map[string]any{op: rule(depth - 1)}
		}
		count := 2
		if r.Intn(3) == 0 {
			count = r.Intn(4)
		}
		args := make([]any, count)
		for i := range args {
			args[i] = rule(depth - 1)
		}
		if r.Intn(12) == 0 {
			return args
		}
		return map[string]any{op: args}
	}

	for i := 0; i < 20000; i++ {
		if !assertSameAsApplyInterface(t, rule(4), map[string]any{"a": value(2), "b": value(2)}) {
			return
		}
	}
}

func TestAddOperatorAffectsPreparedRules(t *testing.T) {
	rules := []string{
		`{"and":[{"==":[{"var":"a"},"x"]},{"!=":[{"var":"a"},"z"]}]}`,
		`{"or":[{"!":{"in":[{"var":"a"},["x","q"]]}},{"===":[{"var":"a"},"q"]},{"!==":[{"var":"a"},"x"]}]}`,
		`{"if":[{"<":[{"var":"n"},1]},"low",{"cat":["n",{"+":[{"var":"n"},1]}]}]}`,
	}
	data := map[string]any{"a": "x", "n": float64(1)}

	for _, op := range []string{"var", "and", "or", "!", "==", "!=", "===", "!==", "in", "if", "<", "+", "cat"} {
		t.Run(op, func(t *testing.T) {
			restoreOperators(t)
			for _, r := range rules {
				rule := decode(t, r)
				prepared, err := Prepare(rule)
				require.NoError(t, err)

				var calls []any
				AddOperator(op, func(values, _ any) any {
					calls = append(calls, values)
					return "x"
				})

				want, err := ApplyInterface(rule, data)
				require.NoError(t, err)
				wantCalls := calls
				calls = nil

				got, err := prepared.Apply(data)
				require.NoError(t, err)
				assert.Equal(t, want, got, r)
				assert.Equal(t, wantCalls, calls, r)
			}
		})
	}
}

func TestOverriddenInCannotChangeThePreparedRule(t *testing.T) {
	restoreOperators(t)
	prepared, err := Prepare(decode(t, `{"in":[{"var":"a"},["x"]]}`))
	require.NoError(t, err)

	AddOperator("in", func(values, _ any) any { return values.([]any)[1] })
	for i := 0; i < 2; i++ {
		got, err := prepared.Apply(map[string]any{"a": "x"})
		require.NoError(t, err)
		require.Equal(t, []any{"x"}, got)
		got.([]any)[0] = "changed"
	}
}
