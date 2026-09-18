package transpiler

import (
	"regexp"
	"strings"
)

var expressions = regexp.MustCompile(`\{\{([^}]*?)\}\}`)
var controlOpen = regexp.MustCompile(`\{#(\w+)(\s+[^}]*?)?\}`)
var controlClose = regexp.MustCompile(`\{/(\w+)\}`)
var multiBlank = regexp.MustCompile(`\n{3,}`)
var multiSpace = regexp.MustCompile(` {2,}`)

var knownSections = []string{"server", "head", "body", "style", "client"}

var sectionPatterns = map[string]*regexp.Regexp{
	"server": regexp.MustCompile(`<server(?:\s[^>]*)?>[\s\S]*?</server>`),
	"head":   regexp.MustCompile(`<head(?:\s[^>]*)?>[\s\S]*?</head>`),
	"body":   regexp.MustCompile(`<body(?:\s[^>]*)?>[\s\S]*?</body>`),
	"style":  regexp.MustCompile(`<style(?:\s[^>]*)?>[\s\S]*?</style>`),
	"client": regexp.MustCompile(`<client(?:\s[^>]*)?>[\s\S]*?</client>`),
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
			case strings.HasPrefix(trimmed, "Component "):
				headerLines = append(headerLines, formatCompHeader(trimmed))
				continue
			case strings.HasPrefix(trimmed, "import "):
				headerLines = append(headerLines, formatImport(trimmed))
				continue
			case strings.HasPrefix(trimmed, "DREEFILE "),
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
	result.WriteString("\n")

	return result.String()
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

func formatCompHeader(line string) string {
	parts := strings.SplitN(line, "(", 2)
	if len(parts) != 2 {
		return line
	}
	name := strings.TrimSpace(strings.TrimPrefix(parts[0], "Component "))
	params := strings.TrimRight(parts[1], ")")
	params = strings.TrimSpace(params)

	if params == "" {
		return "Component " + name
	}

	var formatted []string
	for p := range strings.SplitSeq(params, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		fields := strings.Fields(p)
		if len(fields) == 0 {
			continue
		}
		formatted = append(formatted, strings.Join(fields, " "))
	}
	return "Component " + name + " (" + strings.Join(formatted, ", ") + ")"
}

func formatImport(line string) string {
	line = strings.TrimPrefix(line, "import ")
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "import"
	}
	if len(fields) == 1 {
		return "import " + fields[0]
	}
	return "import " + fields[0] + " " + fields[1]
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

func formatSections(input string) string {
	found := map[string]string{}
	for _, tag := range knownSections {
		re := sectionPatterns[tag]
		m := re.FindString(input)
		if m != "" {
			found[tag] = formatSectionBody(tag, strings.TrimSpace(m))
		}
	}
	if len(found) == 0 {
		return input
	}

	if body, ok := found["body"]; ok {
		body = formatExpressions(body)
		body = formatControlFlow(body)
		found["body"] = body
	}

	var result strings.Builder
	for i, tag := range knownSections {
		body, ok := found[tag]
		if !ok {
			continue
		}
		if i > 0 && result.Len() > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(body)
	}
	return result.String()
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
