package dreefile

import "strings"

// formatCodeSection normalizes a code section (server, style, client) to the
// canonical three-part form — root tag, one-level-indented body, closing tag —
// while leaving an unclosed or single-line section untouched.
func formatCodeSection(raw string) string {
	if !strings.Contains(raw, "\n") {
		return raw
	}
	openEnd := strings.IndexByte(raw, '>')
	if openEnd < 0 {
		return raw
	}
	prefix := raw[:openEnd+1]
	rest := raw[openEnd+1:]
	idx := strings.LastIndex(rest, "</")
	if idx < 0 || !strings.HasSuffix(rest, ">") {
		return raw
	}
	closeTag := strings.TrimSpace(rest[idx:])
	inner := strings.TrimPrefix(rest[:idx], "\n")
	inner = strings.TrimRight(inner, " \t\r\n")
	return strings.TrimRight(indentCodeSection(prefix+"\n"+inner+"\n"+closeTag), " \t\r")
}

// indentCodeSection shifts a code section so its body sits one level under the
// root tag while the code's own relative indentation is preserved. A section
// that contains a multi-line string literal is left untouched, because shifting
// such lines would change the literal's value.
func indentCodeSection(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) < 2 || strings.IndexByte(lines[0], '>') < 0 {
		return text
	}
	last := len(lines)
	if strings.HasPrefix(strings.TrimSpace(lines[last-1]), "</") {
		last--
	}
	inner := lines[1:last]
	if hasMultilineLiteral(inner) {
		return text
	}
	min := -1
	for _, line := range inner {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if w := leadingWS(line); min < 0 || w < min {
			min = w
		}
	}
	if min < 0 {
		min = 0
	}

	var b strings.Builder
	b.WriteString(lines[0])
	b.WriteString("\n")
	for _, line := range inner {
		if strings.TrimSpace(line) == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString(indentUnit)
		b.WriteString(line[min:])
		b.WriteString("\n")
	}
	if last < len(lines) {
		b.WriteString(lines[last])
	}
	return b.String()
}

func leadingWS(line string) int {
	n := 0
	for n < len(line) && (line[n] == ' ' || line[n] == '\t') {
		n++
	}
	return n
}

// hasMultilineLiteral reports whether any string literal stays open across a
// line break. Multi-line literals (raw or escaped continuations) make the
// section unsafe to re-indent.
func hasMultilineLiteral(lines []string) bool {
	var quote byte
	escaped := false
	for _, line := range lines {
		for i := 0; i < len(line); i++ {
			c := line[i]
			if quote == 0 {
				if c == '"' || c == '\'' || c == '`' {
					quote = c
				}
				continue
			}
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
		}
		if quote != 0 {
			return true
		}
		escaped = false
	}
	return false
}
