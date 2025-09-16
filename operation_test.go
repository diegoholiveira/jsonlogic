package jsonlogic

import (
	"io"
	"strings"
	"sync"
	"testing"
)

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
