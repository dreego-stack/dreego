package html

import (
	"fmt"
	"strings"
)

func expressionCode(expr string, filters []string, pos int) (string, bool, error) {
	code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
	raw := false
	for _, f := range filters {
		switch f {
		case "raw":
			raw = true
		case "upper":
			code = fmt.Sprintf("strings.ToUpper(%s)", code)
		default:
			return "", false, fmt.Errorf("unknown filter '%s' at position %d", f, pos)
		}
	}
	return code, raw, nil
}

// conditionCode turns a template condition into a Go boolean expression. Every
// condition is routed through dreego.Truthy so strings, numbers, and slices are
// valid too (empty string / zero / empty collection = false). An optional Go
// init statement before the first top-level ";" is preserved, because the
// truthiness wrapper must only apply to the final expression.
func conditionCode(cond string) string {
	cond = strings.TrimSpace(cond)
	if idx := topLevelSemicolon(cond); idx >= 0 {
		init := strings.TrimSpace(cond[:idx])
		return init + "; dreego.Truthy(" + strings.TrimSpace(cond[idx+1:]) + ")"
	}
	return "dreego.Truthy(" + cond + ")"
}

func topLevelSemicolon(s string) int {
	var quote byte
	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if quote == '`' {
				if c == '`' {
					quote = 0
				}
				continue
			}
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			quote = c
		case ';':
			return i
		}
	}
	return -1
}
