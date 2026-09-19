package lua

func assignable(value expression) bool {
	switch value.(type) {
	case nameExpression, memberExpression, indexExpression:
		return true
	default:
		return false
	}
}
