package jsonlogic

import (
	"math"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func mod(a any, b any) any {
	_a := typing.ToNumber(a)
	_b := typing.ToNumber(b)

	return math.Mod(_a, _b)
}

func abs(a any) any {
	_a := typing.ToNumber(a)

	return math.Abs(_a)
}

func sum(values any) any {
	sum := float64(0)

	for _, n := range values.([]any) {
		sum += typing.ToNumber(n)
	}

	return sum
}

func minus(values any) any {
	_values := values.([]any)

	if len(_values) == 0 {
		return 0
	}

	sum := typing.ToNumber(_values[0])
	for i := 1; len(_values) > i; i++ {
		sum -= typing.ToNumber(_values[i])
	}

	return sum
}

func mult(values any) any {
	sum := float64(1)

	for _, n := range values.([]any) {
		sum *= typing.ToNumber(n)
	}

	return sum
}

func div(values any) any {
	_values := values.([]any)

	if len(_values) == 0 {
		return 0
	}

	sum := typing.ToNumber(_values[0])
	for i := 1; len(_values) > i; i++ {
		sum = sum / typing.ToNumber(_values[i])
	}

	return sum
}
