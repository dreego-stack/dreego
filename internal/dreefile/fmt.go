package dreefile

import (
	"regexp"
	"sort"
	"strings"
)

var expressions = regexp.MustCompile(`\{\{([^}]*?)\}\}`)
var controlOpen = regexp.MustCompile(`\{#(\w+)(\s+[^}]*?)?\}`)
var controlClose = regexp.MustCompile(`\{/(\w+)\}`)
var multiBlank = regexp.MustCompile(`\n{3,}`)
var multiSpace = regexp.MustCompile(` {2,}`)

var fmtSectionOrder = map[string]int{"server": 0, "head": 1, "body": 2, "style": 3, "client": 4}

type fmtSection struct {
	tag    string
	text   string
	gap    string
	offset int
}

func Format(input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(input, "\n")
	var headerLines []string
	var remaining []string
	headerDone := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if !headerDone {
			switch {
			case isLegacyHeaderLine(trimmed):
				end := legacyHeaderEnd(lines, i)
				for j := i; j <= end; j++ {
					headerLines = append(headerLines, strings.TrimRight(lines[j], " \t\r"))
				}
				i = end
				continue
			case trimmed == "DREEFILE",
				strings.HasPrefix(trimmed, "DREEFILE "),
				strings.HasPrefix(trimmed, "COMPONENT "),
				strings.HasPrefix(trimmed, "GOIMPORT"):
				if end := directiveBlockEnd(lines, i); end > i {
					for j := i; j <= end; j++ {
						headerLines = append(headerLines, strings.TrimRight(lines[j], " \t\r"))
					}
					i = end
					continue
				}
				headerLines = append(headerLines, trimmed)
				continue
			case strings.HasPrefix(trimmed, "LAYOUT "):
				headerLines = append(headerLines, formatLayoutLine(trimmed))
				continue
			case trimmed == "":
				headerLines = append(headerLines, "")
				continue
			}
		}

		headerDone = true
		remaining = append(remaining, line)
	}

	body := strings.Join(remaining, "\n")
	body = formatSections(body)
	body = strings.TrimRight(body, " \t\r")
	body = strings.TrimLeft(body, "\n")
	body = multiBlank.ReplaceAllString(body, "\n\n")

	var header strings.Builder
	for _, h := range headerLines {
		header.WriteString(h)
		header.WriteString("\n")
	}
	headerText := multiBlank.ReplaceAllString(header.String(), "\n\n")
	headerText = strings.Trim(headerText, "\n")

	var result strings.Builder
	if headerText != "" {
		result.WriteString(headerText)
		result.WriteString("\n\n")
	}
	result.WriteString(body)

	return strings.TrimRight(result.String(), " \t\r\n") + "\n"
}

// isLegacyHeaderLine reports whether trimmed is one of the removed header forms
// (Component ..., bare import ..., from ... import {...}). Format leaves these
// lines untouched instead of normalizing them, so a legacy file is never
// silently reformatted into something the generate error cannot explain.
func isLegacyHeaderLine(trimmed string) bool {
	return trimmed == "Component" || strings.HasPrefix(trimmed, "Component ") ||
		trimmed == "import" || strings.HasPrefix(trimmed, "import ") ||
		strings.HasPrefix(trimmed, "from ")
}

// legacyHeaderEnd returns the last line index belonging to a legacy header
// block. A single-line form returns start; a legacy from-import brace block is
// kept whole so its names are not moved into the body. A block that never
// closes or that reaches a section line falls back to the start line.
func legacyHeaderEnd(lines []string, start int) int {
	if !strings.Contains(lines[start], "{") {
		return start
	}
	depth := 0
	for j := start; j < len(lines); j++ {
		if j > start && strings.HasPrefix(strings.TrimSpace(lines[j]), "<") {
			return start
		}
		depth += braceDelta(lines[j])
		if depth <= 0 {
			return j
		}
	}
	return start
}

// directiveBlockEnd returns the index of the line that closes a brace block
// opened by the directive at start, or -1 when the first line has no opening
// brace, the block never closes, or a section line is reached first. A
// non-negative result is always > start only for block directives; inline
// directives return start. Callers fall back to a single header line when the
// block does not close, so an unbalanced directive can never swallow the body.
func directiveBlockEnd(lines []string, start int) int {
	if !strings.Contains(lines[start], "{") {
		return -1
	}
	depth := 0
	for j := start; j < len(lines); j++ {
		if j > start && strings.HasPrefix(strings.TrimSpace(lines[j]), "<") {
			return -1
		}
		depth += braceDelta(lines[j])
		if depth <= 0 {
			return j
		}
	}
	return -1
}

