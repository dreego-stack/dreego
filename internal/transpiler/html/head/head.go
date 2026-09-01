package head

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func Gen(html string, bufName string) (string, error) {
	var out strings.Builder
	rest := html
	pos := 0
	for rest != "" {
		open := strings.Index(rest, "{{")
		if open < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral(rest)))
			break
		}
		if open > 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral(rest[:open])))
			pos += open
			rest = rest[open:]
		}
		closeIdx := ir.FindExprEnd(rest[2:])
		if closeIdx < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral(rest)))
			break
		}
		closeIdx += 2
		expr, filters := ir.ParseExpression(strings.TrimSpace(rest[2:closeIdx]))
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
			out.WriteString(fmt.Sprintf("%s.WriteString(dreego.%s(%s))\n", bufName, HeadSafeFunc(html, pos), code))
		}
		pos += closeIdx + 2
		rest = rest[closeIdx+2:]
	}
	return out.String(), nil
}

func HeadSafeFunc(html string, i int) string {
	if i < 0 || i >= len(html) {
		return "SafeText"
	}
	tagStart := strings.LastIndex(html[:i], "<")
	if tagStart < 0 || html[tagStart+1] == '/' || html[tagStart+1] == '!' {
		return "SafeText"
	}
	tagEnd := ir.TagEnd(html[tagStart:])
	if tagEnd < 0 {
		return "SafeText"
	}
	tagEnd += tagStart
	if isMetaRefresh(html[tagStart:tagEnd]) {
		return "SafeRefresh"
	}
	name := ir.AttrNameAt(html[tagStart:tagEnd], i-tagStart)
	if name == "" {
		return "SafeText"
	}
	return ir.AttrContext(name)
}

func isMetaRefresh(tag string) bool {
	name := strings.TrimSpace(ir.AttrValue(tag, "http-equiv"))
	return strings.EqualFold(name, "refresh")
}
