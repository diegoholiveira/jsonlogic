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
	var sum float64

	for _, n := range values.([]interface{}) {
		if sum == 0 {
			sum = toNumber(n)

			continue
		}

		sum -= toNumber(n)
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
	var sum float64

	for i, n := range values.([]interface{}) {
		if i == 0 {
			sum = toNumber(n)

			continue
		}

		sum = sum / toNumber(n)
	}

	return sum
}
