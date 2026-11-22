package jsonlogic

import (
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/diegoholiveira/jsonlogic/v3/internal"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		rule         string
		wantCompiled bool
	}{
		{rule: `{"==":[{"var":"a"},"x"]}`, wantCompiled: true},
		{rule: `{"!==":["x",{"var":"a"}]}`},
		{rule: `{"in":[{"var":"a"},["x","y"]]}`, wantCompiled: true},
		{rule: `{"and":[{"==":[{"var":"a"},"x"]},{"!":{"in":[{"var":"b"},["1"]]}}]}`, wantCompiled: true},
		{rule: `{"or":[{"!=":[{"var":"a"},"x"]},{"===":[{"var":"b"},"1"]}]}`, wantCompiled: true},
		{rule: `{"==":[{"var":"a"},1]}`},
		{rule: `{"==":[{"var":"a"},true]}`},
		{rule: `{"==":[{"var":"a.b"},"x"]}`},
		{rule: `{"in":["x",{"var":"a"}]}`},
		{rule: `{"in":[{"var":"a"},["x",1]]}`},
		{rule: `{"and":[{"==":[{"var":"a"},"x"]},{"<":[{"var":"b"},1]}]}`},
		{rule: `{"and":[]}`},
	}
	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			var rule any
			require.NoError(t, json.Unmarshal([]byte(tt.rule), &rule))
			assert.Equal(t, tt.wantCompiled, compile(rule) != nil)
		})
	}
}

func TestPreparedRuleMatchesTestSuites(t *testing.T) {
	suites := map[string][]internal.Test{
		"Official": internal.GetScenariosFromOfficialTestSuite(),
		"Proposed": internal.GetScenariosFromProposedOfficialTestSuite(),
	}
	for name, tests := range suites {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				want, wantErr := ApplyInterface(test.Rule, test.Data)
				prepared, err := Prepare(test.Rule)
				if err != nil {
					continue
				}
				got, gotErr := prepared.Apply(test.Data)
				assert.Equal(t, wantErr, gotErr, test.Scenario)
				assert.Equal(t, want, got, "%s_%d", test.Scenario, test.Index)
			}
		})
	}
}

func TestCompiledChecksMatchApplyInterface(t *testing.T) {
	rules := []string{
		`{"==":[{"var":"a"},"x"]}`,
		`{"!=":[{"var":"a"},""]}`,
		`{"===":[{"var":"a"},"1"]}`,
		`{"!==":[{"var":"a"},"x"]}`,
		`{"in":[{"var":"a"},["x","1",""]]}`,
		`{"and":[{"==":[{"var":"a"},"x"]},{"in":[{"var":"b"},["1","2"]]}]}`,
		`{"or":[{"!=":[{"var":"a"},"x"]},{"==":[{"var":"b"},"1"]}]}`,
		`{"!":{"in":[{"var":"a"},["x"]]}}`,
		`{"!":[{"==":[{"var":"b"},"1"]}]}`,
	}
	values := []any{"x", "1", "", "y", float64(1), float64(0), true, false, nil, []any{"x"}, map[string]any{"a": "x"}}
	for _, r := range rules {
		var rule any
		require.NoError(t, json.Unmarshal([]byte(r), &rule))
		prepared, err := Prepare(rule)
		require.NoError(t, err)
		require.NotNil(t, prepared.check, r)
		for _, a := range values {
			for _, b := range values {
				data := map[string]any{"a": a, "b": b}
				want, wantErr := ApplyInterface(rule, data)
				got, gotErr := prepared.Apply(data)
				assert.Equal(t, wantErr, gotErr)
				assert.Equal(t, want, got, "%s with %v", r, data)
			}
		}
	}
}

func TestAddOperatorDisablesCompiledChecks(t *testing.T) {
	operatorsLock.RLock()
	original := operators["=="]
	operatorsLock.RUnlock()
	t.Cleanup(func() {
		operatorsLock.Lock()
		operators["=="] = original
		operatorsLock.Unlock()
		atomic.StoreInt32(&compiledChecksDisabled, 0)
	})

	prepared, err := Prepare(map[string]any{"==": []any{map[string]any{"var": "a"}, "x"}})
	require.NoError(t, err)

	AddOperator("==", func(_, _ any) any { return true })

	got, err := prepared.Apply(map[string]any{"a": "y"})
	require.NoError(t, err)
	assert.Equal(t, true, got)
}

func BenchmarkRuleApply(b *testing.B) {
	rule := map[string]any{"and": []any{
		map[string]any{"==": []any{map[string]any{"var": "env"}, "prod"}},
		map[string]any{"in": []any{map[string]any{"var": "id"}, []any{"290636", "290316"}}},
	}}
	data := map[string]any{"env": "prod", "id": "290636", "dsp": "25", "key": "req-1"}
	prepared, err := Prepare(rule)
	require.NoError(b, err)
	for name, apply := range map[string]func() (any, error){
		"ApplyInterface": func() (any, error) { return ApplyInterface(rule, data) },
		"Prepared":       func() (any, error) { return prepared.Apply(data) },
	} {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := apply(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
