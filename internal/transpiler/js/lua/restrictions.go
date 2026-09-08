package lua

var forbiddenCalls = map[string]bool{
	"collectgarbage": true,
	"dofile":         true,
	"load":           true,
	"loadfile":       true,
	"require":        true,
	"setmetatable":   true,
}

var forbiddenRoots = map[string]bool{
	"coroutine": true,
	"debug":     true,
	"io":        true,
	"os":        true,
	"package":   true,
	"socket":    true,
}

func memberRoot(value memberExpression) (string, bool) {
	current := value.object
	for {
		switch root := current.(type) {
		case nameExpression:
			return root.name, true
		case memberExpression:
			current = root.object
		default:
			return "", false
		}
	}
}
