package jsonlogic

import (
	"fmt"
	"sync"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

type OperatorFn func(values, data any) (result any)

type ErrInvalidOperator struct {
	operator string
}

func (e ErrInvalidOperator) Error() string {
	return fmt.Sprintf("The operator \"%s\" is not supported", e.operator)
}

// operators holds custom operators
var operators = make(map[string]OperatorFn)
var operatorsLock = &sync.RWMutex{}

// AddOperator allows for custom operators to be used
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
