package jsonlogic_test

import (
	"encoding/json"
	"testing"

	"github.com/qoala-platform/jsonlogic/v3"
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

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		rule     json.RawMessage
		expected json.RawMessage
	}{
		// Default behavior (ROUND_HALF_UP)
		{
			name:     "default precision and mode",
			rule:     json.RawMessage(`{"round": [3.14159]}`),
			expected: json.RawMessage(`3.1416`),
		},
		{
			name:     "default exactly half",
			rule:     json.RawMessage(`{"round": [3.5]}`),
			expected: json.RawMessage(`3.5000`),
		},

		// ROUND_DOWN tests
		{
			name:     "round down positive",
			rule:     json.RawMessage(`{"round": [3.7, 0, "ROUND_DOWN"]}`),
			expected: json.RawMessage(`3`),
		},
		{
			name:     "round down negative",
			rule:     json.RawMessage(`{"round": [-3.7, 0, "ROUND_DOWN"]}`),
			expected: json.RawMessage(`-3`),
		},
		{
			name:     "round down positive with precision",
			rule:     json.RawMessage(`{"round": [3.77777, 2, "ROUND_DOWN"]}`),
			expected: json.RawMessage(`3.77`),
		},

		// ROUND_HALF_UP tests
		{
			name:     "round half up exactly half",
			rule:     json.RawMessage(`{"round": [3.5, 0, "ROUND_HALF_UP"]}`),
			expected: json.RawMessage(`4`),
		},
		{
			name:     "round half up negative exactly half",
			rule:     json.RawMessage(`{"round": [-3.5, 0, "ROUND_HALF_UP"]}`),
			expected: json.RawMessage(`-4`),
		},
		{
			name:     "round half up just under half",
			rule:     json.RawMessage(`{"round": [3.49999, 0, "ROUND_HALF_UP"]}`),
			expected: json.RawMessage(`3`),
		},

		// ROUND_HALF_EVEN tests
		{
			name:     "round half even to even",
			rule:     json.RawMessage(`{"round": [2.5, 0, "ROUND_HALF_EVEN"]}`),
			expected: json.RawMessage(`2`),
		},
		{
			name:     "round half even to odd",
			rule:     json.RawMessage(`{"round": [3.5, 0, "ROUND_HALF_EVEN"]}`),
			expected: json.RawMessage(`4`),
		},
		{
			name:     "round half even negative to even",
			rule:     json.RawMessage(`{"round": [-2.5, 0, "ROUND_HALF_EVEN"]}`),
			expected: json.RawMessage(`-2`),
		},

		// ROUND_CEILING tests
		{
			name:     "round ceiling positive",
			rule:     json.RawMessage(`{"round": [3.1, 0, "ROUND_CEILING"]}`),
			expected: json.RawMessage(`4`),
		},
		{
			name:     "round ceiling negative",
			rule:     json.RawMessage(`{"round": [-3.1, 0, "ROUND_CEILING"]}`),
			expected: json.RawMessage(`-3`),
		},

		// ROUND_FLOOR tests
		{
			name:     "round floor positive",
			rule:     json.RawMessage(`{"round": [3.9, 0, "ROUND_FLOOR"]}`),
			expected: json.RawMessage(`3`),
		},
		{
			name:     "round floor negative",
			rule:     json.RawMessage(`{"round": [-3.9, 0, "ROUND_FLOOR"]}`),
			expected: json.RawMessage(`-4`),
		},

		// ROUND_UP tests
		{
			name:     "round up positive",
			rule:     json.RawMessage(`{"round": [3.1, 0, "ROUND_UP"]}`),
			expected: json.RawMessage(`4`),
		},
		{
			name:     "round up negative",
			rule:     json.RawMessage(`{"round": [-3.1, 0, "ROUND_UP"]}`),
			expected: json.RawMessage(`-4`),
		},

		// ROUND_HALF_DOWN tests
		{
			name:     "round half down exactly half",
			rule:     json.RawMessage(`{"round": [3.5, 0, "ROUND_HALF_DOWN"]}`),
			expected: json.RawMessage(`3`),
		},
		{
			name:     "round half down negative exactly half",
			rule:     json.RawMessage(`{"round": [-3.5, 0, "ROUND_HALF_DOWN"]}`),
			expected: json.RawMessage(`-3`),
		},
		{
			name:     "round half down just over half",
			rule:     json.RawMessage(`{"round": [3.50001, 0, "ROUND_HALF_DOWN"]}`),
			expected: json.RawMessage(`4`),
		},

		// ROUND_05UP tests
		{
			name:     "round 05up on 0.5",
			rule:     json.RawMessage(`{"round": [2.5, 0, "ROUND_05UP"]}`),
			expected: json.RawMessage(`3`),
		},
		{
			name:     "round 05up on 0.0",
			rule:     json.RawMessage(`{"round": [2.0, 0, "ROUND_05UP"]}`),
			expected: json.RawMessage(`2`),
		},
		{
			name:     "round 05up negative on 0.5",
			rule:     json.RawMessage(`{"round": [-2.5, 0, "ROUND_05UP"]}`),
			expected: json.RawMessage(`-3`),
		},

		// Precision tests
		{
			name:     "negative precision",
			rule:     json.RawMessage(`{"round": [3141.59, -2, "ROUND_HALF_UP"]}`),
			expected: json.RawMessage(`3100`),
		},
		{
			name:     "high precision",
			rule:     json.RawMessage(`{"round": [3.14159265359, 8, "ROUND_HALF_UP"]}`),
			expected: json.RawMessage(`3.14159265`),
		},

		// Invalid mode test
		{
			name:     "invalid mode falls back to ROUND_HALF_UP",
			rule:     json.RawMessage(`{"round": [3.5, 0, "INVALID_MODE"]}`),
			expected: json.RawMessage(`4`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := jsonlogic.ApplyRaw(tt.rule, nil)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, string(tt.expected), string(output))
		})
	}
}
