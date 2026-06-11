package jsonlogic_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestFilterParseTheSubjectFromFirstPosition(t *testing.T) {
	rule := strings.NewReader(`{"filter": [
		[1,2,3,4,5],
		{"%":[{"var":""},2]}
	]}`)

	var result bytes.Buffer

	err := jsonlogic.Apply(rule, nil, &result)
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,3,5]`, result.String())
}

func TestFilterParseTheSubjectFromNullValue(t *testing.T) {
	rule := strings.NewReader(`{"filter": [
		null,
		{"%":[{"var":""},2]}
	]}`)

	var result bytes.Buffer

	err := jsonlogic.Apply(rule, nil, &result)
	assert.Nil(t, err)
	assert.JSONEq(t, `[]`, result.String())
}

func TestReduceSkipNullValues(t *testing.T) {
	rule := strings.NewReader(`{"reduce": [
		[1,2,null,4,5],
		{"+":[{"var":"current"}, {"var":"accumulator"}]},
		0
	]}`)

	var result bytes.Buffer

	err := jsonlogic.Apply(rule, nil, &result)
	assert.Nil(t, err)
	assert.JSONEq(t, `12`, result.String())
}

func TestReduceBoolValues(t *testing.T) {
	rule := strings.NewReader(`{"reduce": [
		[true,false,true,null],
		{"or":[{"var":"current"}, {"var":"accumulator"}]},
		false
	]}`)

	var result bytes.Buffer

	err := jsonlogic.Apply(rule, nil, &result)
	assert.Nil(t, err)
	assert.JSONEq(t, `true`, result.String())
}

func TestReduceStringValues(t *testing.T) {
	rule := strings.NewReader(`{"reduce": [
		["a",null,"b"],
		{"cat":[{"var":"current"}, {"var":"accumulator"}]},
		""
	]}`)

	var result bytes.Buffer

	err := jsonlogic.Apply(rule, nil, &result)
	assert.Nil(t, err)
	assert.JSONEq(t, `"ba"`, result.String())
}

func TestFilterWithMissingLogicArgument(t *testing.T) {
	// filter needs [array, logic]; omitting the logic argument must not panic.
	rule := strings.NewReader(`{"filter": [[1,2,3]]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, nil, &result)
	assert.NoError(t, err)
	assert.JSONEq(t, `[]`, result.String())
}

func TestMapWithMissingLogicArgument(t *testing.T) {
	// map needs [array, logic]; omitting the logic argument must not panic.
	rule := strings.NewReader(`{"map": [[1,2,3]]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, nil, &result)
	assert.NoError(t, err)
	assert.JSONEq(t, `[]`, result.String())
}

func TestReduceWithMissingInitialValue(t *testing.T) {
	// reduce needs [array, logic, initial]; omitting initial value must not panic.
	rule := strings.NewReader(`{"reduce": [
		[1,2,3],
		{"+":[{"var":"current"},{"var":"accumulator"}]}
	]}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, nil, &result)
	assert.NoError(t, err)
	assert.JSONEq(t, `0`, result.String())
}

// TestAllSeesOuterDataVar proves that a condition inside "all" can reference a
// variable from the outer data context, not just from the current element.
// The primitive elements (5) have no "expected" key, so the only way to resolve
// {"var":"expected"} is from outer data — the bug resolves it against the element
// instead and gets nil, making every element fail the condition.
func TestAllSeesOuterDataVar(t *testing.T) {
	rule := strings.NewReader(`{"all": [{"var": "items"}, {"==": [{"var": ""}, {"var": "expected"}]}]}`)
	data := strings.NewReader(`{"items": [5, 5, 5], "expected": 5}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, data, &result)
	assert.NoError(t, err)
	assert.JSONEq(t, `true`, result.String())
}

// TestAllSeesOuterDataVarWithRelativeElementVar is the "all" analog of TestIssue81.
// The condition mixes {"var":".B"} (current element's field) with {"var":"B"} (outer data).
// The bug resolves both against the element, making {"var":"B"} shadow the outer value
// and causing the != to evaluate as 1!=1 (false) instead of 1!=2 (true).
func TestAllSeesOuterDataVarWithRelativeElementVar(t *testing.T) {
	rule := strings.NewReader(`{"all": [{"var": "A"}, {"!=": [{"var": ".B"}, {"var": "B"}]}]}`)
	data := strings.NewReader(`{"A": [{"B": 1}, {"B": 1}], "B": 2}`)

	var result bytes.Buffer
	err := jsonlogic.Apply(rule, data, &result)
	assert.NoError(t, err)
	assert.JSONEq(t, `true`, result.String())
}
