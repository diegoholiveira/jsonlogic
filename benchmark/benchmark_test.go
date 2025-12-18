package benchmark

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

var TestCases = []struct {
	name  string
	logic string
	data  string
}{
	{
		name:  "baseline_noop",
		logic: `true`,
		data:  `{}`,
	},
	{
		name:  "simple_equal",
		logic: `{"==": [1, 1]}`,
		data:  `{}`,
	},
	{
		name:  "complex_condition",
		logic: `{"and": [{"<": [{"var": "temp"}, 110]}, {"==": [{"var": "pie.filling"}, "apple"]}]}`,
		data:  `{"temp": 100, "pie": {"filling": "apple"}}`,
	},
	{
		name:  "nested_var",
		logic: `{"var": ["deeply.nested.variable", 99]}`,
		data:  `{"deeply": {"nested": {"variable": 42}}}`,
	},
	{
		name:  "array_operations",
		logic: `{"map": [{"var": "integers"}, {"*": [{"var": ""}, 2]}]}`,
		data:  `{"integers": [1, 2, 3, 4, 5]}`,
	},
	{
		name: "complex_logic",
		logic: `{"if": [
            {"<": [{"var": "age"}, 18]},
            "Too young",
            {"and": [
                {"<": [{"var": "age"}, 65]},
                {">=": [{"var": "age"}, 18]}
            ]},
            "Adult",
            "Senior"
        ]}`,
		data: `{"age": 25}`,
	},
	{
		name:  "reduce_sum",
		logic: `{"reduce": [{"var": "numbers"}, {"+": [{"var": "accumulator"}, {"var": "current"}]}, 0]}`,
		data:  `{"numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`,
	},
	{
		name:  "filter_even_numbers",
		logic: `{"filter": [{"var": "numbers"}, {"==": [{"%": [{"var": ""}, 2]}, 0]}]}`,
		data:  `{"numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`,
	},
	{
		name:  "contains_all",
		logic: `{"contains_all": [{"var": "tags"}, ["urgent", "reviewed"]]}`,
		data:  `{"tags": ["urgent", "reviewed", "approved", "processed"]}`,
	},
	{
		name:  "contains_any",
		logic: `{"contains_any": [{"var": "permissions"}, ["admin", "superuser"]]}`,
		data:  `{"permissions": ["user", "editor", "admin"]}`,
	},
	{
		name:  "contains_none",
		logic: `{"contains_none": [{"var": "flags"}, ["banned", "suspended"]]}`,
		data:  `{"flags": ["active", "verified", "premium"]}`,
	},
	{
		name:  "all_validation",
		logic: `{"all": [{"var": "users"}, {">": [{"var": ".age"}, 18]}]}`,
		data:  `{"users": [{"age": 25}, {"age": 30}, {"age": 22}, {"age": 19}]}`,
	},
	{
		name:  "some_validation",
		logic: `{"some": [{"var": "items"}, {"<": [{"var": ".price"}, 100]}]}`,
		data:  `{"items": [{"price": 150}, {"price": 75}, {"price": 200}]}`,
	},
	{
		name:  "string_concatenation",
		logic: `{"cat": [{"var": "firstName"}, " ", {"var": "lastName"}]}`,
		data:  `{"firstName": "John", "lastName": "Doe"}`,
	},
	{
		name:  "substring_extraction",
		logic: `{"substr": [{"var": "text"}, 0, 10]}`,
		data:  `{"text": "The quick brown fox jumps over the lazy dog"}`,
	},
	{
		name:  "max_operation",
		logic: `{"max": [85, 92, 78, 95, 88]}`,
		data:  `{}`,
	},
	{
		name:  "min_operation",
		logic: `{"min": [19.99, 15.50, 22.00, 12.99]}`,
		data:  `{}`,
	},
	{
		name:  "modulo_operation",
		logic: `{"%": [{"var": "value"}, 3]}`,
		data:  `{"value": 17}`,
	},
	{
		name:  "or_operation",
		logic: `{"or": [{"<": [{"var": "age"}, 18]}, {">": [{"var": "age"}, 65]}]}`,
		data:  `{"age": 70}`,
	},
	{
		name:  "merge_arrays",
		logic: `{"merge": [{"var": "array1"}, {"var": "array2"}]}`,
		data:  `{"array1": [1, 2, 3], "array2": [4, 5, 6]}`,
	},
	{
		name:  "missing_fields",
		logic: `{"missing": ["name", "email", "phone"]}`,
		data:  `{"name": "John", "email": "john@example.com"}`,
	},
	{
		name:  "deeply_nested_operations",
		logic: `{"and": [{"filter": [{"var": "users"}, {">": [{"var": ".age"}, 18]}]}, {"map": [{"var": "items"}, {"*": [{"var": ".price"}, 1.1]}]}]}`,
		data:  `{"users": [{"age": 25}, {"age": 30}], "items": [{"price": 10}, {"price": 20}]}`,
	},
}

