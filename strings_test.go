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
