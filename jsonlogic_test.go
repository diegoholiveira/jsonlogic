package jsonlogic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/diegoholiveira/jsonlogic/v3/internal"
	"github.com/stretchr/testify/assert"
)

func TestRulesFromJsonLogic(t *testing.T) {
	tests := internal.GetScenariosFromOfficialTestSuite()

	for i, test := range tests {
		t.Run(fmt.Sprintf("Scenario_%d", i), func(t *testing.T) {
			var result bytes.Buffer

			err := Apply(test.Rule, test.Data, &result)
			if err != nil {
				t.Fatal(err)
			}

			if b, err := ioutil.ReadAll(test.Expected); err == nil {
				assert.JSONEq(t, string(b), result.String())
			}
		})
	}
}

func TestDivWithOnlyOneValue(t *testing.T) {
	rule := strings.NewReader(`{"/":[4]}`)
	data := strings.NewReader(`null`)

	var result bytes.Buffer

	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `4`, result.String())
}

func TestSetAValue(t *testing.T) {
	rule := strings.NewReader(`{
		"map": [
			{"var": "objects"},
			{"set": [
				{"var": ""},
				"age",
				{"+": [{"var": ".age"}, 2]}
			]}
		]
	}`)

	data := strings.NewReader(`{
		"objects": [
			{"age": 100, "location": "north"},
			{"age": 500, "location": "south"}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
		{"age": 102, "location": "north"},
		{"age": 502, "location": "south"}
	]`

	assert.JSONEq(t, expected, result.String())
}

func TestLocalContext(t *testing.T) {
	rule := strings.NewReader(`{
		"filter": [
			{"var": "people"},
			{"==": [
				{"var": ".age"},
				{"min": {"map": [
					{"var": "people"},
					{"var": ".age"}
				]}}
			]}
		]
	}`)

	data := strings.NewReader(`{
		"people": [
			{"age":18, "name":"John"},
			{"age":20, "name":"Luke"},
			{"age":18, "name":"Mark"}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
		{"age": 18, "name": "John"},
		{"age": 18, "name": "Mark"}
	]`

	assert.JSONEq(t, expected, result.String())
}

func TestMapWithZeroValue(t *testing.T) {
	rule := strings.NewReader(`{
		"filter": [
			{"var": "people"},
			{"==": [
				{"var": ".age"},
				{"min": {"map": [
					{"var": "people"},
					{"var": ".age"}
				]}}
			]}
		]
	}`)

	data := strings.NewReader(`{
		"people": [
			{"age":0, "name":"John"}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
		{"age": 0, "name": "John"}
	]`

	assert.JSONEq(t, expected, result.String())
}

func TestListOfRanges(t *testing.T) {
	rule := strings.NewReader(`{
		"filter": [
			{"var": "people"},
			{"in": [
				{"var": ".age"},
				[
					[12, 18],
					[22, 28],
					[32, 38]
				]
			]}
		]
	}`)

	data := strings.NewReader(`{
		"people": [
			{"age":18, "name":"John"},
			{"age":20, "name":"Luke"},
			{"age":18, "name":"Mark"}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
		{"age": 18, "name": "John"},
		{"age": 18, "name": "Mark"}
	]`

	assert.JSONEq(t, expected, result.String())
}

func TestInSortedOperator(t *testing.T) {
	rule := strings.NewReader(`{
		"filter": [
			{"var": "people"},
			{"in_sorted": [
				{"var": ".age"},
				[
					11.00,
					[12, 14],
					[13, 18],
					2,
					"20",
					[32, 38],
					"a",
					["b", "d"]
				]
			]}
		]
	}`)

	data := strings.NewReader(`{
		"people": [
			{"age":"18", "name":"John"},
			{"age":20, "name":"Luke"},
			{"age":18, "name":"Mark"},
			{"age":40, "name":"Donald"},
			{"age":11, "name":"Mickey"},
			{"age":"1", "name":"Minnie"},
			{"age":2, "name":"Mario"},
			{"age":"a", "name":"Mario"},
			{"age":"c", "name":"Princess"}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
		{"age":"18", "name": "John"},
		{"age":20, "name":"Luke"},
		{"age":18, "name": "Mark"},
		{"age":11, "name":"Mickey"},
		{"age":2, "name":"Mario"},
		{"age":"a", "name":"Mario"},
		{"age":"c", "name":"Princess"}
	]`

	assert.JSONEq(t, expected, result.String())
}

func TestSomeWithLists(t *testing.T) {
	rule := strings.NewReader(`{
		"some": [
			[511, 521, 811],
			{"in":[
				{"var":""},
				[1, 2, 3, 511]
			]}
		]
	}`)

	data := strings.NewReader(`{}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, "true", result.String())
}

func TestAllWithLists(t *testing.T) {
	rule := strings.NewReader(`{
		"all": [
			[511, 521, 811],
			{"in":[
				{"var":""},
				[511, 521, 811, 3]
			]}
		]
	}`)

	data := strings.NewReader("{}")

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, "true", result.String())
}

func TestAllWithArrayOfMapData(t *testing.T) {
	data := strings.NewReader(`[
		{
		  "P1": "A",
		  "P2":"a"
		},

		{
		  "P1": "B",
		  "P2":"b"
		}
	  ]`)
	rule := strings.NewReader(`
	  {
		"all": [
		  { "var": "" },
		  { "in": [ {"var": "P1"} , ["A","B"]] }
		]
	  }
	`)
	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEq(t, "true", result.String())
}

func TestNoneWithLists(t *testing.T) {
	rule := strings.NewReader(`{
		"none": [
			[511, 521, 811],
			{"in":[
				{"var":""},
				[1, 2]
			]}
		]
	}`)

	data := strings.NewReader("{}")

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, "true", result.String())
}

func TestInOperatorWorksWithMaps(t *testing.T) {
	rule := strings.NewReader(`{
		"some": [
			[511,521,811],
			{"in": [
				{"var": ""},
				{"map": [
					{"var": "my_list"},
					{"var": ".service_id"}
				]}
			]}
		]
	}`)

	data := strings.NewReader(`{
		"my_list": [
			{"service_id": 511},
			{"service_id": 771},
			{"service_id": 521},
			{"service_id": 181}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, "true", result.String())
}

func TestAbsoluteValue(t *testing.T) {
	rule := strings.NewReader(`{
		"abs": { "var": "test.number" }
	}`)

	data := strings.NewReader(`{
		"test": {
			"number": -2
		}
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, "2", result.String())
}

func TestMergeArrayOfArrays(t *testing.T) {
	rule := strings.NewReader(`{
		"merge": [
			[
				[
					"18800000",
					"18800969"
				]
			],
			[
				[
					"19840000",
					"19840969"
				]
			]
		]
	}`)
	data := strings.NewReader(`{}`)

	expectedResult := "[[\"18800000\",\"18800969\"],[\"19840000\",\"19840969\"]]"

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, expectedResult, result.String())
}

func TestDataWithDefaultValueWithApplyRaw(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"+": [
			1,
			2
		]
	}`)

	var expected json.RawMessage = json.RawMessage("3")

	output, err := ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}

func TestDataWithDefaultValueWithApplyInterface(t *testing.T) {
	rule := map[string]interface{}{
		"+": []interface{}{
			float64(1),
			float64(2),
		},
	}

	expected := float64(3)
	output, err := ApplyInterface(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, output.(float64))
}

func TestMissingOperators(t *testing.T) {
	rule := map[string]interface{}{
		"sum": []interface{}{
			float64(1),
			float64(2),
		},
	}

	_, err := ApplyInterface(rule, nil)

	assert.EqualError(t, err, "The operator \"sum\" is not supported")
}

func TestZeroDivision(t *testing.T) {
	logic := strings.NewReader(`{"/":[0,10]}`)
	data := strings.NewReader(`{}`)
	var result bytes.Buffer

	Apply(logic, data, &result) // nolint:errcheck

	assert.JSONEq(t, `0`, result.String())
}

func TestSliceWithOnlyWithNumbersAsKey(t *testing.T) {
	rule := strings.NewReader(`{"var": "people.0"}`)

	data := strings.NewReader(`{
		"people": [
			{"age":18, "name":"John"},
			{"age":20, "name":"Luke"},
			{"age":18, "name":"Mark"}
		]
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"age": 18, "name": "John"}`

	assert.JSONEq(t, expected, result.String())
}

func TestMapWithOnlyWithNumbersAsKey(t *testing.T) {
	rule := strings.NewReader(`{"var": "people.103"}`)

	data := strings.NewReader(`{
		"people": {
			"100": {"age":18, "name":"John"},
			"101": {"age":20, "name":"Luke"},
			"103": {"age":18, "name":"Mark"}
		}
	}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"age": 18, "name": "Mark"}`

	assert.JSONEq(t, expected, result.String())
}

func TestBetweenIsBiggerEq(t *testing.T) {
	rule := strings.NewReader(`{
		"filter": [
			[1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
			{">=": [8, {"var": ""}, 3]}
		]
	}`)

	data := strings.NewReader(`{}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[3, 4, 5, 6, 7, 8]`

	assert.JSONEq(t, expected, result.String())
}

func TestBetweenIsBigger(t *testing.T) {
	rule := strings.NewReader(`{
		"filter": [
			[1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
			{">": [8, {"var": ""}, 3]}
		]
	}`)

	data := strings.NewReader(`{}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[4, 5, 6, 7]`

	assert.JSONEq(t, expected, result.String())
}

func TestUnaryOperation(t *testing.T) {
	logic := strings.NewReader(`{"and":[{"!":{"var":"var_not_in_data"}}]}`)
	data := strings.NewReader(`{"some_key": "value"}`)

	var result bytes.Buffer
	assert.Nil(t, Apply(logic, data, &result))

	assert.JSONEq(t, `true`, result.String())
}

func TestInOperatorAgainstNil(t *testing.T) {
	rule := strings.NewReader(`{"filter":[{"var": "accounts"},{"and":[{"in":["abc",{"var":"tags.tag-1"}]}]}]}`)
	data := strings.NewReader(`{"accounts":[{"name":"account-1","tags":{"tag-1":"abc"}}, {"name":"account-2","tags":{"tag-2":"xyz"}}]}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[
		{
			"name": "account-1",
			"tags": {
				"tag-1": "abc"
			}
		}
	]`

	assert.JSONEq(t, expected, result.String())
}

func TestReduceFilterAndContains(t *testing.T) {
	rule := strings.NewReader(`{"reduce":[{"filter":[{"var":"data.level1.level2"},{"==":[{"var":"access"},true]}]},{"or":[{"var":"current.access"},{"var":"accumulator"}]},false]}`)
	data := strings.NewReader(`{"data":{"level1":{"level2":[{"access":true }]}}}}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `true`

	assert.JSONEq(t, expected, result.String())
}

func TestReduceFilterAndNotContains(t *testing.T) {
	rule := strings.NewReader(`{"reduce":[{"filter":[{"var":"data.level1.level2"},{"==":[{"var":"access"},true]}]},{"or":[{"var":"current.access"},{"var":"accumulator"}]},false]}`)
	data := strings.NewReader(`{"data":{"level1":{"level2":[{"access":false }]}}}}`)

	var result bytes.Buffer
	err := Apply(rule, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `false`

	assert.JSONEq(t, expected, result.String())
}

func TestReduceWithUnsupportedValue(t *testing.T) {
	b := []byte(`{"reduce":[{"filter":[{"var":"data"},{"==":[{"var":""},""]}]},{"cat":[{"var":"current"},{"var":"accumulator"}]},null]}`)

	rule := map[string]interface{}{}
	_ = json.Unmarshal(b, &rule)
	data := map[string]interface{}{
		"data": []interface{}{"str"},
	}

	_, err := ApplyInterface(rule, data)
	assert.EqualError(t, err, "The type \"<nil>\" is not supported")
}

func TestAddOperator(t *testing.T) {
	AddOperator("strlen", func(values, data interface{}) interface{} {
		v, ok := values.(string)

		if ok {
			return len(v)
		}
		return 0
	})
	logic := strings.NewReader(`{ "strlen": { "var": "foo" } }`)
	data := strings.NewReader(`{"foo": "bar"}`)

	var result bytes.Buffer
	err := Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `3`

	assert.JSONEq(t, expected, result.String())
}

func TestIssue50(t *testing.T) {
	logic := strings.NewReader(`{"<": ["abc", 3]}`)
	data := strings.NewReader(`{}`)

	var result bytes.Buffer
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	err := Apply(logic, data, &result)
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
	_ = Apply(logic, data, &result)
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

	output, err := GetJsonLogicWithSolvedVars(rule, data)

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

	err := Apply(rule, data, &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := `true`
	assert.JSONEq(t, expected, result.String())
}
