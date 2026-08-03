package core

import (
	"fmt"
	"strings"
)

func genHead(html string, bufName string) string {
	var out strings.Builder
	rest := html
	for rest != "" {
		open := strings.IndexByte(rest, '{')
		if open < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, goLiteral(rest)))
			break
		}
		if open > 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, goLiteral(rest[:open])))
			rest = rest[open:]
		}
		closeIdx := strings.IndexByte(rest, '}')
		if closeIdx < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, goLiteral(rest)))
			break
		}
		expr, filters := parseExpression(rest[1:closeIdx])
		code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
		raw := false
		for _, f := range filters {
			switch f {
			case "raw":
				raw = true
			case "upper":
				code = fmt.Sprintf("strings.ToUpper(%s)", code)
			}
		}
		if raw {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, code))
		} else {
			out.WriteString(fmt.Sprintf("%s.WriteString(html.EscapeString(%s))\n", bufName, code))
		}
		rest = rest[closeIdx+1:]
	}
	return out.String()
}
