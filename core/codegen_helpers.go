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
	scopeRange(&result, css, 0, len(css), prefix)
	return strings.TrimSpace(result.String())
}

func scopeRange(b *strings.Builder, css string, start, end int, prefix string) {
	i := start
	for i < end {
		if css[i] == '{' {
			sel := strings.TrimSpace(css[start:i])
			close := matchBrace(css, i, end)
			innerStart, innerEnd := i+1, close
			b.WriteString(scopedHeader(sel, prefix))
			b.WriteByte('{')
			if strings.HasPrefix(sel, "@keyframes") {
				b.WriteString(css[innerStart:innerEnd])
			} else if !strings.HasPrefix(sel, "@") && sel != "" {
				b.WriteString(css[innerStart:innerEnd])
			} else {
				scopeRange(b, css, innerStart, innerEnd, prefix)
			}
			b.WriteByte('}')
			start = close + 1
			i = close + 1
			continue
		}
		i++
	}
	if start < end {
		b.WriteString(css[start:end])
	}
}

func scopedHeader(sel string, prefix string) string {
	if sel == "" {
		return ""
	}
	if strings.HasPrefix(sel, "@") {
		return sel
	}
	return scopeSelector(sel, prefix)
}

func scopeSelector(sel string, prefix string) string {
	parts := splitTopLevelComma(sel)
	scoped := make([]string, len(parts))
	for i, p := range parts {
		scoped[i] = prefix + strings.TrimSpace(p)
	}
	return strings.Join(scoped, ",\n")
}

func splitTopLevelComma(sel string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(sel); i++ {
		switch sel[i] {
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				parts = append(parts, sel[start:i])
				start = i + 1
			}
		}
	}
	return append(parts, sel[start:])
}

func matchBrace(css string, open, end int) int {
	depth := 1
	for i := open + 1; i < end; i++ {
		switch css[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return end
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
