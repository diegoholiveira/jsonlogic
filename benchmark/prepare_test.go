package benchmark

import (
	"encoding/json"
	"fmt"
	"testing"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

type prepareCase struct {
	name  string
	logic string
	data  string
}

func prepareCases(b *testing.B) []prepareCase {
	b.Helper()

	ids := make([]any, 1000)
	for i := range ids {
		ids[i] = fmt.Sprintf("user-%d", i)
	}
	largeIn, err := json.Marshal(map[string]any{"in": []any{map[string]any{"var": "user.id"}, ids}})
	if err != nil {
		b.Fatal(err)
	}

	var cases []prepareCase
	for _, tc := range comprehensiveCases {
		cases = append(cases, prepareCase{tc.name, tc.logic, tc.data})
	}

	return append(cases,
		prepareCase{
			name:  "feature_flag",
			logic: `{"and": [{"==": [{"var": "env"}, "prod"]}, {"in": [{"var": "id"}, ["290636", "290316"]]}]}`,
			data:  `{"env": "prod", "id": "290636", "dsp": "25", "key": "req-1"}`,
		},
		prepareCase{
			name: "feature_flag_targeting",
			logic: `{
				"if": [
					{"or": [
						{"in": [{"var": "user.id"}, ["u-1", "u-2", "u-3", "u-4", "u-5", "u-6", "u-7", "u-8"]]},
						{"and": [
							{"==": [{"var": "env"}, "staging"]},
							{"in": ["@example.com", {"var": "user.email"}]}
						]}
					]},
					"internal",
					{"and": [
						{"!": {"in": [{"var": "user.country"}, ["CN", "RU", "IR", "KP"]]}},
						{"!==": [{"var": "user.plan"}, "free"]},
						{">=": [{"var": "app.version"}, 42]},
						{"<": [{"%": [{"var": "user.bucket"}, 100]}, 25]}
					]},
					"treatment",
					"control"
				]
			}`,
			data: `{"env": "prod", "user": {"id": "u-99", "email": "ana@mail.test", "country": "BR", "plan": "pro", "bucket": 1318}, "app": {"version": 57}}`,
		},
		prepareCase{
			name: "object_result",
			logic: `{
				"if": [
					{"==": [{"var": "env"}, "prod"]},
					{"plan": "pro", "limits": {"a": 1, "b": 2}},
					{"plan": "free", "limits": {"a": 0, "b": 0}}
				]
			}`,
			data: `{"env": "prod"}`,
		},
		prepareCase{
			name:  "literal_objects",
			logic: `{"map": [[{"p": 1}, {"p": 2}, {"p": 3}, {"p": 4}], {"*": [{"var": "p"}, 2]}]}`,
			data:  `{}`,
		},
		prepareCase{
			name:  "in_large_list",
			logic: string(largeIn),
			data:  `{"user": {"id": "user-999"}}`,
		},
	)
}

// BenchmarkPrepare compares the ways of applying the same rule. Prepare is the cost
// paid once for PreparedApply.
func BenchmarkPrepare(b *testing.B) {
	for _, tc := range prepareCases(b) {
		b.Run(tc.name, func(b *testing.B) {
			var rule, data any
			if err := json.Unmarshal([]byte(tc.logic), &rule); err != nil {
				b.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.data), &data); err != nil {
				b.Fatal(err)
			}
			prepared, err := jsonlogic.Prepare(rule)
			if err != nil {
				b.Fatal(err)
			}

			applyInterface := func() (any, error) { return jsonlogic.ApplyInterface(rule, data) }
			preparedApply := func() (any, error) { return prepared.Apply(data) }

			b.Run("Apply", func(b *testing.B) { runBenchmark(b, tc.logic, tc.data) })
			b.Run("ApplyInterface", func(b *testing.B) { runSerial(b, applyInterface) })
			b.Run("PreparedApply", func(b *testing.B) { runSerial(b, preparedApply) })
			b.Run("ApplyInterfaceParallel", func(b *testing.B) { runParallel(b, applyInterface) })
			b.Run("PreparedApplyParallel", func(b *testing.B) { runParallel(b, preparedApply) })
			b.Run("Prepare", func(b *testing.B) {
				runSerial(b, func() (any, error) { return jsonlogic.Prepare(rule) })
			})
		})
	}
}

func runSerial(b *testing.B, run func() (any, error)) {
	b.Helper()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := run(); err != nil {
			b.Fatal(err)
		}
	}
}

func runParallel(b *testing.B, run func() (any, error)) {
	b.Helper()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := run(); err != nil {
				b.Error(err)
				return
			}
		}
	})
}
