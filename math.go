package jsonlogic

import (
	"math"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func mod(values, data any) any {
	_values := parseValues(values, data).([]any)

	a, b := _values[0], _values[1]

	_a := typing.ToNumber(a)
	_b := typing.ToNumber(b)

	return math.Mod(_a, _b)
}

func abs(values, data any) any {
	values = parseValues(values, data)
	if typing.IsSlice(values) {
		return math.Abs(typing.ToNumber(values.([]any)[0]))
	}

	return math.Abs(typing.ToNumber(values))
}

func sum(values, data any) any {
	values = parseValues(values, data)
	if !typing.IsSlice(values) {
		return typing.ToNumber(values)
	}

	inputSlice := values.([]any)
	sliceLen := len(inputSlice)

	if sliceLen == 0 {
		return float64(0)
	}

	if sliceLen == 1 {
		return typing.ToNumber(inputSlice[0])
	}

	sum := float64(0)
	for _, n := range inputSlice {
		sum += typing.ToNumber(n)
	}

	return sum
}

func minus(values, data any) any {
	_values := parseValues(values, data).([]any)

	if len(_values) == 0 {
		return 0
	}

	if len(_values) == 1 {
		return -1 * typing.ToNumber(_values[0])
	}

	sum := typing.ToNumber(_values[0])
	for i := 1; len(_values) > i; i++ {
		sum -= typing.ToNumber(_values[i])
	}

	return sum
}

func mult(values, data any) any {
	values = parseValues(values, data)

	sum := float64(1)

	for _, n := range values.([]any) {
		sum *= typing.ToNumber(n)
	}

	return sum
}

func div(values, data any) any {
	_values := parseValues(values, data).([]any)

	if len(_values) == 0 {
		return 0
	}

	sum := typing.ToNumber(_values[0])
	for i := 1; len(_values) > i; i++ {
		sum = sum / typing.ToNumber(_values[i])
	}

	return sum
}

func max(values, data any) any {
	values = parseValues(values, data)
	parsed := values.([]any)
	size := len(parsed)
	if size == 0 {
		return nil
	}

	bigger := typing.ToNumber(parsed[0])

	for i := 1; i < size; i++ {
		_n := typing.ToNumber(parsed[i])
		if _n > bigger {
			bigger = _n
		}
	}

	return bigger
}

func min(values, data any) any {
	values = parseValues(values, data)
	parsed := values.([]any)
	size := len(parsed)
	if size == 0 {
		return nil
	}

	smallest := typing.ToNumber(parsed[0])

	for i := 1; i < size; i++ {
		_n := typing.ToNumber(parsed[i])
		if smallest > _n {
			smallest = _n
		}
	}

	return smallest
}
