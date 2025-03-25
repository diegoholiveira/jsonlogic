package jsonlogic

import (
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3/internal/typing"
)

func substr(values, data any) any {
	values = parseValues(values, data)
	parsed := values.([]any)

	runes := []rune(typing.ToString(parsed[0]))

	from := int(typing.ToNumber(parsed[1]))
	length := len(runes)

	if from < 0 {
		from = length + from
	}

	if from < 0 || from > length {
		// case from is still negative, we must stop right now and return the original string
		return string(runes)
	}

	if len(parsed) == 3 {
		length = int(typing.ToNumber(parsed[2]))
	}

	var to int
	if length < 0 {
		length = len(runes) + length
		to = length
	} else {
		to = from + length
	}

	if to > len(runes) {
		to = len(runes)
	}

	return string(runes[from:to])
}

func concat(values, data any) any {
	values = parseValues(values, data)
	if typing.IsString(values) {
		return values
	}

	inputSlice := values.([]any)

	if len(inputSlice) == 0 {
		return ""
	}

	if len(inputSlice) == 1 {
		return typing.ToString(inputSlice[0])
	}

	var s strings.Builder

	for _, text := range inputSlice {
		s.WriteString(typing.ToString(text))
	}

	return strings.TrimSpace(s.String())
}
