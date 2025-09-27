package jsonlogic

import (
	"fmt"
	"sync"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

// OperatorFn defines the signature for custom operator functions.
// It takes values and data as input and returns a result.
type OperatorFn func(values, data any) (result any)

// ErrInvalidOperator represents an error when an unsupported operator is used.
// It contains the operator name that caused the error.
type ErrInvalidOperator struct {
	operator string
}

func (e ErrInvalidOperator) Error() string {
	return fmt.Sprintf("The operator \"%s\" is not supported", e.operator)
}

// operators holds custom operators
var operators = make(map[string]OperatorFn)

var operatorsLock = &sync.RWMutex{}

// AddOperator registers a custom operator with the given key and function.
// The operator function will be called with parsed values and the original data context.
//
// Parameters:
//   - key: the operator name to register (e.g., "custom_op")
//   - cb: the function to execute when the operator is encountered
//
// Concurrency: This function is safe for concurrent use as it properly locks the operators map.
func AddOperator(key string, cb OperatorFn) {
	operatorsLock.Lock()
	defer operatorsLock.Unlock()

	operators[key] = func(values, data any) any {
		return cb(parseValues(values, data), data)
	}
}

func operation(operator string, values, data any) any {
	operatorsLock.RLock()
	opFn, found := operators[operator]
	operatorsLock.RUnlock()
	if found {
		return opFn(values, data)
	}

	panic(ErrInvalidOperator{
		operator: operator,
	})
}

func init() {
	operatorsLock.Lock()
	defer operatorsLock.Unlock()

	operators["and"] = _and
	operators["or"] = _or
	operators["filter"] = filter
	operators["map"] = _map
	operators["reduce"] = reduce
	operators["all"] = all
	operators["none"] = none
	operators["some"] = some
	operators["in"] = _in
	operators["missing"] = missing
	operators["missing_some"] = missingSome
	operators["var"] = getVar
	operators["set"] = setProperty
	operators["cat"] = concat
	operators["substr"] = substr
	operators["merge"] = merge
	operators["if"] = conditional
	operators["?:"] = conditional
	operators["max"] = max
	operators["min"] = min
	operators["+"] = sum
	operators["-"] = minus
	operators["*"] = mult
	operators["/"] = div
	operators["%"] = mod
	operators["abs"] = abs
	operators["!"] = negative
	operators["!!"] = func(v, d any) any { return !typing.IsTrue(negative(v, d)) }
	operators["==="] = hardEquals
	operators["!=="] = func(v, d any) any { return !hardEquals(v, d).(bool) }
	operators["<"] = isLessThan
	operators["<="] = isLessOrEqualThan
	operators[">"] = isGreaterThan
	operators[">="] = isGreaterOrEqualThan
	operators["=="] = isEqual
	operators["!="] = func(v, d any) any { return !isEqual(v, d).(bool) }
}
