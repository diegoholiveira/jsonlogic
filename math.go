package jsonlogic

import "math"

func mod(values, data any) any {
	parsed := parseValues(values, data).([]any)

	a := toNumber(parsed[0])
	b := toNumber(parsed[1])

	return math.Mod(a, b)
}

func abs(values, data any) any {
	parsed := parseValues(values, data)
	parsedAsSlice, ok := parsed.([]any)
	if !ok {
		return math.Abs(toNumber(parsed))
	}

	if len(parsedAsSlice) == 0 {
		return float64(0)
	}

	return math.Abs(toNumber(parsedAsSlice[0]))
}

func sum(values, data any) any {
	parsed := parseValues(values, data)
	parsedAsSlice, ok := parsed.([]any)
	if !ok {
		return toNumber(parsed)
	}

	sum := float64(0)

	for _, n := range parsedAsSlice {
		sum += toNumber(n)
	}

	return sum
}

func minus(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok || len(parsed) == 0 {
		return 0
	}

	if len(parsed) == 1 {
		return -1 * toNumber(parsed[0])
	}

	sum := toNumber(parsed[0])

	for i := 1; len(parsed) > i; i++ {
		sum -= toNumber(parsed[i])
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
		sum *= toNumber(n)
	}

	return sum
}

func div(values, data any) any {
	parsed, ok := parseValues(values, data).([]any)
	if !ok || len(parsed) == 0 {
		return 0
	}

	sum := toNumber(parsed[0])

	for i := 1; len(parsed) > i; i++ {
		sum = sum / toNumber(parsed[i])
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

	bigger := toNumber(parsed[0])

	for i := 1; i < size; i++ {
		if n := toNumber(parsed[i]); n > bigger {
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

	smallest := toNumber(parsed[0])

	for i := 1; i < size; i++ {
		if n := toNumber(parsed[i]); smallest > n {
			smallest = n
		}
	}

	return smallest
}
