package jsonlogic

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONLogicValidator(t *testing.T) {
	scenarios := map[string]struct {
		IsValid bool
		Rule    io.Reader
	}{
		"invalid rule": {
			IsValid: false,
			Rule:    strings.NewReader(`{"a", "b"}`),
		},
		"invalid operator": {
			IsValid: false,
			Rule:    strings.NewReader(`{"filt":[[10, 1, 100], {">=":[{"var":""},2]}]}`),
		},
		"invalid condition inside a filter": {
			IsValid: false,
			Rule:    strings.NewReader(`{"filter":[{"var":"integers"}, {"=": [{"var":""}, [10]]}]}`),
		},
		"primitive is a valid rule": {
			IsValid: true,
			Rule:    strings.NewReader(`10`),
		},
		"set must be valid": {
			IsValid: true,
			Rule: strings.NewReader(`{
				"map": [
					{"var": "objects"},
					{"set": [
						{"var": ""},
						"age",
						{"+": [{"var": ".age"}, 2]}
					]}
				]
			}`),
		},
	}

	for name, scenario := range scenarios {
		t.Run(fmt.Sprintf("SCENARIO:%s", name), func(t *testing.T) {
			assert.Equal(t, scenario.IsValid, IsValid(scenario.Rule))
		})
	}
}
