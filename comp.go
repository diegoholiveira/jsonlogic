package jsonlogic

import (
	"reflect"
)

func less(a, b interface{}) bool {
	if isNumber(a) || isNumber(b) {
		return toNumber(b) > toNumber(a)
	}

	// comparison to a nil value is falsy
	if a == nil || b == nil {
		return false
	}

	return toString(b) > toString(a)
}

func hardEquals(a, b interface{}) bool {
	ra := reflect.ValueOf(a).Kind()
	rb := reflect.ValueOf(b).Kind()

	if ra != rb {
		return false
	}

	return equals(a, b)
}

func equals(a, b interface{}) bool {
	if isNumber(a) {
		return toNumber(a) == toNumber(b)
	}

	// comparison to a nil value is falsy
	if a == nil || b == nil {
		return false
	}

	if isBool(a) {
		return isTrue(a) == isTrue(b)
	}

	return toString(a) == toString(b)
}
