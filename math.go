package jsonlogic

import (
	"math"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func mod(values, data any) any {
	parsed := parseValues(values, data).([]any)

	a := typing.ToNumber(parsed[0])
	b := typing.ToNumber(parsed[1])

	return math.Mod(a, b)
}

func abs(values, data any) any {
	parsed := parseValues(values, data)
	parsedAsSlice, ok := parsed.([]any)
	if !ok {
		return math.Abs(typing.ToNumber(parsed))
	}

	if len(parsedAsSlice) == 0 {
		return float64(0)
	}

	return math.Abs(typing.ToNumber(parsedAsSlice[0]))
}

func sum(values, data any) any {
	parsed := parseValues(values, data)
	parsedAsSlice, ok := parsed.([]any)
	if !ok {
		return typing.ToNumber(parsed)
	}

	sum := float64(0)

	for _, n := range parsedAsSlice {
		sum += typing.ToNumber(n)
	}

	return sum
}

func minus(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok || len(parsed) == 0 {
		return 0
	}

	if len(parsed) == 1 {
		return -1 * typing.ToNumber(parsed[0])
	}

	sum := typing.ToNumber(parsed[0])

	for i := 1; len(parsed) > i; i++ {
		sum -= typing.ToNumber(parsed[i])
	}

	return sum
}

func mult(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok || len(parsed) == 0 {
		return float64(1)
	}

	sum := float64(1)

	for _, n := range parsed {
		sum *= typing.ToNumber(n)
	}

	return sum
}

func div(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok || len(parsed) == 0 {
		return 0
	}

	sum := typing.ToNumber(parsed[0])

	for i := 1; len(parsed) > i; i++ {
		sum = sum / typing.ToNumber(parsed[i])
	}

	return sum
}

func max(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok {
		return nil
	}

	size := len(parsed)
	if size == 0 {
		return nil
	}

	bigger := typing.ToNumber(parsed[0])

	for i := 1; i < size; i++ {
		if n := typing.ToNumber(parsed[i]); n > bigger {
			bigger = n
		}
	}

	return bigger
}

func min(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok {
		return nil
	}

	size := len(parsed)
	if size == 0 {
		return nil
	}

	smallest := typing.ToNumber(parsed[0])

	for i := 1; i < size; i++ {
		if n := typing.ToNumber(parsed[i]); smallest > n {
			smallest = n
		}
	}

	return smallest
}
