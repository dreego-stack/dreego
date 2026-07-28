package core

import (
	"regexp"
	"strings"

)

var expressions = regexp.MustCompile(`\{([^!#/][^}]*?)\}`)
var controlOpen = regexp.MustCompile(`\{#(\w+)(\s+[^}]*?)?\}`)
var controlClose = regexp.MustCompile(`\{/(\w+)\}`)
var multiBlank = regexp.MustCompile(`\n{3,}`)
var multiSpace = regexp.MustCompile(` {2,}`)

var knownSections = []string{"head", "go", "div", "script", "style"}

var sectionPatterns = map[string]*regexp.Regexp{
	"head":   regexp.MustCompile(`<head>[\s\S]*?</head>`),
	"go":     regexp.MustCompile(`<go>[\s\S]*?</go>`),
	"div":    regexp.MustCompile(`<div>[\s\S]*?</div>`),
	"script": regexp.MustCompile(`<script>[\s\S]*?</script>`),
	"style":  regexp.MustCompile(`<style>[\s\S]*?</style>`),
}

func Format(input string) string {
	lines := strings.Split(input, "\n")
	var headerLines []string
	var remaining []string
	headerDone := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !headerDone && strings.HasPrefix(trimmed, "Component ") {
			headerLines = append(headerLines, formatCompHeader(trimmed))
			continue
		}
		if !headerDone && strings.HasPrefix(trimmed, "import ") {
			headerLines = append(headerLines, formatImport(trimmed))
			continue
		}
		if !headerDone && trimmed == "" {
			headerLines = append(headerLines, "")
			continue
		}
		headerDone = true
		remaining = append(remaining, line)
	}

	body := strings.Join(remaining, "\n")
	body = formatSections(body)
	body = strings.TrimRight(body, " \t")
	body = strings.TrimLeft(body, "\n")
	body = multiBlank.ReplaceAllString(body, "\n\n")

	var result strings.Builder
	for _, h := range headerLines {
		result.WriteString(h)
		result.WriteString("\n")
	}
	if len(headerLines) > 0 {
		result.WriteString("\n")
	}
	result.WriteString(body)
	result.WriteString("\n")

	return result.String()
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
	for _, p := range strings.Split(params, ",") {
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
		inner := m[1 : len(m)-1]
		inner = strings.TrimSpace(inner)
		inner = multiSpace.ReplaceAllString(inner, " ")
		inner = strings.ReplaceAll(inner, " |", "|")
		inner = strings.ReplaceAll(inner, "| ", "|")
		return "{" + inner + "}"
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

	if body, ok := found["div"]; ok {
		body = formatExpressions(body)
		body = formatControlFlow(body)
		found["div"] = body
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
	prefix := "<" + tag + ">"
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
