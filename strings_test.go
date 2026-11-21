package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlogic "github.com/diegoholiveira/jsonlogic/v3"
)

func TestCat(t *testing.T) {
	testCases := []struct {
		name     string
		rule     string
		data     string
		expected string
	}{
		{
			name:     "Empty string",
			rule:     `{"cat": ""}`,
			data:     `{}`,
			expected: `""`,
		},
		{
			name:     "Empty array",
			rule:     `{"cat": []}`,
			data:     `{}`,
			expected: `""`,
		},
		{
			name:     "Single string",
			rule:     `{"cat": "hello"}`,
			data:     `{}`,
			expected: `"hello"`,
		},
		{
			name:     "Multiple strings",
			rule:     `{"cat": ["hello", " ", "world"]}`,
			data:     `{}`,
			expected: `"hello world"`,
		},
		{
			name:     "Keeps leading and trailing spaces",
			rule:     `{"cat": [" hello", "world "]}`,
			data:     `{}`,
			expected: `" helloworld "`,
		},
		{
			name:     "Keeps a space from data",
			rule:     `{"cat": ["Hello, ", {"var": "name"}]}`,
			data:     `{"name": ""}`,
			expected: `"Hello, "`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := json.RawMessage(tc.rule)
			data := json.RawMessage(tc.data)
			expected := json.RawMessage(tc.expected)

			output, err := jsonlogic.ApplyRaw(rule, data)
			if err != nil {
				t.Fatal(err)
			}

			assert.JSONEq(t, string(expected), string(output))
		})
	}
}

func TestSubstr(t *testing.T) {
	// Expected values match json-logic-js, which uses String.prototype.substr.
	testCases := []struct {
		rule     string
		expected string
	}{
		{`{"substr": ["jsonlogic", 4]}`, `"logic"`},
		{`{"substr": ["jsonlogic", -5]}`, `"logic"`},
		{`{"substr": ["jsonlogic", 1, 3]}`, `"son"`},
		{`{"substr": ["jsonlogic", 4, -2]}`, `"log"`},
		{`{"substr": ["jsonlogic", 9]}`, `""`},
		{`{"substr": ["jsonlogic", 10]}`, `""`},
		{`{"substr": ["jsonlogic", 10, 2]}`, `""`},
		{`{"substr": ["jsonlogic", -10]}`, `"jsonlogic"`},
		{`{"substr": ["jsonlogic", -10, 2]}`, `"js"`},
		{`{"substr": ["jsonlogic", 0, 0]}`, `""`},
		{`{"substr": ["jsonlogic", 0, 20]}`, `"jsonlogic"`},
		{`{"substr": ["jsonlogic", 5, -5]}`, `""`},
		{`{"substr": ["jsonlogic", 0, -10]}`, `""`},
		{`{"substr": ["", 0, -1]}`, `""`},
		{`{"substr": ["héllo", 1, -1]}`, `"éll"`},
	}

	for _, tc := range testCases {
		t.Run(tc.rule, func(t *testing.T) {
			output, err := jsonlogic.ApplyRaw(json.RawMessage(tc.rule), json.RawMessage(`{}`))
			if err != nil {
				t.Fatal(err)
			}

			assert.JSONEq(t, tc.expected, string(output))
		})
	}
}
