package jsonlogic

// node is an operator call that Prepare unpacked from its single-key map. The operator
// is still looked up by name when evaluated, so AddOperator affects prepared rules.
type node struct {
	operator string
	values   any
}

// stringSet is a list of strings that Prepare indexed for "in".
type stringSet struct {
	list  []any
	items map[string]struct{}
}

func newStringSet(value any) (*stringSet, bool) {
	list, ok := value.([]any)
	if !ok || len(list) == 0 {
		return nil, false
	}

	items := make(map[string]struct{}, len(list))
	for _, item := range list {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		items[s] = struct{}{}
	}

	return &stringSet{list: list, items: items}, true
}

func (s *stringSet) contains(item string) bool {
	_, found := s.items[item]
	return found
}

// optimize replaces the operator calls of a rule with nodes. It only goes into values
// that are always evaluated: anything that can be returned or inspected as plain data,
// such as a map with many keys or the arguments of "filter", must stay as it is.
func optimize(rule any) any {
	switch r := rule.(type) {
	case map[string]any:
		if len(r) != 1 {
			return rule
		}
		for operator, values := range r {
			return &node{operator: operator, values: optimizeValues(operator, values)}
		}
	case []any:
		out := make([]any, len(r))
		for i, value := range r {
			out[i] = optimize(value)
		}
		return out
	}

	return rule
}

func optimizeValues(operator string, values any) any {
	list, isList := values.([]any)

	switch operator {
	case "and", "or", "if", "?:":
		// These require a list, and report the type of any other value.
		if !isList {
			return values
		}
	case "in":
		if isList && len(list) == 2 {
			if set, ok := newStringSet(list[1]); ok {
				return []any{optimize(list[0]), set}
			}
		}
	case "!", "!!", "==", "!=", "===", "!==", "<", "<=", ">", ">=",
		"+", "-", "*", "/", "%", "min", "max", "abs",
		"var", "missing", "missing_some", "set", "cat", "substr", "merge",
		"contains_all", "contains_any", "contains_none":
	default:
		return values
	}

	return optimize(values)
}
