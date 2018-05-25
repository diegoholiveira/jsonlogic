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
	result, _ := BoolApply(true, nil)
	if !result {
		t.Fatal("Always should always pass")
	}
}

func TestNeverShouldNeverPass(t *testing.T) {
	result, _ := BoolApply(false, nil)
	if result {
		t.Fatal("Always should never pass")
	}
}

func TestSimpleComparisonWithInteger(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte(`{
		"==": [1, 1]
	}`), &rules)

	result, _ := BoolApply(rules, nil)
	if !result {
		t.Fatal("A simple comparison is expected to be true")
	}
}

func TestSimpleComparisonWithString(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte(`{
		"==": ["a", "a"]
	}`), &rules)

	result, _ := BoolApply(rules, nil)
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

	result, _ := BoolApply(rules, nil)
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

	result, _ := IntApply(rules, interface{}(data))
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

	result, _ := IntApply(rules, interface{}(data))
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

	result, _ := IntApply(rules, interface{}(data))
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

	result, _ := BoolApply(rules, interface{}(data))
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

	result, _ := StringApply(rules, interface{}(data))
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

	result, _ := StringApply(rules, interface{}(data))
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

	result, _ := BoolApply(rules, interface{}(data))
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

	log.Println("Logic ", logic)
	log.Println("Data ", data)
	log.Println("Expected ", fmt.Sprintf("%v %T", expected, expected))

	err := Apply(logic, data, &result)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Result ", fmt.Sprintf("%v %T", result, result))

	if !reflect.DeepEqual(expected, result) {
		t.Fatal("The value expected is not what we expected")
	}
}
