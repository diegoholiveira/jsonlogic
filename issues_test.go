package jsonlogic_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestIssue50(t *testing.T) {
	logic := strings.NewReader(`{"<": ["abc", 3]}`)
	data := strings.NewReader(`{}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `false`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue51_example1(t *testing.T) {
	logic := strings.NewReader(`{"==":[{"var":"test"},true]}`)
	data := strings.NewReader(`{}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `false`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue51_example2(t *testing.T) {
	logic := strings.NewReader(`{"==":[{"var":"test"},"true"]}`)
	data := strings.NewReader(`{"test": true}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `false`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue52_example1(t *testing.T) {
	data := strings.NewReader(`{}`)
	logic := strings.NewReader(`{"substr": ["jsonlogic", -10]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `"jsonlogic"`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue52_example2(t *testing.T) {
	data := strings.NewReader(`{}`)
	logic := strings.NewReader(`{"substr": ["jsonlogic", 10]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `""`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue58_example(t *testing.T) {
	data := strings.NewReader(`{"foo": "bar"}`)
	logic := strings.NewReader(`{"if":[
		{"==":[{"var":"foo"},"bar"]},{"foo":"is_bar","path":"foo_is_bar"},
		{"foo":"not_bar","path":"default_object"}
	]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"foo":"is_bar","path":"foo_is_bar"}`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue70(t *testing.T) {
	data := strings.NewReader(`{"people": [
		{"age":18, "name":"John"},
		{"age":20, "name":"Luke"},
		{"age":18, "name":"Mark"}
]}`)
	logic := strings.NewReader(`{"filter": [
	{"var": ["people"]},
	{"==": [{"var": ["age"]}, 18]}
]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
    {"age": 18, "name": "John"},
    {"age": 18, "name": "Mark"}
]`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue71_example_empty_min(t *testing.T) {
	data := strings.NewReader(`{}`)
	logic := strings.NewReader(`{"min":[]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `null`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue71_example_empty_max(t *testing.T) {
	data := strings.NewReader(`{}`)
	logic := strings.NewReader(`{"max":[]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `null`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue71_example_max(t *testing.T) {
	data := strings.NewReader(`{}`)
	logic := strings.NewReader(`{"max":[-3, -2]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `-2`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue74(t *testing.T) {
	logic := strings.NewReader(`{"if":[ false, {"var":"values.0.categories"}, "else" ]}`)
	data := strings.NewReader(`{ "values": [] }`)

	var result bytes.Buffer
	_ = jsonlogic.Apply(logic, data, &result)
	expected := `"else"`
	assert.JSONEq(t, expected, result.String())
}

func TestJsonLogicWithSolvedVars(t *testing.T) {
	rule := json.RawMessage(`{
		"or":[
		{
			"and":[
				{"==": [{ "var":"is_foo" }, true ]},
				{"==": [{ "var":"is_bar" }, true ]},
				{">=": [{ "var":"foo" }, 17179869184 ]},
				{"==": [{ "var":"bar" }, 0 ]}
			]
      	},
      	{
			"and":[
				{"==": [{ "var":"is_bar" }, true ]},
				{"==": [{ "var":"is_foo" }, false ]},
				{"==": [{ "var":"foo" }, 34359738368 ]},
				{"==": [{ "var":"bar" }, 0 ]}
			]
      	}]
    }`)

	data := json.RawMessage(`{"foo": 34359738368, "bar": 10, "is_foo": false, "is_bar": true}`)

	output, err := jsonlogic.GetJsonLogicWithSolvedVars(rule, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{
		"or":[
		{
			"and":[
				{ "==":[ false, true ] },
				{ "==":[ true, true ] },
				{ ">=":[ 34359738368, 17179869184 ] },
				{ "==":[ 10, 0 ] }
			]
		},
		{
			"and":[
				{ "==":[ true, true ] },
				{ "==":[ false, false ] },
				{ "==":[ 34359738368, 34359738368 ] },
				{ "==":[ 10, 0 ] }
			]
		}]
	}`

	assert.JSONEq(t, expected, string(output))
}

func TestIssue79(t *testing.T) {
	rule := strings.NewReader(
		`{"and": [
        {"in": [
          {"var": "flow"},
          ["BRAND"]
        ]},
        {"or": [
          {"if": [
            {"missing": ["gender"]},
            true,
            false
          ]},
          {"some": [
            {"var": "gender"},
            {"==": [
              {"var": null},
              "men"
            ]}
          ]}
        ]}
      ]}`,
	)

	data := strings.NewReader(`{"category":["sneakers"],"flow":"BRAND","gender":["men"],"market":"US"}`)

	var result bytes.Buffer

	err := jsonlogic.Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `true`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue83(t *testing.T) {
	rule := `{
	  "map": [
	    {"var": "listOfLists"},
	    {"in": ["item_a", {"var": ""}]}
	  ]
	}`

	data := `{
	  "listOfLists": [
	    ["item_a", "item_b", "item_c"],
	    ["item_b", "item_c"],
	    ["item_a", "item_c"]
	  ]
	}`

	var result bytes.Buffer

	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(data), &result)

	if assert.Nil(t, err) {
		expected := `[true,false,true]`
		assert.JSONEq(t, expected, result.String())
	}
}

func TestIssue81(t *testing.T) {
	rule := `{
      "some": [
        {"var": "A"},
        {"!=": [
          {"var": ".B"},
          {"var": "B"}
        ]}
      ]}
         `

	data := `{"A":[{"B":1}], "B":2}`

	var result bytes.Buffer

	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(data), &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `true`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue96(t *testing.T) {
	rule := `{"map":[
      {"var":"integers"},
	  {"*":[{"var":[""]},2]}
    ]}`

	data := `{"integers": [1,2,3]}`

	var result bytes.Buffer

	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(data), &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[2, 4, 6]`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue98(t *testing.T) {
	rule := `{"or": [{"and": [true]}]}`
	data := `{}`

	var result bytes.Buffer

	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(data), &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `true`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue110(t *testing.T) {
	logic := strings.NewReader(`{ "map":[{"var": "arr"},{"var":["xxx", "default"]}]}`)
	data := strings.NewReader(`{"arr": [{"xxx": "111","yyy": "222"},{"xxx": "333","yyy": "444"}]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `["111","333"]`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue125_InOperatorWithVarsInSlice(t *testing.T) {
	// This test demonstrates the issue: vars within slices are not resolved
	rule := strings.NewReader(`{"in": [{"var": "needle"}, [{"var": "item1"}, {"var": "item2"}]]}`)
	data := strings.NewReader(`{"needle":"foo", "item1":"bar", "item2":"foo"}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	// Should be true because "foo" should be found in the resolved array ["bar", "foo"]
	// Currently fails because it compares "foo" against unresolved [{"var": "item1"}, {"var": "item2"}]
	expected := `true`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue125_CustomOperatorWithVarsInSlice(t *testing.T) {
	// Add a custom operator that processes slice elements
	jsonlogic.AddOperator("contains_any", func(values, data any) any {
		parsed := values.([]any)
		needle := parsed[0]
		haystack := parsed[1].([]any)

		for _, item := range haystack {
			if item == needle {
				return true
			}
		}
		return false
	})

	rule := strings.NewReader(`{"contains_any": [{"var": "needle"}, [{"var": "item1"}, {"var": "item2"}]]}`)
	data := strings.NewReader(`{"needle":"foo", "item1":"bar", "item2":"foo"}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	// Should be true because "foo" should be found in the resolved array ["bar", "foo"]
	// Currently fails because the custom operator receives unresolved [{"var": "item1"}, {"var": "item2"}]
	expected := `true`
	assert.JSONEq(t, expected, result.String())
}

func TestIssue135(t *testing.T) {
	cases := []struct {
		name     string
		rule     string
		expected string
	}{
		{
			name:     "or returns last operand when all are falsy",
			rule:     `{"or":[null,0]}`,
			expected: `0`,
		},
		{
			name:     "and returns last operand when all are truthy",
			rule:     `{"and":[1,"result"]}`,
			expected: `"result"`,
		},
		{
			name:     "and returns last truthy operand, not max",
			rule:     `{"and":[3,1]}`,
			expected: `1`,
		},
		{
			name:     "and example from jsonlogic.com docs",
			rule:     `{"and":[true,"a",3]}`,
			expected: `3`,
		},
		{
			name:     "and returns first falsy operand (empty string)",
			rule:     `{"and":[true,"",3]}`,
			expected: `""`,
		},
		{
			name:     "or short-circuits on first truthy operand",
			rule:     `{"or":[1,0]}`,
			expected: `1`,
		},
		{
			name:     "and with non-empty slice continues past it",
			rule:     `{"and":[[1,2,3],true]}`,
			expected: `true`,
		},
		{
			name:     "and with empty slice returns it",
			rule:     `{"and":[[],true]}`,
			expected: `[]`,
		},
		{
			name:     "or returns null when all operands are null",
			rule:     `{"or":[null,null]}`,
			expected: `null`,
		},
		{
			name:     "or returns empty array as last falsy operand",
			rule:     `{"or":[false,[]]}`,
			expected: `[]`,
		},
		{
			name:     "and with empty operand list returns null",
			rule:     `{"and":[]}`,
			expected: `null`,
		},
		{
			name:     "or with empty operand list returns null",
			rule:     `{"or":[]}`,
			expected: `null`,
		},
		{
			name:     "or with single falsy operand returns it",
			rule:     `{"or":[0]}`,
			expected: `0`,
		},
		{
			name:     "and with single falsy operand returns it",
			rule:     `{"and":[0]}`,
			expected: `0`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var result bytes.Buffer
			err := jsonlogic.Apply(strings.NewReader(tc.rule), strings.NewReader(`{}`), &result)
			if err != nil {
				t.Fatal(err)
			}
			// `or` should return the first truthy operand or the last operand;
			// `and` should return the first falsy operand or the last operand.
			assert.JSONEq(t, tc.expected, result.String())
		})
	}
}

func TestIssue153_CustomOperatorNumericTypes(t *testing.T) {
	numbers := map[string]any{
		"int": int(3), "int8": int8(3), "int16": int16(3), "int32": int32(3), "int64": int64(3),
		"uint": uint(3), "uint8": uint8(3), "uint16": uint16(3), "uint32": uint32(3), "uint64": uint64(3),
		"float32": float32(3), "float64": float64(3),
	}

	// OP is replaced by a call to the custom operator, which returns 3.
	cases := []struct {
		rule     string
		expected any
	}{
		{`OP`, float64(3)},
		{`{">":[OP,0]}`, true},
		{`{"<":[OP,10]}`, true},
		{`{">=":[OP,3]}`, true},
		{`{"==":[OP,3]}`, true},
		{`{"===":[OP,3]}`, true},
		{`{"!==":[OP,3]}`, false},
		{`{"+":[OP,1]}`, float64(4)},
		{`{"-":[OP,1]}`, float64(2)},
		{`{"*":[OP,2]}`, float64(6)},
		{`{"/":[OP,2]}`, float64(1.5)},
		{`{"%":[OP,2]}`, float64(1)},
		{`{"max":[OP,1]}`, float64(3)},
		{`{"min":[OP,5]}`, float64(3)},
		{`{"abs":[OP]}`, float64(3)},
		{`{"if":[OP,"yes","no"]}`, "yes"},
		{`{"!":[OP]}`, false},
		{`{"!!":[OP]}`, true},
		{`{"and":[OP,"x"]}`, "x"},
		{`{"or":[OP,"x"]}`, float64(3)},
		{`{"in":[OP,[1,2,3]]}`, true},
		{`{"in":[3,[OP]]}`, true},
		{`{"cat":["a",OP]}`, "a3"},
		{`{"substr":["abcdef",OP]}`, "def"},
		{`{"merge":OP}`, []any{float64(3)}},
		{`{"reduce":[[1,2],{"+":[{"var":"current"},{"var":"accumulator"}]},OP]}`, float64(6)},
	}

	for name, number := range numbers {
		operator := "issue153_" + name
		jsonlogic.AddOperator(operator, func(values, data any) any { return number })

		for _, tt := range cases {
			t.Run(name+" "+tt.rule, func(t *testing.T) {
				var rule any
				err := json.Unmarshal([]byte(strings.ReplaceAll(tt.rule, "OP", `{"`+operator+`":[]}`)), &rule)
				assert.NoError(t, err)

				output, err := jsonlogic.ApplyInterface(rule, nil)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, output)

				prepared, err := jsonlogic.Prepare(rule)
				assert.NoError(t, err)
				output, err = prepared.Apply(nil)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, output)
			})
		}
	}
}

func TestIssue153_CustomOperatorContainers(t *testing.T) {
	list := []any{1, int64(2), 3}
	object := map[string]any{"a": 1, "nested": []any{map[string]any{"b": uint8(2)}}, "s": "x", "n": nil}
	jsonlogic.AddOperator("issue153_list", func(values, data any) any { return list })
	jsonlogic.AddOperator("issue153_object", func(values, data any) any { return object })

	cases := []struct {
		rule     string
		expected any
	}{
		{`{"issue153_list":[]}`, []any{float64(1), float64(2), float64(3)}},
		{`{"map":[{"issue153_list":[]},{"+":[{"var":""},1]}]}`, []any{float64(2), float64(3), float64(4)}},
		{`{"filter":[{"issue153_list":[]},{">":[{"var":""},1]}]}`, []any{float64(2), float64(3)}},
		{`{"some":[{"issue153_list":[]},{">":[{"var":""},2]}]}`, true},
		{`{"all":[{"issue153_list":[]},{">":[{"var":""},0]}]}`, true},
		{`{"none":[{"issue153_list":[]},{">":[{"var":""},2]}]}`, false},
		{`{"contains_all":[{"issue153_list":[]},[1,3]]}`, true},
		{`{"issue153_object":[]}`, map[string]any{
			"a":      float64(1),
			"nested": []any{map[string]any{"b": float64(2)}},
			"s":      "x",
			"n":      nil,
		}},
	}

	for _, tt := range cases {
		t.Run(tt.rule, func(t *testing.T) {
			var rule any
			assert.NoError(t, json.Unmarshal([]byte(tt.rule), &rule))

			output, err := jsonlogic.ApplyInterface(rule, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}

	// Values owned by the operator must not be modified.
	assert.Equal(t, []any{1, int64(2), 3}, list)
	assert.Equal(t, map[string]any{"a": 1, "nested": []any{map[string]any{"b": uint8(2)}}, "s": "x", "n": nil}, object)
}

func TestIssue153_CustomOperatorOtherTypesAreUnchanged(t *testing.T) {
	type point struct{ X int }

	for _, value := range []any{nil, true, "x", []int{1, 2}, point{X: 1}, json.Number("1")} {
		jsonlogic.AddOperator("issue153_other", func(values, data any) any { return value })

		output, err := jsonlogic.ApplyInterface(map[string]any{"issue153_other": []any{}}, nil)
		assert.NoError(t, err)
		assert.Equal(t, value, output)
	}
}
