package jsonlogic

func filter(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	var subject interface{}

	if isSlice(parsed[0]) {
		subject = parsed[0]
	}

	if isMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	result := make([]interface{}, 0)

	if subject == nil {
		return result
	}

	logic := solveVars(parsed[1], data)

	for _, value := range subject.([]interface{}) {
		v := parseValues(logic, value)

		if isTrue(v) {
			result = append(result, value)
		}
	}

	return result
}

func _map(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	var subject interface{}

	if isSlice(parsed[0]) {
		subject = parsed[0]
	}

	if isMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	result := make([]interface{}, 0)

	if subject == nil {
		return result
	}

	logic := solveVars(parsed[1], data)

	for _, value := range subject.([]interface{}) {
		v := parseValues(logic, value)

		if isTrue(v) || isNumber(v) {
			result = append(result, v)
		}
	}

	return result
}

func reduce(values, data interface{}) interface{} {
	parsed := values.([]interface{})

	var subject interface{}

	if isSlice(parsed[0]) {
		subject = parsed[0]
	}

	if isMap(parsed[0]) {
		subject = apply(parsed[0], data)
	}

	if subject == nil {
		return float64(0)
	}

	context := map[string]interface{}{
		"current":     float64(0),
		"accumulator": toNumber(parsed[2]),
	}

	for _, value := range subject.([]interface{}) {
		if value == nil {
			continue
		}

		context["current"] = value

		v := apply(parsed[1], context)

		context["accumulator"] = toNumber(v)
	}

	return context["accumulator"]
}
