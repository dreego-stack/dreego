package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func goLiteral(s string) string {
	if strings.Contains(s, "`") {
		return strconv.Quote(s)
	}
	return "`" + s + "`"
}

func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	var result strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		result.WriteByte(byte(unicode.ToUpper(rune(p[0]))))
		if len(p) > 1 {
			result.WriteString(strings.ToLower(p[1:]))
		}
	}
	return result.String()
}

func scopeCSS(css string, hash string) string {
	prefix := fmt.Sprintf("[data-scope=%s] ", hash)
	var result strings.Builder
	depth := 0
	atDepth := 0
	selectorStart := 0

	for i := 0; i < len(css); i++ {
		ch := css[i]

		if ch == '{' {
			selector := strings.TrimSpace(css[selectorStart:i])
			if depth == 0 {
				if strings.HasPrefix(selector, "@") {
					atDepth = 1
					result.WriteString(selector)
				} else {
					result.WriteString(prefix)
					result.WriteString(selector)
				}
			} else if depth == atDepth {
				result.WriteString(prefix)
				result.WriteString(selector)
			} else {
				result.WriteString(selector)
			}
			result.WriteByte('{')
			depth++
			selectorStart = i + 1
			continue
		}

		if ch == '}' {
			result.WriteByte('}')
			depth--
			if depth < atDepth {
				atDepth = 0
			}
			selectorStart = i + 1
			continue
		}
	}

	trailing := strings.TrimSpace(css[selectorStart:])
	if trailing != "" {
		result.WriteString(trailing)
	}

	return strings.TrimSpace(result.String())
}

func extractAttrValues(attrs string) string {
	if attrs == "" {
		return ""
	}
	var vals []string
	inQuote := false
	braceDepth := 0
	start := 0
	for i := 0; i < len(attrs); i++ {
		ch := attrs[i]
		if ch == '"' && braceDepth == 0 {
			inQuote = !inQuote
		}
		if ch == '{' && !inQuote {
			braceDepth++
		}
		if ch == '}' && !inQuote {
			braceDepth--
		}
		if ch == ' ' && !inQuote && braceDepth == 0 {
			if start < i {
				vals = append(vals, attrVal(attrs[start:i]))
			}
			start = i + 1
		}
	}
	if start < len(attrs) {
		vals = append(vals, attrVal(attrs[start:]))
	}
	return strings.Join(vals, ", ")
}

func attrVal(part string) string {
	eq := strings.IndexByte(part, '=')
	if eq < 0 {
		return fmt.Sprintf("%q", part)
	}
	val := strings.TrimSpace(part[eq+1:])
	if val == "" {
		return fmt.Sprintf("%q", "")
	}
	if val[0] == '{' && val[len(val)-1] == '}' {
		return val[1 : len(val)-1]
	}
	val = strings.Trim(val, "\"")
	return fmt.Sprintf("%q", val)
}
