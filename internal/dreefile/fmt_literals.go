package dreefile

import "strings"

// applyOutsideLiterals calls f on every span of s that is not inside a
// double-quoted, single-quoted, or back-quoted literal, and copies literal
// spans through unchanged. It keeps formatter normalization from rewriting
// string contents such as {{ "a  b" }} or {#if x == "a | b"}.
func applyOutsideLiterals(s string, f func(string) string) string {
	var b strings.Builder
	segStart := 0
	i := 0
	for i < len(s) {
		c := s[i]
		if c != '"' && c != '\'' && c != '`' {
			i++
			continue
		}
		b.WriteString(f(s[segStart:i]))
		j := literalEnd(s, i)
		b.WriteString(s[i:j])
		i = j
		segStart = j
	}
	b.WriteString(f(s[segStart:]))
	return b.String()
}

// literalEnd returns the index just past the literal that starts at i.
func literalEnd(s string, i int) int {
	quote := s[i]
	j := i + 1
	for j < len(s) {
		switch {
		case quote == '`':
			if s[j] == '`' {
				return j + 1
			}
		case s[j] == '\\' && j+1 < len(s):
			j++
		case s[j] == quote:
			return j + 1
		}
		j++
	}
	return len(s)
}

func normalizeExpressionInner(inner string) string {
	inner = strings.TrimSpace(inner)
	return applyOutsideLiterals(inner, func(seg string) string {
		seg = multiSpace.ReplaceAllString(seg, " ")
		seg = strings.ReplaceAll(seg, " |", "|")
		seg = strings.ReplaceAll(seg, "| ", "|")
		return seg
	})
}

func normalizeControlTag(m string) string {
	return applyOutsideLiterals(m, func(seg string) string {
		return multiSpace.ReplaceAllString(seg, " ")
	})
}

// isCodeSection reports whether a section's body is literal code or text that
// the formatter must preserve byte for byte.
func isCodeSection(tag string) bool {
	switch tag {
	case "server", "client", "style":
		return true
	}
	return false
}

func collapseBlankLines(s string) string {
	return multiBlank.ReplaceAllString(s, "\n\n")
}
