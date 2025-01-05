package jsonlogic

import "math"

func mod(a any, b any) any {
	_a := toNumber(a)
	_b := toNumber(b)

	return math.Mod(_a, _b)
}

func abs(a any) any {
	_a := toNumber(a)

	return math.Abs(_a)
}

func sum(values any) any {
	sum := float64(0)

	for _, n := range values.([]any) {
		sum += toNumber(n)
	}

	return sum
}

func minus(values any) any {
	_values := toSliceOfNumbers(values)

	sum := _values[0]
	for i := 1; len(_values) > i; i++ {
		sum -= _values[i]
	}

	return sum
}

func mult(values any) any {
	sum := float64(1)

	for _, n := range values.([]any) {
		sum *= toNumber(n)
	}

	return sum
}

func div(values any) any {
	_values := toSliceOfNumbers(values)

	sum := _values[0]
	for i := 1; len(_values) > i; i++ {
		sum = sum / _values[i]
	}

	return sum
}
