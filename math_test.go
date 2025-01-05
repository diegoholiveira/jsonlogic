package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/stretchr/testify/assert"
)

func TestSubOperation(t *testing.T) {
	var rule json.RawMessage = json.RawMessage(`{
		"-": [
			0,
			10
		]
	}`)

	var expected json.RawMessage = json.RawMessage("-10")

	output, err := jsonlogic.ApplyRaw(rule, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expected), string(output))
}
