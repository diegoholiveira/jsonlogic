package jsonlogic

import "math"

func mod(a interface{}, b interface{}) interface{} {
	_a := toNumber(a)
	_b := toNumber(b)

	return math.Mod(_a, _b)
}

func abs(a interface{}) interface{} {
	_a := toNumber(a)

	return math.Abs(_a)
}

func sum(values interface{}) interface{} {
	sum := float64(0)

	for _, n := range values.([]interface{}) {
		sum += toNumber(n)
	}

	return sum
}

func minus(values interface{}) interface{} {
	_values := toSliceOfNumbers(values)

	sum := _values[0]
	for i := 1; len(_values) > i; i++ {
		sum -= _values[i]
	}

	return sum
}

func mult(values interface{}) interface{} {
	sum := float64(1)

	for _, n := range values.([]interface{}) {
		sum *= toNumber(n)
	}

	return sum
}

func div(values interface{}) interface{} {
	_values := toSliceOfNumbers(values)

	sum := _values[0]
	for i := 1; len(_values) > i; i++ {
		sum = sum / _values[i]
	}

	return sum
}
