package core

import (
	"fmt"
	"strings"
)

func genHead(html string, bufName string) (string, error) {
	var out strings.Builder
	rest := html
	pos := 0
	for rest != "" {
		open := strings.Index(rest, "{{")
		if open < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, goLiteral(rest)))
			break
		}
		if open > 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, goLiteral(rest[:open])))
			pos += open
			rest = rest[open:]
		}
		closeIdx := strings.Index(rest[2:], "}}")
		if closeIdx < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, goLiteral(rest)))
			break
		}
		closeIdx += 2
		expr, filters := parseExpression(strings.TrimSpace(rest[2:closeIdx]))
		code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
		raw := false
		for _, f := range filters {
			switch f {
			case "raw":
				raw = true
			case "upper":
				code = fmt.Sprintf("strings.ToUpper(%s)", code)
			default:
				return "", fmt.Errorf("unknown filter '%s' at position %d", f, pos)
			}
		}
		if raw {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, code))
		} else {
			out.WriteString(fmt.Sprintf("%s.WriteString(html.EscapeString(%s))\n", bufName, code))
		}
		pos += closeIdx + 2
		rest = rest[closeIdx+2:]
	}
	return out.String(), nil
}
