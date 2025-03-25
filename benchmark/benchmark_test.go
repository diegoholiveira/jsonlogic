package benchmark

import (
	"bytes"
	"fmt"
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
}

func BenchmarkJSONLogic(b *testing.B) {
	for _, tc := range TestCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				logic := strings.NewReader(tc.logic)
				data := strings.NewReader(tc.data)
				var result bytes.Buffer
				err := jsonlogic.Apply(logic, data, &result)
				if err != nil {
					fmt.Printf("\n\nError: %+v\n\n", err)
					b.Fatal(err)
				}
			}
		})
	}
}
