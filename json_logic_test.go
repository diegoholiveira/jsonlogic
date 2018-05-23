package json_logic

import (
	"encoding/json"
	"testing"
)

func TestAlwaysShouldAlwaysPass(t *testing.T) {
	result, _ := Apply(true, nil)
	if !result {
		t.Fatal("Always should always pass")
	}
}

func TestNeverShouldNeverPass(t *testing.T) {
	result, _ := Apply(false, nil)
	if result {
		t.Fatal("Always should never pass")
	}
}

func TestRootElement(t *testing.T) {
	rules := []int{1, 1}

	_, err := Apply(rules, nil)
	if err == nil {
		t.Fatal("We must force the root element to be an object")
	}
}

func TestSimpleComparisonWithInteger(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte("{\"==\":[1, 1]}"), &rules)

	result, _ := Apply(rules, nil)
	if !result {
		t.Fatal("A simple comparison is expected to be true")
	}
}

func TestSimpleComparisonWithString(t *testing.T) {
	var rules interface{}
	json.Unmarshal([]byte("{\"==\":[\"a\", \"a\"]}"), &rules)

	result, _ := Apply(rules, nil)
	if !result {
		t.Fatal("A simple comparison is expected to be true")
	}
}
