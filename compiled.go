package jsonlogic

import (
	"strings"
	"sync/atomic"
)

// compiledChecksDisabled is set to 1 once AddOperator replaces an operator that compiled checks implement.
var compiledChecksDisabled int32

// check evaluates a prepared rule directly; ok is false when it cannot decide and apply must.
type check func(data map[string]any) (match, ok bool)

func (c check) run(data any) (match, ok bool) {
	m, isMap := data.(map[string]any)
	if c == nil || !isMap || atomic.LoadInt32(&compiledChecksDisabled) == 1 {
		return false, false
	}
	return c(m)
}

func disableCompiledChecksFor(op string) {
	switch op {
	case "and", "or", "!", "==", "!=", "===", "!==", "in":
		atomic.StoreInt32(&compiledChecksDisabled, 1)
	}
}

func compile(rule any) check {
	m, ok := rule.(map[string]any)
	if !ok || len(m) != 1 {
		return nil
	}
	for op, args := range m {
		switch op {
		case "and", "or":
			return compileLogical(op == "and", args)
		case "!":
			return compileNot(args)
		case "==", "!=", "===", "!==":
			return compileEquals(op[0] == '!', args)
		case "in":
			return compileIn(args)
		}
	}
	return nil
}

func compileLogical(isAnd bool, args any) check {
	list, ok := args.([]any)
	if !ok || len(list) == 0 {
		return nil
	}
	checks := make([]check, len(list))
	for i, a := range list {
		if checks[i] = compile(a); checks[i] == nil {
			return nil
		}
	}
	return func(data map[string]any) (bool, bool) {
		for _, c := range checks {
			match, ok := c(data)
			if !ok {
				return false, false
			}
			if match != isAnd {
				return match, true
			}
		}
		return isAnd, true
	}
}

func compileNot(args any) check {
	if list, ok := args.([]any); ok && len(list) == 1 {
		args = list[0]
	}
	c := compile(args)
	if c == nil {
		return nil
	}
	return func(data map[string]any) (bool, bool) {
		match, ok := c(data)
		return !match, ok
	}
}

func compileEquals(negate bool, args any) check {
	list, ok := args.([]any)
	if !ok || len(list) != 2 {
		return nil
	}
	name, ok := plainVar(list[0])
	if !ok {
		return nil
	}
	lit, ok := list[1].(string)
	if !ok {
		return nil
	}
	return func(data map[string]any) (bool, bool) {
		v, ok := data[name].(string)
		return (v == lit) != negate, ok
	}
}

func compileIn(args any) check {
	list, ok := args.([]any)
	if !ok || len(list) != 2 {
		return nil
	}
	name, ok := plainVar(list[0])
	if !ok {
		return nil
	}
	values, ok := list[1].([]any)
	if !ok {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		s, ok := v.(string)
		if !ok {
			return nil
		}
		set[s] = struct{}{}
	}
	return func(data map[string]any) (bool, bool) {
		v, ok := data[name].(string)
		_, found := set[v]
		return found, ok
	}
}

func plainVar(v any) (string, bool) {
	m, ok := v.(map[string]any)
	if !ok || len(m) != 1 {
		return "", false
	}
	name, ok := m["var"].(string)
	if !ok || name == "" {
		return "", false
	}
	return name, !strings.Contains(name, ".")
}
