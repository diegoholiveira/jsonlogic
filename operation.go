package jsonlogic

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/diegoholiveira/jsonlogic/v3/internal/javascript"
)

// OperatorFn defines the signature for custom operator functions.
// It takes values and data as input and returns a result.
// See AddOperator for the types it should return.
type OperatorFn func(values, data any) (result any)

// ErrInvalidOperator represents an error when an unsupported operator is used.
// It contains the operator name that caused the error.
type ErrInvalidOperator struct {
	operator string
}

func (e ErrInvalidOperator) Error() string {
	return fmt.Sprintf("The operator \"%s\" is not supported", e.operator)
}

// operators holds the registered operators as a map[string]OperatorFn. A stored map
// is never modified, so evaluating a rule reads it without locking.
var operators atomic.Value

// operatorsLock serializes writers, which replace the map with an updated copy.
var operatorsLock sync.Mutex

func lookupOperator(key string) (OperatorFn, bool) {
	opFn, found := operators.Load().(map[string]OperatorFn)[key]
	return opFn, found
}

// AddOperator registers a custom operator with the given key and function.
// The operator function will be called with parsed values and the original data context.
//
// Parameters:
//   - key: the operator name to register (e.g., "custom_op")
//   - cb: the function to execute when the operator is encountered
//
// The result should be one of the types encoding/json uses to unmarshal into an
// interface value: bool, float64, string, []any, map[string]any or nil. As a
// convenience, the other built-in Go integer and float types are converted to
// float64, including inside []any and map[string]any. Integers beyond 2^53 may lose
// precision. Any other type is left unchanged and built-in operators may not
// handle it.
//
// The callback must treat values and data as read-only. Values written into them are
// neither checked nor converted.
//
// Concurrency: This function is safe for concurrent use. Each call copies the operators
// map so that evaluating rules needs no lock.
func AddOperator(key string, cb OperatorFn) {
	operatorsLock.Lock()
	defer operatorsLock.Unlock()

	current := operators.Load().(map[string]OperatorFn)
	updated := make(map[string]OperatorFn, len(current)+1)
	for name, opFn := range current {
		updated[name] = opFn
	}
	updated[key] = func(values, data any) any {
		return normalizeNumbers(cb(parseValues(values, data), data))
	}
	operators.Store(updated)
}

func operation(operator string, values, data any) any {
	if opFn, found := lookupOperator(operator); found {
		return opFn(values, data)
	}

	panic(ErrInvalidOperator{
		operator: operator,
	})
}

func init() {
	builtin := make(map[string]OperatorFn)

	builtin["and"] = _and
	builtin["or"] = _or
	builtin["filter"] = filter
	builtin["map"] = _map
	builtin["reduce"] = reduce
	builtin["all"] = all
	builtin["none"] = none
	builtin["some"] = some
	builtin["in"] = _in
	builtin["missing"] = missing
	builtin["missing_some"] = missingSome
	builtin["var"] = getVar
	builtin["set"] = setProperty
	builtin["cat"] = concat
	builtin["substr"] = substr
	builtin["merge"] = merge
	builtin["if"] = conditional
	builtin["?:"] = conditional
	builtin["max"] = max
	builtin["min"] = min
	builtin["+"] = sum
	builtin["-"] = minus
	builtin["*"] = mult
	builtin["/"] = div
	builtin["%"] = mod
	builtin["abs"] = abs
	builtin["!"] = negative
	builtin["!!"] = func(v, d any) any { return !javascript.IsTrue(negative(v, d)) }
	builtin["==="] = hardEquals
	builtin["!=="] = func(v, d any) any { return !hardEquals(v, d).(bool) }
	builtin["<"] = isLessThan
	builtin["<="] = isLessOrEqualThan
	builtin[">"] = isGreaterThan
	builtin[">="] = isGreaterOrEqualThan
	builtin["=="] = isEqual
	builtin["!="] = func(v, d any) any { return !isEqual(v, d).(bool) }

	/* CUSTOM OPERATORS */
	builtin["contains_all"] = func(v, d any) any { return containsAll(parseValues(v, d)) }
	builtin["contains_any"] = func(v, d any) any { return containsAny(parseValues(v, d)) }
	builtin["contains_none"] = func(v, d any) any { return containsNone(parseValues(v, d)) }

	operators.Store(builtin)
}
