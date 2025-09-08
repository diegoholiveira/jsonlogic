package jsonlogic_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/jsonlogic/v3"
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

	expected := `"jsonlogic"`
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
