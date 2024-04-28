package issues_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/jsonlogic/v3"
)

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
