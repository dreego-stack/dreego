package lexer

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func applyDreefileLine(h *ir.FileHeader, rest string) error {
	kind := ""
	if fields := strings.Fields(rest); len(fields) > 0 {
		kind = fields[0]
	}
	switch kind {
	case "page":
		h.Kind = ir.FileKindPage
	case "component":
		h.Kind = ir.FileKindComponent
		h.Props, h.Slots = parseSignature(rest)
	case "layout":
		h.Kind = ir.FileKindLayout
	default:
		return fmt.Errorf("invalid DREEFILE value %q: value must be one of page, component, layout", kind)
	}
	return nil
}

func parseLeadingQuoted(s string) (value string, end int, ok bool) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '"' {
		return "", 0, false
	}
	idx := strings.IndexByte(s[1:], '"')
	if idx < 0 {
		return "", 0, false
	}
	return s[1 : 1+idx], 1 + idx + 1, true
}

func parseSignature(rest string) ([]ir.Prop, []string) {
	open := strings.IndexByte(rest, '(')
	if open < 0 {
		return nil, nil
	}
	closeIdx := strings.IndexByte(rest[open:], ')')
	if closeIdx < 0 {
		return nil, nil
	}
	props := parseProps(strings.TrimSpace(rest[open+1 : open+closeIdx]))
	tail := strings.TrimSpace(rest[open+closeIdx+1:])
	return props, parseSlotList(tail)
}

func parseSlotList(tail string) []string {
	if !strings.HasPrefix(tail, "(") || !strings.HasSuffix(tail, ")") {
		return nil
	}
	var slots []string
	for s := range strings.SplitSeq(tail[1:len(tail)-1], ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			slots = append(slots, s)
		}
	}
	return slots
}

func parseQuotedValue(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return ""
}

func parseLayoutLine(line string) string {
	return parseQuotedValue(strings.TrimPrefix(line, "LAYOUT "))
}

func parseComponentImport(lines []string) (*ir.Import, int) {
	line := strings.TrimSpace(lines[0])
	rest := strings.TrimSpace(strings.TrimPrefix(line, "COMPONENT "))
	path, pathEnd, ok := parseLeadingQuoted(rest)
	if !ok || path == "" {
		return nil, 0
	}
	tail := strings.TrimSpace(rest[pathEnd:])
	if !strings.HasPrefix(tail, "IMPORT") {
		return nil, 0
	}
	tail = strings.TrimSpace(strings.TrimPrefix(tail, "IMPORT"))
	if !strings.HasPrefix(tail, "{") {
		return nil, 0
	}
	items, consumed := parseBraceList(append([]string{tail}, lines[1:]...))
	if consumed == 0 {
		return nil, 0
	}
	imp := &ir.Import{Path: path}
	var aliases map[string]string
	for _, item := range items {
		fields := strings.Fields(item)
		if len(fields) == 1 {
			imp.Names = append(imp.Names, fields[0])
			continue
		}
		if len(fields) == 3 && fields[1] == "as" {
			if aliases == nil {
				aliases = map[string]string{}
			}
			aliases[fields[2]] = fields[0]
			imp.Names = append(imp.Names, fields[2])
			continue
		}
		return nil, 0
	}
	imp.Aliases = aliases
	return imp, consumed
}

func parseGoImport(lines []string) ([]string, int) {
	line := strings.TrimSpace(lines[0])
	tail := strings.TrimSpace(strings.TrimPrefix(line, "GOIMPORT"))
	if !strings.HasPrefix(tail, "{") {
		return nil, 0
	}
	items, consumed := parseBraceList(append([]string{tail}, lines[1:]...))
	if consumed == 0 {
		return nil, 0
	}
	paths := make([]string, 0, len(items))
	for _, item := range items {
		path := strings.Trim(strings.TrimSpace(item), `"`)
		if path == "" {
			return nil, 0
		}
		paths = append(paths, path)
	}
	return paths, consumed
}

func parseBraceList(lines []string) (items []string, consumed int) {
	startLine, startCol := -1, 0
	for i := range lines {
		if c := strings.IndexByte(lines[i], '{'); c >= 0 {
			startLine, startCol = i, c
			break
		}
	}
	if startLine < 0 {
		return nil, 0
	}
	var b strings.Builder
	b.WriteString(lines[startLine][startCol+1:])
	i := startLine + 1
	for {
		s := b.String()
		if c := strings.IndexByte(s, '}'); c >= 0 {
			b.Reset()
			b.WriteString(s[:c])
			break
		}
		if i >= len(lines) {
			return nil, 0
		}
		b.WriteString("\n")
		b.WriteString(lines[i])
		i++
	}
	for part := range strings.SplitSeq(b.String(), ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items, i
}
