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
