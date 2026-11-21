package jsonlogic

import "strings"

func substr(values, data any) any {
	values = parseValues(values, data)
	parsed := values.([]any)

	runes := []rune(toString(parsed[0]))
	length := len(runes)

	// Same as String.prototype.substr in json-logic-js: a negative start
	// counts from the end (clamped to 0) and a start past the end gives "".
	from := int(toNumber(parsed[1]))
	if from < 0 {
		from += length
		if from < 0 {
			from = 0
		}
	}
	if from > length {
		from = length
	}

	to := length
	if len(parsed) == 3 {
		n := int(toNumber(parsed[2]))
		if n < 0 {
			// a negative length drops that many characters from the end
			to = length + n
		} else if n < length-from {
			to = from + n
		}
	}

	if to <= from {
		return ""
	}

	return string(runes[from:to])
}

func concat(values, data any) any {
	values = parseValues(values, data)
	if _, ok := values.(string); ok {
		return values
	}

	inputSlice := values.([]any)

	if len(inputSlice) == 0 {
		return ""
	}

	if len(inputSlice) == 1 {
		return toString(inputSlice[0])
	}

	var s strings.Builder

	for _, text := range inputSlice {
		s.WriteString(toString(text))
	}

	return s.String()
}
