package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRecoverWithNonErrorPanicValue verifies that a panic carrying a value that
// does NOT implement the error interface (e.g. a custom operator doing
// panic("boom")) is converted into a returned error rather than escaping the
// library.
//
// Before the fix, applyInterfaceUnguarded's recover did `err = e.(error)`,
// which itself panicked on a non-error value and escaped, crashing the caller.
func TestRecoverWithNonErrorPanicValue(t *testing.T) {
	AddOperator("panic_with_string", func(values, data any) any {
		panic("boom: not an error value")
	})

	rule := map[string]any{
		"panic_with_string": []any{},
	}
	data := map[string]any{}

	var (
		gotResult any
		gotErr    error
		escaped   any
		didPanic  bool
	)

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
				escaped = r
			}
		}()
		gotResult, gotErr = ApplyInterface(rule, data)
	}()

	assert.False(t, didPanic,
		"panic escaped the library's recover instead of being returned as an error: %v", escaped)
	assert.Error(t, gotErr,
		"a non-error panic must be converted into a returned error")
	assert.Nil(t, gotResult)
}