func braceDelta(line string) int {
	return strings.Count(line, "{") - strings.Count(line, "}")
}

func formatLayoutLine(line string) string {
	path := strings.TrimSpace(strings.TrimPrefix(line, "LAYOUT "))
	if len(path) >= 2 && path[0] == '"' && path[len(path)-1] == '"' {
		return `LAYOUT "` + path[1:len(path)-1] + `"`
	}
	return line
}

func formatExpressions(input string) string {
	return expressions.ReplaceAllStringFunc(input, func(m string) string {
		inner := m[2 : len(m)-2]
		inner = strings.TrimSpace(inner)
		inner = multiSpace.ReplaceAllString(inner, " ")
		inner = strings.ReplaceAll(inner, " |", "|")
		inner = strings.ReplaceAll(inner, "| ", "|")
		return "{{ " + inner + " }}"
	})
}

func formatControlFlow(input string) string {
	input = controlOpen.ReplaceAllStringFunc(input, func(m string) string {
		m = multiSpace.ReplaceAllString(m, " ")
		return m
	})
	input = controlClose.ReplaceAllStringFunc(input, func(m string) string {
		m = multiSpace.ReplaceAllString(m, " ")
		return m
	})
	return input
}

// formatSections normalizes a section file without inventing or dropping
// content. Section boundaries come from the real lexer, so a document-level
// <head> nested inside a body-level <body> layout stays part of its outer
// <body> instead of being hoisted to a root section. Text before, between, and
// after sections is preserved, and a file the lexer rejects is left unchanged.
func formatSections(input string) string {
	sections := collectFmtSections(input)
	if len(sections) == 0 {
		return input
	}

	reorder := true
	for _, s := range sections[1:] {
		if strings.TrimSpace(s.gap) != "" {
			reorder = false
			break
		}
	}
	ordered := sections
	if reorder {
		ordered = append([]fmtSection(nil), sections...)
		sort.SliceStable(ordered, func(i, j int) bool { return fmtSectionOrder[ordered[i].tag] < fmtSectionOrder[ordered[j].tag] })
	}

	var result strings.Builder
	if strings.TrimSpace(sections[0].gap) != "" {
		result.WriteString(sections[0].gap)
	}
	for i, s := range ordered {
		if i > 0 {
			if !reorder && strings.TrimSpace(s.gap) != "" {
				result.WriteString(s.gap)
			} else {
				result.WriteString("\n\n")
			}
		}
		text := formatSectionBody(s.tag, s.text)
		if s.tag == "body" {
			text = formatControlFlow(formatExpressions(text))
		}
		result.WriteString(text)
	}
	if tail := input[sections[len(sections)-1].offset:]; strings.TrimSpace(tail) != "" {
		result.WriteString(tail)
	}
	return result.String()
}

// collectFmtSections returns each root section with its exact source text, the
// verbatim gap before it, and its end offset. Section nesting comes from the
// lexer, so a section-like tag inside another section stays nested.
func collectFmtSections(input string) []fmtSection {
	toks, err := Lex(input)
	if err != nil {
		return nil
	}
	var out []fmtSection
	prev, depth, start, tag := 0, 0, -1, ""
	for _, tok := range toks {
		switch {
		case tok.Type == TokenTagOpen && depth == 0 && isFmtSection(tok.Tag):
			start, tag, depth = tok.Pos, tok.Tag, 1
		case depth > 0 && tok.Type == TokenTagOpen && tok.Tag == tag:
			depth++
		case depth > 0 && tok.Type == TokenTagClose && tok.Tag == tag:
			depth--
			if depth == 0 {
				end := tok.Pos + len("</"+tag+">")
				out = append(out, fmtSection{tag: tag, text: input[start:end], gap: input[prev:start], offset: end})
				prev = end
			}
		}
	}
	return out
}

func isFmtSection(tag string) bool {
	_, ok := fmtSectionOrder[tag]
	return ok
}

func formatSectionBody(tag, raw string) string {
	openEnd := strings.IndexByte(raw, '>')
	if openEnd < 0 {
		return raw
	}
	prefix := raw[:openEnd+1]
	suffix := "</" + tag + ">"
	raw = strings.TrimPrefix(raw, prefix)
	raw = strings.TrimSuffix(raw, suffix)

	lines := strings.Split(raw, "\n")
	start := 0
	end := len(lines) - 1
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	for end >= 0 && strings.TrimSpace(lines[end]) == "" {
		end--
	}

	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString("\n")
	for i := start; i <= end; i++ {
		line := strings.TrimRight(lines[i], " \t")
		if strings.TrimSpace(line) == "" {
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	sb.WriteString(suffix)
	return sb.String()
}
