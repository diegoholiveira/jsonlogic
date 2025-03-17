package jsonlogic

import (
	"math"

	"github.com/qoala-platform/jsonlogic/v3/internal/typing"
)

// Rounding modes
const (
	ROUND_DOWN      = "ROUND_DOWN"      // Round towards zero
	ROUND_HALF_UP   = "ROUND_HALF_UP"   // Round to nearest, ties away from zero
	ROUND_HALF_EVEN = "ROUND_HALF_EVEN" // Round to nearest, ties to even
	ROUND_CEILING   = "ROUND_CEILING"   // Round towards positive infinity
	ROUND_FLOOR     = "ROUND_FLOOR"     // Round towards negative infinity
	ROUND_UP        = "ROUND_UP"        // Round away from zero
	ROUND_HALF_DOWN = "ROUND_HALF_DOWN" // Round to nearest, ties towards zero
	ROUND_05UP      = "ROUND_05UP"      // Round zero or five away from zero
)

// isValidRoundingMode checks if the provided mode is valid
func isValidRoundingMode(mode string) bool {
	switch mode {
	case ROUND_DOWN, ROUND_HALF_UP, ROUND_HALF_EVEN, ROUND_CEILING,
		ROUND_FLOOR, ROUND_UP, ROUND_HALF_DOWN, ROUND_05UP:
		return true
	default:
		return false
	}
}

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

func max(values any) any {
	converted := values.([]any)
	size := len(converted)
	if size == 0 {
		return nil
	}

	bigger := typing.ToNumber(converted[0])

	for i := 1; i < size; i++ {
		_n := typing.ToNumber(converted[i])
		if _n > bigger {
			bigger = _n
		}
	}

	return bigger
}

func min(values any) any {
	converted := values.([]any)
	size := len(converted)
	if size == 0 {
		return nil
	}

	smallest := typing.ToNumber(converted[0])

	for i := 1; i < size; i++ {
		_n := typing.ToNumber(converted[i])
		if smallest > _n {
			smallest = _n
		}
	}

	return smallest
}

func floor(value any) any {
	converted := value.(float64)
	return math.Floor(converted)
}

func ceil(value any) any {
	converted := value.(float64)
	return math.Ceil(converted)
}

func round(values any) any {
	parsed := values.([]any)
	if len(parsed) == 0 {
		return 0
	}

	num := typing.ToNumber(parsed[0])
	precision := 4.0      // default precision
	mode := ROUND_HALF_UP // default mode

	// Check if precision is specified
	if len(parsed) > 1 {
		precision = typing.ToNumber(parsed[1])
	}

	// Check if rounding mode is specified
	if len(parsed) > 2 && typing.IsString(parsed[2]) {
		requestedMode := typing.ToString(parsed[2])
		if isValidRoundingMode(requestedMode) {
			mode = requestedMode
		}
	}

	// Calculate multiplier for the given precision
	multiplier := math.Pow(10, precision)
	scaled := num * multiplier

	var rounded float64
	switch mode {
	case ROUND_DOWN:
		if scaled < 0 {
			rounded = math.Ceil(scaled)
		} else {
			rounded = math.Floor(scaled)
		}

	case ROUND_HALF_UP:
		if scaled < 0 {
			rounded = math.Ceil(scaled - 0.5)
		} else {
			rounded = math.Floor(scaled + 0.5)
		}

	case ROUND_HALF_EVEN:
		fraction := math.Abs(scaled - math.Floor(scaled))
		if fraction == 0.5 {
			// If the number is exactly halfway, round to the nearest even number
			floor := math.Floor(scaled)
			if math.Mod(floor, 2) == 0 {
				// If floor is even, round down
				rounded = floor
			} else {
				// If floor is odd, round up
				rounded = math.Ceil(scaled)
			}
		} else if fraction > 0.5 {
			rounded = math.Ceil(scaled)
		} else {
			rounded = math.Floor(scaled)
		}

	case ROUND_CEILING:
		rounded = math.Ceil(scaled)

	case ROUND_FLOOR:
		rounded = math.Floor(scaled)

	case ROUND_UP:
		if scaled < 0 {
			rounded = math.Floor(scaled)
		} else {
			rounded = math.Ceil(scaled)
		}

	case ROUND_HALF_DOWN:
		if scaled < 0 {
			// For negative numbers
			fraction := math.Abs(scaled - math.Ceil(scaled))
			if fraction <= 0.5 {
				// If exactly half or less, round towards zero (ceil for negative)
				rounded = math.Ceil(scaled)
			} else {
				// If more than half, round away from zero (floor for negative)
				rounded = math.Floor(scaled)
			}
		} else {
			// For positive numbers
			fraction := scaled - math.Floor(scaled)
			if fraction <= 0.5 {
				// If exactly half or less, round towards zero (floor for positive)
				rounded = math.Floor(scaled)
			} else {
				// If more than half, round away from zero (ceil for positive)
				rounded = math.Ceil(scaled)
			}
		}

	case ROUND_05UP:
		fraction := math.Abs(scaled - math.Floor(scaled))
		if fraction == 0.5 || fraction == 0.0 {
			if scaled < 0 {
				rounded = math.Floor(scaled)
			} else {
				rounded = math.Ceil(scaled)
			}
		} else {
			rounded = math.Floor(scaled)
		}
	}

	// Return to original scale
	return rounded / multiplier
}
