package jsonlogic

import (
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// restoreOperators undoes the AddOperator calls of a test.
func restoreOperators(t *testing.T) {
	t.Helper()
	saved := operators.Load()
	t.Cleanup(func() { operators.Store(saved) })
}

// TestConcurrentApplyAndAddOperator validates that validating rules and adding operators concurrently
// doesn't cause fatal errors or deadlocks.
func TestConcurrentValidationAndAddOperator(t *testing.T) {
	var wg sync.WaitGroup
	numRoutines := 10
	numIterations := 100

	// Start multiple goroutines to validate rules concurrently
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				rule := `{"==": [1, 1]}`
				_ = IsValid(strings.NewReader(rule))
			}
		}()
	}

	// Start a goroutine to add a new operator concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < numIterations; j++ {
			AddOperator("test_op", func(values, data any) any {
				return "test"
			})
		}
	}()

	wg.Wait()
}

// TestConcurrentApplyAndAddOperator validates that applying rules and adding operators concurrently
// doesn't cause fatal errors or deadlocks.
func TestConcurrentApplyAndAddOperator(t *testing.T) {
	var wg sync.WaitGroup
	numRoutines := 10
	numIterations := 100

	// Start multiple goroutines to apply rules concurrently
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				rule := `{"==": [1, 1]}`
				data := `{}`
				_ = Apply(strings.NewReader(rule), strings.NewReader(data), io.Discard)
			}
		}()
	}

	// Start a goroutine to add a new operator concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < numIterations; j++ {
			AddOperator("test_op", func(values, data any) any {
				return "test"
			})
		}
	}()

	wg.Wait()
}

// Writers replace the operators map, so concurrent registrations must not lose each other.
func TestConcurrentAddOperatorKeepsEveryOperator(t *testing.T) {
	restoreOperators(t)
	names := make([]string, 50)
	var wg sync.WaitGroup
	for i := range names {
		names[i] = "concurrent_op_" + strings.Repeat("x", i)
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			AddOperator(name, func(_, _ any) any { return name })
		}(names[i])
	}
	wg.Wait()

	for _, name := range append(names, "==") {
		assert.True(t, ValidateJsonLogic(map[string]any{name: []any{}}), name)
	}
}

func TestOperatorCanRegisterAnotherOperator(t *testing.T) {
	restoreOperators(t)
	AddOperator("registers_other", func(_, _ any) any {
		AddOperator("registered_by_other", func(_, _ any) any { return "ok" })
		return true
	})

	got, err := ApplyInterface(map[string]any{"and": []any{
		map[string]any{"registers_other": []any{}},
		map[string]any{"registered_by_other": []any{}},
	}}, nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", got)
}
