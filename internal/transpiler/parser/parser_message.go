package parser

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func ParseMessageExpression(raw string, pos int) (string, []ir.MessageArgument, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil, fmt.Errorf("message key is required at position %d", pos)
	}
	keyEnd := strings.IndexAny(raw, " \t\r\n")
	if keyEnd < 0 {
		keyEnd = len(raw)
	}
	key := raw[:keyEnd]
	if !validMessageKey(key) {
		return "", nil, fmt.Errorf("invalid message key %q at position %d", key, pos)
	}

	seen := make(map[string]bool)
	var args []ir.MessageArgument
	rest := strings.TrimSpace(raw[keyEnd:])
	for rest != "" {
		equals := strings.IndexByte(rest, '=')
		space := strings.IndexAny(rest, " \t\r\n")
		if equals <= 0 || space >= 0 && space < equals {
			invalid := rest
			if space >= 0 {
				invalid = rest[:space]
			}
			return "", nil, fmt.Errorf("invalid message argument %q at position %d", invalid, pos)
		}
		name := rest[:equals]
		if !validMessageName(name) {
			return "", nil, fmt.Errorf("invalid message argument %q at position %d", name, pos)
		}
		expression, remaining := splitMessageArgument(rest[equals+1:])
		if expression == "" {
			return "", nil, fmt.Errorf("message argument %q requires an expression at position %d", name, pos)
		}
		if seen[name] {
			return "", nil, fmt.Errorf("duplicate message argument %q at position %d", name, pos)
		}
		seen[name] = true
		args = append(args, ir.MessageArgument{Name: name, Expression: expression})
		rest = strings.TrimSpace(remaining)
	}
	return key, args, nil
}

func splitMessageArgument(value string) (string, string) {
	depth := 0
	var quote byte
	escaped := false
	for index := 0; index < len(value); index++ {
		char := value[index]
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if char == '\\' && quote != '`' {
				escaped = true
				continue
			}
			if char == quote {
				quote = 0
			}
			continue
		}
		switch char {
		case '\'', '"', '`':
			quote = char
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		default:
			if depth == 0 && isMessageSpace(char) {
				next := strings.TrimLeft(value[index:], " \t\r\n")
				equals := strings.IndexByte(next, '=')
				space := strings.IndexAny(next, " \t\r\n")
				if equals > 0 && (space < 0 || equals < space) && validMessageName(next[:equals]) {
					return strings.TrimSpace(value[:index]), next
				}
			}
		}
	}
	return strings.TrimSpace(value), ""
}

func isMessageSpace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\r' || char == '\n'
}

func parseMessageExpression(raw string, pos int) (string, []ir.MessageArgument, error) {
	return ParseMessageExpression(raw, pos)
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