func performWarmupRuns() {
	runtime.GC()

	for _, tc := range TestCases {
		for i := 0; i < 10; i++ {
			logic := strings.NewReader(tc.logic)
			data := strings.NewReader(tc.data)
			var result bytes.Buffer
			_ = jsonlogic.Apply(logic, data, &result)
		}
	}

	runtime.GC()
}

// Helper function to reduce duplication in benchmarks
func runBenchmark(b *testing.B, logic, data string) {
	// Pre-convert to bytes to avoid string overhead in loop
	logicBytes := []byte(logic)
	dataBytes := []byte(data)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logicReader := bytes.NewReader(logicBytes)
		dataReader := bytes.NewReader(dataBytes)
		var result bytes.Buffer
		err := jsonlogic.Apply(logicReader, dataReader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkComprehensive runs a focused suite of complex, realistic benchmarks
// that exercise multiple operators and represent real-world usage patterns.
// This is the default benchmark suite for version comparisons.
func BenchmarkComprehensive(b *testing.B) {
	cases := []struct {
		name  string
		logic string
		data  string
	}{
		{
			name: "user_validation",
			logic: `{
				"and": [
					{">=": [{"var": "user.age"}, 18]},
					{"in": [{"var": "user.country"}, ["US", "CA", "UK", "AU"]]},
					{"or": [
						{"==": [{"var": "user.subscription"}, "premium"]},
						{"<": [{"+": [{"var": "user.loginCount"}, 1]}, 100]}
					]}
				]
			}`,
			data: `{"user": {"age": 25, "country": "US", "subscription": "premium", "loginCount": 50}}`,
		},
		{
			name: "data_pipeline",
			logic: `{
				"reduce": [
					{"filter": [
						{"var": "orders"},
						{">=": [{"var": ".amount"}, 100]}
					]},
					{"+": [{"var": "accumulator"}, {"var": "current.amount"}]},
					0
				]
			}`,
			data: `{"orders": [{"amount": 50}, {"amount": 150}, {"amount": 200}, {"amount": 75}, {"amount": 120}]}`,
		},
		{
			name: "business_rules",
			logic: `{
				"if": [
					{"and": [
						{">": [{"var": "order.total"}, 1000]},
						{"==": [{"var": "customer.tier"}, "gold"]}
					]},
					{"*": [{"var": "order.total"}, 0.8]},
					{">": [{"var": "order.total"}, 500]},
					{"*": [{"var": "order.total"}, 0.9]},
					{"var": "order.total"}
				]
			}`,
			data: `{"order": {"total": 1200}, "customer": {"tier": "gold"}}`,
		},
		{
			name: "array_validation",
			logic: `{
				"and": [
					{"all": [{"var": "items"}, {">": [{"var": ".quantity"}, 0]}]},
					{"some": [{"var": "items"}, {"<": [{"var": ".price"}, 50]}]},
					{"none": [{"var": "items"}, {"==": [{"var": ".status"}, "cancelled"]}]}
				]
			}`,
			data: `{"items": [{"quantity": 2, "price": 30, "status": "active"}, {"quantity": 1, "price": 75, "status": "active"}]}`,
		},
		{
			name: "string_processing",
			logic: `{
				"and": [
					{"in": ["error", {"var": "message"}]},
					{">": [{"var": "severity"}, 5]},
					{"==": [{"substr": [{"var": "code"}, 0, 3]}, "ERR"]}
				]
			}`,
			data: `{"message": "System error detected", "severity": 8, "code": "ERR-500"}`,
		},
		{
			name: "custom_operators",
			logic: `{
				"and": [
					{"contains_all": [{"var": "required_permissions"}, ["read", "write"]]},
					{"contains_any": [{"var": "user_roles"}, ["admin", "moderator"]]},
					{"contains_none": [{"var": "flags"}, ["banned", "suspended"]]}
				]
			}`,
			data: `{"required_permissions": ["read", "write", "execute"], "user_roles": ["admin", "user"], "flags": ["active", "verified"]}`,
		},
		{
			name: "complex_data_transform",
			logic: `{
				"map": [
					{"filter": [
						{"var": "products"},
						{"and": [
							{"in": [{"var": ".category"}, ["electronics", "accessories"]]},
							{">": [{"var": ".stock"}, 0]}
						]}
					]},
					{"*": [{"var": ".price"}, 1.1]}
				]
			}`,
			data: `{"products": [{"category": "electronics", "price": 100, "stock": 5}, {"category": "clothing", "price": 50, "stock": 10}, {"category": "accessories", "price": 25, "stock": 0}]}`,
		},
		{
			name: "nested_conditions",
			logic: `{
				"if": [
					{"and": [
						{"missing": ["name", "email"]},
						{">": [{"var": "age"}, 0]}
					]},
					{"cat": ["Missing required fields for user ", {"var": "id"}]},
					{"or": [
						{"<": [{"var": "age"}, 13]},
						{"missing_some": [1, ["parent_email", "guardian_name"]]}
					]},
					"Parental consent required",
					"Valid user"
				]
			}`,
			data: `{"name": "John", "email": "john@example.com", "age": 25, "id": "12345"}`,
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			runBenchmark(b, tc.logic, tc.data)
		})
	}
}

// BenchmarkDetailed runs all detailed benchmarks (65 total).
// Use this for comprehensive testing of individual operators.
// For version comparisons, use BenchmarkComprehensive instead.
func BenchmarkDetailed(b *testing.B) {
	performWarmupRuns()

	for _, tc := range TestCases {
		b.Run(tc.name, func(b *testing.B) {
			runBenchmark(b, tc.logic, tc.data)
		})
	}
}

// Parallel benchmarks for testing concurrent performance
func BenchmarkJSONLogicParallel(b *testing.B) {
	parallelCases := []struct {
		name  string
		logic string
		data  string
	}{
		{
			name:  "simple_equal",
			logic: `{"==": [1, 1]}`,
			data:  `{}`,
		},
		{
			name:  "map",
			logic: `{"map": [{"var": "integers"}, {"*": [{"var": ""}, 2]}]}`,
			data:  `{"integers": [1, 2, 3, 4, 5]}`,
		},
		{
			name:  "reduce",
			logic: `{"reduce": [{"var": "numbers"}, {"+": [{"var": "accumulator"}, {"var": "current"}]}, 0]}`,
			data:  `{"numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`,
		},
		{
			name:  "filter",
			logic: `{"filter": [{"var": "numbers"}, {"==": [{"%": [{"var": ""}, 2]}, 0]}]}`,
			data:  `{"numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`,
		},
	}

	for _, tc := range parallelCases {
		b.Run(tc.name, func(b *testing.B) {
			logicBytes := []byte(tc.logic)
			dataBytes := []byte(tc.data)

			b.ResetTimer()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					logic := bytes.NewReader(logicBytes)
					data := bytes.NewReader(dataBytes)
					var result bytes.Buffer
					err := jsonlogic.Apply(logic, data, &result)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// Size variation benchmarks to test scalability
func BenchmarkArrayOperationsScaling(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"small_10", 10},
		{"medium_100", 100},
		{"large_1000", 1000},
	}

	generateIntArray := func(size int) string {
		result := "["
		for i := 0; i < size; i++ {
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf("%d", i+1)
		}
		result += "]"
		return result
	}

	// Map operation scaling
	b.Run("map", func(b *testing.B) {
		for _, s := range sizes {
			b.Run(s.name, func(b *testing.B) {
				logic := `{"map": [{"var": "integers"}, {"*": [{"var": ""}, 2]}]}`
				data := fmt.Sprintf(`{"integers": %s}`, generateIntArray(s.size))
				runBenchmark(b, logic, data)
			})
		}
	})

	// Filter operation scaling
	b.Run("filter", func(b *testing.B) {
		for _, s := range sizes {
			b.Run(s.name, func(b *testing.B) {
				logic := `{"filter": [{"var": "numbers"}, {"==": [{"%": [{"var": ""}, 2]}, 0]}]}`
				data := fmt.Sprintf(`{"numbers": %s}`, generateIntArray(s.size))
				runBenchmark(b, logic, data)
			})
		}
	})

	// Reduce operation scaling
	b.Run("reduce", func(b *testing.B) {
		for _, s := range sizes {
			b.Run(s.name, func(b *testing.B) {
				logic := `{"reduce": [{"var": "numbers"}, {"+": [{"var": "accumulator"}, {"var": "current"}]}, 0]}`
				data := fmt.Sprintf(`{"numbers": %s}`, generateIntArray(s.size))
				runBenchmark(b, logic, data)
			})
		}
	})

	// All operation scaling
	b.Run("all", func(b *testing.B) {
		for _, s := range sizes {
			b.Run(s.name, func(b *testing.B) {
				logic := `{"all": [{"var": "numbers"}, {">": [{"var": ""}, 0]}]}`
				data := fmt.Sprintf(`{"numbers": %s}`, generateIntArray(s.size))
				runBenchmark(b, logic, data)
			})
		}
	})

	// Some operation scaling
	b.Run("some", func(b *testing.B) {
		for _, s := range sizes {
			b.Run(s.name, func(b *testing.B) {
				logic := `{"some": [{"var": "numbers"}, {">": [{"var": ""}, 500]}]}`
				data := fmt.Sprintf(`{"numbers": %s}`, generateIntArray(s.size))
				runBenchmark(b, logic, data)
			})
		}
	})

	// None operation scaling
	b.Run("none", func(b *testing.B) {
		for _, s := range sizes {
			b.Run(s.name, func(b *testing.B) {
				logic := `{"none": [{"var": "numbers"}, {"<": [{"var": ""}, 0]}]}`
				data := fmt.Sprintf(`{"numbers": %s}`, generateIntArray(s.size))
				runBenchmark(b, logic, data)
			})
		}
	})
}

// Categorical benchmarks for easier filtering
func BenchmarkMathOperations(b *testing.B) {
	cases := []struct {
		name  string
		logic string
		data  string
	}{
		{"add", `{"+": [5, 3]}`, `{}`},
		{"subtract", `{"-": [10, 3]}`, `{}`},
		{"multiply", `{"*": [4, 5]}`, `{}`},
		{"divide", `{"/": [20, 4]}`, `{}`},
		{"modulo", `{"%": [17, 3]}`, `{}`},
		{"max", `{"max": [85, 92, 78, 95, 88]}`, `{}`},
		{"min", `{"min": [19.99, 15.50, 22.00, 12.99]}`, `{}`},
		{"abs", `{"abs": [-42]}`, `{}`},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			runBenchmark(b, tc.logic, tc.data)
		})
	}
}

func BenchmarkStringOperations(b *testing.B) {
	cases := []struct {
		name  string
		logic string
		data  string
	}{
		{
			"concat",
			`{"cat": [{"var": "firstName"}, " ", {"var": "lastName"}]}`,
			`{"firstName": "John", "lastName": "Doe"}`,
		},
		{
			"substr",
			`{"substr": [{"var": "text"}, 0, 10]}`,
			`{"text": "The quick brown fox jumps over the lazy dog"}`,
		},
		{
			"in_string",
			`{"in": ["quick", {"var": "text"}]}`,
			`{"text": "The quick brown fox"}`,
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			runBenchmark(b, tc.logic, tc.data)
		})
	}
}

