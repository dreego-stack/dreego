package parser

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func parseMessageExpression(raw string, pos int) (string, []ir.MessageArgument, error) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("message key is required at position %d", pos)
	}
	key := parts[0]
	if !validMessageKey(key) {
		return "", nil, fmt.Errorf("invalid message key %q at position %d", key, pos)
	}

	seen := make(map[string]bool, len(parts)-1)
	args := make([]ir.MessageArgument, 0, len(parts)-1)
	for _, part := range parts[1:] {
		name, expression, ok := strings.Cut(part, "=")
		if !ok || !validMessageName(name) {
			return "", nil, fmt.Errorf("invalid message argument %q at position %d", part, pos)
		}
		if expression == "" {
			return "", nil, fmt.Errorf("message argument %q requires an expression at position %d", name, pos)
		}
		if seen[name] {
			return "", nil, fmt.Errorf("duplicate message argument %q at position %d", name, pos)
		}
		seen[name] = true
		args = append(args, ir.MessageArgument{Name: name, Expression: expression})
	}
	return key, args, nil
}

func validMessageKey(key string) bool {
	for _, segment := range strings.Split(key, ".") {
		if !validMessageName(segment) {
			return false
		}
	}
	return true
}

func validMessageName(name string) bool {
	if name == "" || !isMessageLetter(name[0]) {
		return false
	}
	for i := 1; i < len(name); i++ {
		if !isMessageLetter(name[i]) && (name[i] < '0' || name[i] > '9') && name[i] != '-' && name[i] != '_' {
			return false
		}
	}
	return true
}

func isMessageLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}
