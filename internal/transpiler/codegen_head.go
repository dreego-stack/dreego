package transpiler

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
		closeIdx := findExprEnd(rest[2:])
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
			out.WriteString(fmt.Sprintf("%s.WriteString(dreego.%s(%s))\n", bufName, headSafeFunc(html, pos), code))
		}
		pos += closeIdx + 2
		rest = rest[closeIdx+2:]
	}
	return out.String(), nil
}

// headSafeFunc picks the safe-value rule for an expression placeholder in a
// <head> section. The head HTML is a string, so the same tag scanning used for
// component text sections applies: placeholders inside URL, script, or style
// attribute values get the corresponding rule; everything else is text.
func headSafeFunc(html string, i int) string {
	if i < 0 || i >= len(html) {
		return "SafeText"
	}
	tagStart := strings.LastIndex(html[:i], "<")
	if tagStart < 0 || html[tagStart+1] == '/' || html[tagStart+1] == '!' {
		return "SafeText"
	}
	tagEnd := tagEnd(html[tagStart:])
	if tagEnd < 0 {
		return "SafeText"
	}
	tagEnd += tagStart
	if isMetaRefresh(html[tagStart:tagEnd]) {
		return "SafeRefresh"
	}
	name := attrNameAt(html[tagStart:tagEnd], i-tagStart)
	if name == "" {
		return "SafeText"
	}
	return attrContext(name)
}

// isMetaRefresh reports whether a <meta> tag carries http-equiv="refresh", whose
// content attribute embeds a URL after "url=". Attribute names and the "=" are
// matched case-insensitively and tolerate whitespace, like an HTML parser.
func isMetaRefresh(tag string) bool {
	name := strings.TrimSpace(attrValue(tag, "http-equiv"))
	return strings.EqualFold(name, "refresh")
}

// attrValue returns the value of the first occurrence of an attribute in a tag
// string, or "" when absent. Matching is case-insensitive and tolerates
// whitespace around "=" and the attribute name, like an HTML parser.
func attrValue(tag string, attr string) string {
	name := strings.ToLower(attr)
	lower := strings.ToLower(tag)
	i := 0
	for {
		idx := strings.Index(lower[i:], name)
		if idx < 0 {
			return ""
		}
		start := i + idx
		if start > 0 && !isAttrSpace(lower[start-1]) {
			i = start + len(name)
			continue
		}
		pos := skipAttrSpace(lower, start+len(name))
		if pos >= len(tag) || tag[pos] != '=' {
			i = pos
			continue
		}
		pos = skipAttrSpace(tag, pos+1)
		if pos >= len(tag) {
			return ""
		}
		if tag[pos] == '"' || tag[pos] == '\'' {
			quote := tag[pos]
			pos++
			end := strings.IndexByte(tag[pos:], quote)
			if end < 0 {
				return ""
			}
			return tag[pos : pos+end]
		}
		end := pos
		for end < len(tag) && !isAttrSpace(tag[end]) && tag[end] != '>' {
			end++
		}
		return tag[pos:end]
	}
}

// isAttrSpace reports whether c is an ASCII whitespace byte.
func isAttrSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// skipAttrSpace returns the index of the first non-whitespace byte at or after
// i.
func skipAttrSpace(s string, i int) int {
	for i < len(s) && isAttrSpace(s[i]) {
		i++
	}
	return i
}
