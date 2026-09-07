package lua

func assignable(value expression) bool {
	switch value.(type) {
	case nameExpression, memberExpression:
		return true
	default:
		return false
	}
}
