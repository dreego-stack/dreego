package core

import (
	"fmt"
	"strings"
)

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
