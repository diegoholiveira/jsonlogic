package jsonlogic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestAlwaysShouldAlwaysPass(t *testing.T) {
	var result bool
	err := Apply(true, nil, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("Always should always pass")
	}
}

func TestNeverShouldNeverPass(t *testing.T) {
	var result bool
	err := Apply(false, nil, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatal("Always should never pass")
	}
}

func TestSimpleComparisonWithInteger(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte(`{
		"==": [1, 1]
	}`), &rules)

	var result bool
	err := Apply(rules, nil, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("A simple comparison is expected to be true")
	}
}

func TestSimpleComparisonWithString(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte(`{
		"==": ["a", "a"]
	}`), &rules)

	var result bool
	err := Apply(rules, nil, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("A simple comparison is expected to be true")
	}
}

func TestComposedComparisons(t *testing.T) {
	var rules interface{}

	json.Unmarshal([]byte(`{
		"and": [
			{"==": [1,1]},
			{"==": [1,2]}
		]
	}`), &rules)

	var result bool
	err := Apply(rules, nil, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatal("The composed comparison is expected to be false")
	}
}

func TestSimpleVar(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"var": "a"
	}`), &rules)

	json.Unmarshal([]byte(`{
		"a": 10
	}`), &data)

	var result float64
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != 10 {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestSimpleVarWithoutSyntacticSugar(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"var": ["a"]
	}`), &rules)

	json.Unmarshal([]byte(`{
		"a": 10
	}`), &data)

	var result float64
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != 10 {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestVariableWithDefaultValue(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"var": ["z", 20]
	}`), &rules)

	json.Unmarshal([]byte(`{
		"a": 10
	}`), &data)

	var result float64
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != 20 {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestSimpleVarComparison(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"==": [
			{"var": "a"},
			10
		]
	}`), &rules)

	json.Unmarshal([]byte(`{
		"a": 10
	}`), &data)

	var result bool
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestComposedVar(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"var": "champ.name"
	}`), &rules)

	json.Unmarshal([]byte(`{
		"champ": {
			"name": "Diego"
		}
	}`), &data)

	var result string
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "Diego" {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestIndexedVar(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"var": 1
	}`), &rules)

	json.Unmarshal([]byte(`[
		"apple",
		"banana",
		"carrot"
	]`), &data)

	var result string
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "banana" {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestComplexRule(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"and": [
			{"<": [{"var": "temp"}, 110]},
			{"==": [{"var": "pie.filling"}, "apple"]}
		]
	}`), &rules)

	json.Unmarshal([]byte(`{
		"temp": 100,
		"pie": {
			"filling": "apple"
		}
	}`), &data)

	var result bool
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestRulesFromJsonLogic(t *testing.T) {
	response, err := http.Get("http://jsonlogic.com/tests.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	buffer, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	var scenarios []interface{}

	err = json.Unmarshal(buffer, &scenarios)
	if err != nil {
		log.Println(err)
		return
	}

	for _, scenario := range scenarios {
		if reflect.ValueOf(scenario).Kind() == reflect.String {
			continue
		}

		validateScenario(t, scenario)
	}
}

func validateScenario(t *testing.T, scenario interface{}) {
	var result interface{}

	logic := scenario.([]interface{})[0]
	data := scenario.([]interface{})[1]
	expected := scenario.([]interface{})[2]

	err := Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, result) {
		fmt.Println("Logic ", logic)
		fmt.Println("Data ", data)
		fmt.Println("Expected ", fmt.Sprintf("%v %T", expected, expected))
		fmt.Println("Result ", fmt.Sprintf("%v %T", result, result))

		t.Fatal("The value expected is not what we expected")
	}
}

func TestSetAValue(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"map": [
			{"var": "objects"},
			{"set": [
				{"var": ""},
				"age",
				{"+": [{"var": ".age"}, 2]}
			]}
		]
	}`), &rules)

	json.Unmarshal([]byte(`{
		"objects": [
			{"age": 100, "location": "north"},
			{"age": 500, "location": "south"}
		]
	}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`[
		{"age": 102, "location": "north"},
		{"age": 502, "location": "south"}
	]`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("We expect a new object with new data")
	}
}

func TestLocalContext(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
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
	}`), &rules)

	json.Unmarshal([]byte(`{
		"people": [
			{"age":18, "name":"John"},
			{"age":20, "name":"Luke"},
			{"age":18, "name":"Mark"}
		]
	}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`[
		{"age": 18, "name": "John"},
		{"age": 18, "name": "Mark"}
	]`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("filter do not have access to all scope")
	}
}

func TestListOfRanges(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
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
	}`), &rules)

	json.Unmarshal([]byte(`{
		"people": [
			{"age":18, "name":"John"},
			{"age":20, "name":"Luke"},
			{"age":18, "name":"Mark"}
		]
	}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`[
		{"age": 18, "name": "John"},
		{"age": 18, "name": "Mark"}
	]`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("filter do not have access to all scope")
	}
}

func TestSomeWithLists(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"some": [
			[511, 521, 811],
			{"in":[
				{"var":""},
				[1, 2, 3, 511]
			]}
		]
	}`), &rules)

	json.Unmarshal([]byte(`{}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`true`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("list must work with some")
	}
}

func TestAllWithLists(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"all": [
			[511, 521, 811],
			{"in":[
				{"var":""},
				[511, 521, 811, 3]
			]}
		]
	}`), &rules)

	json.Unmarshal([]byte(`{}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`true`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("list must work with some")
	}
}

func TestNoneWithLists(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
		"none": [
			[511, 521, 811],
			{"in":[
				{"var":""},
				[1, 2]
			]}
		]
	}`), &rules)

	json.Unmarshal([]byte(`{}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`true`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("list must work with some")
	}
}

func TestInOperatorWorksWithMaps(t *testing.T) {
	var rules interface{}
	var data interface{}

	json.Unmarshal([]byte(`{
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
	}`), &rules)

	json.Unmarshal([]byte(`{
		"my_list": [
			{"service_id": 511},
			{"service_id": 771},
			{"service_id": 521},
			{"service_id": 181}
		]
	}`), &data)

	var result interface{}
	err := Apply(rules, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	var expected interface{}
	json.Unmarshal([]byte(`true`), &expected)

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("maps must work with in operator")
	}
}