func BenchmarkLogicOperations(b *testing.B) {
	cases := []struct {
		name  string
		logic string
		data  string
	}{
		{
			"and",
			`{"and": [{"<": [{"var": "temp"}, 110]}, {"==": [{"var": "status"}, "ok"]}]}`,
			`{"temp": 100, "status": "ok"}`,
		},
		{
			"or",
			`{"or": [{"<": [{"var": "age"}, 18]}, {">": [{"var": "age"}, 65]}]}`,
			`{"age": 70}`,
		},
		{
			"not",
			`{"!": [false]}`,
			`{}`,
		},
		{
			"if",
			`{"if": [{"<": [{"var": "age"}, 18]}, "minor", "adult"]}`,
			`{"age": 25}`,
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			runBenchmark(b, tc.logic, tc.data)
		})
	}
}

func BenchmarkCustomOperators(b *testing.B) {
	cases := []struct {
		name  string
		logic string
		data  string
	}{
		{
			"contains_all",
			`{"contains_all": [{"var": "tags"}, ["urgent", "reviewed"]]}`,
			`{"tags": ["urgent", "reviewed", "approved", "processed"]}`,
		},
		{
			"contains_any",
			`{"contains_any": [{"var": "permissions"}, ["admin", "superuser"]]}`,
			`{"permissions": ["user", "editor", "admin"]}`,
		},
		{
			"contains_none",
			`{"contains_none": [{"var": "flags"}, ["banned", "suspended"]]}`,
			`{"flags": ["active", "verified", "premium"]}`,
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			runBenchmark(b, tc.logic, tc.data)
		})
	}
}
