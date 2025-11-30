package jsonlogic

import (
	"bytes"
	"strings"
	"testing"
)

func TestContainsAll(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		data     string
		expected string
	}{
		{
			name:     "all elements present",
			rule:     `{"contains_all": [["a", "b", "c"], ["a", "b"]]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "all elements present - exact match",
			rule:     `{"contains_all": [["a", "b"], ["a", "b"]]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "some elements missing",
			rule:     `{"contains_all": [["a", "b"], ["a", "b", "c"]]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "empty required array",
			rule:     `{"contains_all": [["a", "b", "c"], []]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "empty search array",
			rule:     `{"contains_all": [[], ["a"]]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "with variable",
			rule:     `{"contains_all": [{"var": "selected"}, ["vip", "premium"]]}`,
			data:     `{"selected": ["vip", "premium", "gold"]}`,
			expected: "true",
		},
		{
			name:     "with variable - missing element",
			rule:     `{"contains_all": [{"var": "selected"}, ["vip", "diamond"]]}`,
			data:     `{"selected": ["vip", "premium", "gold"]}`,
			expected: "false",
		},
		{
			name:     "with numbers",
			rule:     `{"contains_all": [[1, 2, 3, 4], [1, 3]]}`,
			data:     `{}`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bytes.Buffer
			err := Apply(strings.NewReader(tt.rule), strings.NewReader(tt.data), &result)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.TrimSpace(result.String()) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		data     string
		expected string
	}{
		{
			name:     "one element present",
			rule:     `{"contains_any": [["a", "b", "c"], ["x", "b"]]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "multiple elements present",
			rule:     `{"contains_any": [["a", "b", "c"], ["a", "c"]]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "no elements present",
			rule:     `{"contains_any": [["a", "b", "c"], ["x", "y"]]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "empty check array",
			rule:     `{"contains_any": [["a", "b", "c"], []]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "empty search array",
			rule:     `{"contains_any": [[], ["a"]]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "with variable",
			rule:     `{"contains_any": [{"var": "tags"}, ["urgent", "important"]]}`,
			data:     `{"tags": ["normal", "urgent"]}`,
			expected: "true",
		},
		{
			name:     "with variable - no match",
			rule:     `{"contains_any": [{"var": "tags"}, ["urgent", "important"]]}`,
			data:     `{"tags": ["normal", "low"]}`,
			expected: "false",
		},
		{
			name:     "with numbers",
			rule:     `{"contains_any": [[1, 2, 3], [5, 3, 7]]}`,
			data:     `{}`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bytes.Buffer
			err := Apply(strings.NewReader(tt.rule), strings.NewReader(tt.data), &result)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.TrimSpace(result.String()) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestContainsNone(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		data     string
		expected string
	}{
		{
			name:     "no elements present",
			rule:     `{"contains_none": [["a", "b", "c"], ["x", "y"]]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "one element present",
			rule:     `{"contains_none": [["a", "b", "c"], ["x", "b"]]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "all elements present",
			rule:     `{"contains_none": [["a", "b", "c"], ["a", "b"]]}`,
			data:     `{}`,
			expected: "false",
		},
		{
			name:     "empty check array",
			rule:     `{"contains_none": [["a", "b", "c"], []]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "empty search array",
			rule:     `{"contains_none": [[], ["a"]]}`,
			data:     `{}`,
			expected: "true",
		},
		{
			name:     "with variable - blocked words not present",
			rule:     `{"contains_none": [{"var": "content"}, ["spam", "blocked"]]}`,
			data:     `{"content": ["hello", "world"]}`,
			expected: "true",
		},
		{
			name:     "with variable - blocked word present",
			rule:     `{"contains_none": [{"var": "content"}, ["spam", "blocked"]]}`,
			data:     `{"content": ["hello", "spam"]}`,
			expected: "false",
		},
		{
			name:     "with numbers",
			rule:     `{"contains_none": [[1, 2, 3], [7, 8, 9]]}`,
			data:     `{}`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bytes.Buffer
			err := Apply(strings.NewReader(tt.rule), strings.NewReader(tt.data), &result)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.TrimSpace(result.String()) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
