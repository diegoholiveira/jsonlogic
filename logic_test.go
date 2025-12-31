package jsonlogic_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestAndReturnsFirstFalsyArgument(t *testing.T) {
	rule := `{"and":["", true, "unused"]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `""`, result.String())
}

func TestAndReturnsFirstFalsyEmptyArray(t *testing.T) {
	rule := `{"and":[[], 1, "done"]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `[]`, result.String())
}

func TestAndReturnsNullWhenNoArguments(t *testing.T) {
	rule := `{"and":[]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `null`, result.String())
}

func TestAndReturnsFirstFalsyNumber(t *testing.T) {
	rule := `{"and":[0, 2, 3]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `0`, result.String())
}

func TestAndReturnsLastArgumentWhenAllTruthy(t *testing.T) {
	rule := `{"and":[true, 1, "done"]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `"done"`, result.String())
}

func TestAndReturnsLastArgumentWhenAllTruthyNoNumbers(t *testing.T) {
	rule := `{"and":[true, "yes"]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `"yes"`, result.String())
}

func TestAndReturnsFalsyAfterTruthyArray(t *testing.T) {
	rule := `{"and":[[1], 0, "unused"]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `0`, result.String())
}

func TestOrReturnsFirstTruthyArgument(t *testing.T) {
	rule := `{"or":["ok", 0, false]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `"ok"`, result.String())
}

func TestOrReturnsNullWhenNoArguments(t *testing.T) {
	rule := `{"or":[]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `null`, result.String())
}

func TestOrReturnsFirstTruthyNumber(t *testing.T) {
	rule := `{"or":[1, false, "later"]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `1`, result.String())
}

func TestOrReturnsFirstTruthyArray(t *testing.T) {
	rule := `{"or":[[1, 2], 0, false]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `[1,2]`, result.String())
}

func TestOrReturnsLastArgumentWhenAllFalsy(t *testing.T) {
	rule := `{"or":[0, false, ""]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `""`, result.String())
}

func TestOrReturnsLastFalsyArrayWhenAllFalsy(t *testing.T) {
	rule := `{"or":[0, false, []]}`

	var result bytes.Buffer
	err := jsonlogic.Apply(strings.NewReader(rule), strings.NewReader(`{}`), &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `[]`, result.String())
}
