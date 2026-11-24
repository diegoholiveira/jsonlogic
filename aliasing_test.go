package jsonlogic_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func decodeJSON(t *testing.T, s string) any {
	t.Helper()
	var v any
	require.NoError(t, json.Unmarshal([]byte(s), &v))
	return v
}

// modify changes every object and list inside a result.
func modify(v any) {
	switch v := v.(type) {
	case map[string]any:
		for _, item := range v {
			modify(item)
		}
		for key := range v {
			delete(v, key)
		}
		v["cat"] = "modified"
	case []any:
		for i := range v {
			modify(v[i])
			v[i] = "modified"
		}
	}
}

// assertResultIsNotShared checks that modifying a result changes neither a prepared rule
// nor the rule given to ApplyInterface.
func assertResultIsNotShared(t *testing.T, rule string, data string) bool {
	t.Helper()

	prepared, err := jsonlogic.Prepare(decodeJSON(t, rule))
	require.NoError(t, err)
	result, wantErr := prepared.Apply(decodeJSON(t, data))
	want := fmt.Sprintf("%#v", result) // it may hold NaN
	modify(result)
	result, err = prepared.Apply(decodeJSON(t, data))
	ok := assert.Equal(t, wantErr, err, rule)
	ok = assert.Equal(t, want, fmt.Sprintf("%#v", result), "prepared rule %s changed", rule) && ok

	logic := decodeJSON(t, rule)
	result, _ = jsonlogic.ApplyInterface(logic, decodeJSON(t, data))
	modify(result)
	return assert.Equal(t, decodeJSON(t, rule), logic, "rule changed by ApplyInterface") && ok
}

func TestResultDoesNotShareMemoryWithTheRule(t *testing.T) {
	rules := []string{
		`{"a":1,"b":2}`,
		`{"if":[true,{"a":1,"b":2}]}`,
		`{"if":[true,[{"a":1,"b":2},1]]}`,
		`{"and":[true,{"a":1,"b":[1,{"c":2}]}]}`,
		`{"or":[false,{"a":1,"b":2}]}`,
		`[{"a":1,"b":2}]`,
		`{"merge":[[{"a":1,"b":2}],[3]]}`,
		`{"set":[{"a":1,"b":2},"c",{"d":1,"e":2}]}`,
		`{"var":["gone",{"a":1,"b":2}]}`,
		`{"map":[[1],{"a":1,"b":2}]}`,
		`{"map":[[{"a":1,"b":2}],{"var":""}]}`,
		`{"map":[[{"a":[1]}],{"var":"a"}]}`,
		`{"map":[[[1,2]],{"var":""}]}`,
		`{"filter":[[{"a":1,"b":2}],true]}`,
		`{"filter":[[{"a":1}],true]}`,
		`{"filter":[[{}],true]}`,
		`{"filter":[[[1,2]],true]}`,

		`{"filter":[[1,2],true]}`,
		`{"if":[true,[1,2]]}`,
		`{"if":[true,[]]}`,
		`{"merge":[[1],[2]]}`,
		`[1,2]`,
	}
	for _, rule := range rules {
		t.Run(rule, func(t *testing.T) {
			assertResultIsNotShared(t, rule, `{}`)
		})
	}
}

func TestRandomResultsDoNotShareMemoryWithTheRule(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	scalars := []any{"x", "a", "b", "", float64(1), float64(0), true, nil}
	ops := []string{
		"and", "or", "if", "!", "==", "in", "var", "+", "cat", "merge", "set",
		"missing", "filter", "map", "reduce", "all", "some", "none",
	}

	var value, rule func(depth int) any
	value = func(depth int) any {
		k := r.Intn(10)
		switch {
		case depth == 0 || k < 4:
			return scalars[r.Intn(len(scalars))]
		case k < 6:
			list := make([]any, r.Intn(4))
			for i := range list {
				list[i] = value(depth - 1)
			}
			return list
		case k < 9:
			return map[string]any{"a": value(depth - 1), "var": value(depth - 1)}
		default:
			return map[string]any{ops[r.Intn(len(ops))]: value(depth - 1)}
		}
	}
	rule = func(depth int) any {
		if depth == 0 || r.Intn(4) == 0 {
			return value(2)
		}
		args := make([]any, 2+r.Intn(2))
		for i := range args {
			args[i] = rule(depth - 1)
		}
		return map[string]any{ops[r.Intn(len(ops))]: args}
	}

	for i := 0; i < 5000; i++ {
		data := toJSON(map[string]any{"a": value(2), "b": value(2)})
		if !assertResultIsNotShared(t, toJSON(rule(3)), data) {
			return
		}
	}
}

func TestModifiedResultIsNotEvaluatedAsLogic(t *testing.T) {
	prepared, err := jsonlogic.Prepare(decodeJSON(t, `{"if":[true,{"name":"x","var":"a"}]}`))
	require.NoError(t, err)
	data := map[string]any{"a": "from data"}

	result, err := prepared.Apply(data)
	require.NoError(t, err)
	// A map with a single key is an operator call.
	delete(result.(map[string]any), "name")

	result, err = prepared.Apply(data)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "x", "var": "a"}, result)
}

func TestModifiedResultDoesNotBypassTheTypeCheck(t *testing.T) {
	prepared, err := jsonlogic.Prepare(decodeJSON(t, `{"if":[true,{"a":1,"b":2}]}`))
	require.NoError(t, err)

	result, err := prepared.Apply(nil)
	require.NoError(t, err)
	result.(map[string]any)["a"] = struct{}{}

	result, err = prepared.Apply(nil)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"a": float64(1), "b": float64(2)}, result)
}

func TestConcurrentResultsCanBeModified(t *testing.T) {
	prepared, err := jsonlogic.Prepare(decodeJSON(t, `{"if":[true,{"plan":"free","limit":1}]}`))
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(user int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				result, err := prepared.Apply(nil)
				if !assert.NoError(t, err) {
					return
				}
				result.(map[string]any)["user"] = user
			}
		}(i)
	}
	wg.Wait()
}
