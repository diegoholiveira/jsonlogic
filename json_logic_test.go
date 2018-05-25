package json_logic

import (
	"encoding/json"
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

func TestRootElement(t *testing.T) {
	rules := []int{1, 1}

	_, err := BoolApply(rules, nil)
	if err == nil {
		t.Fatal("We must force the root element to be an object")
	}
}

func TestSimpleComparisonWithInteger(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte("{\"==\":[1, 1]}"), &rules)

	result, _ := BoolApply(rules, nil)
	if !result {
		t.Fatal("A simple comparison is expected to be true")
	}
}

func TestSimpleComparisonWithString(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte("{\"==\":[\"a\", \"a\"]}"), &rules)

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
	json.Unmarshal([]byte("{\"var\":\"a\"}"), &rules)

	data := map[string]interface{}{"a": 10}

	result, _ := IntApply(rules, interface{}(data))
	if result != 10 {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestSimpleVarWithoutSyntacticSugar(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte("{\"var\":[\"a\"]}"), &rules)

	data := map[string]interface{}{"a": 10}

	result, _ := IntApply(rules, interface{}(data))
	if result != 10 {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestVariableWithDefaultValue(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte("{\"var\":[\"z\", 20]}"), &rules)

	data := map[string]interface{}{"a": 10}

	result, _ := IntApply(rules, interface{}(data))
	if result != 20 {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestSimpleVarComparison(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte(`{
		"==": [
			{"var": "a"},
			10
		]
	}`), &rules)

	data := map[string]interface{}{"a": 10}

	result, _ := BoolApply(rules, interface{}(data))
	if !result {
		t.Fatal("The value expected must be equal the value of the context")
	}
}

func TestComposedVar(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte(`{
		"var": "champ.name"
	}`), &rules)

	data := map[string]interface{}{"champ": interface{}(map[string]interface{}{"name": "Diego"})}

	result, _ := StringApply(rules, interface{}(data))
	if result != "Diego" {
		t.Fatal("The value expected must be equal the value of the context")
	}
}
