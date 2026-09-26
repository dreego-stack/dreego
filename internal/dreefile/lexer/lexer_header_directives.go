package lexer

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
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

func parseLayoutLine(line string) (string, error) {
	raw := strings.TrimSpace(strings.TrimPrefix(line, "LAYOUT"))
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return "", fmt.Errorf("invalid LAYOUT value %q: path must be a quoted string", raw)
	}
	path := raw[1 : len(raw)-1]
	if path == "" {
		return "", fmt.Errorf("invalid LAYOUT value: path must not be empty")
	}
	return path, nil
}

func parseProfileLine(line string) (string, error) {
	raw := strings.TrimSpace(strings.TrimPrefix(line, "PROFILE"))
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return "", fmt.Errorf("invalid PROFILE value %q: name must be a quoted string", raw)
	}
	name := raw[1 : len(raw)-1]
	if name == "" {
		return "", fmt.Errorf("invalid PROFILE value: name must not be empty")
	}
	return name, nil
}

func parseComponentImport(lines []string) (*ir.Import, int, error) {
	line := lines[0]
	kw := strings.Index(line, "COMPONENT")
	if kw < 0 {
		return nil, 0, nil
	}
	raw := line[kw+len("COMPONENT"):]
	lead := len(raw) - len(strings.TrimLeft(raw, " \t"))
	rest := raw[lead:]
	restOffset := kw + len("COMPONENT") + lead
	path, pathEnd, ok := parseLeadingQuoted(rest)
	if !ok || path == "" {
		return nil, 0, nil
	}
	tail := rest[pathEnd:]
	impIdx := strings.Index(tail, "IMPORT")
	if impIdx < 0 {
		return nil, 0, nil
	}
	afterImport := strings.TrimSpace(tail[impIdx+len("IMPORT"):])
	if !strings.HasPrefix(afterImport, "{") {
		return nil, 0, nil
	}
	from := restOffset + pathEnd + impIdx + len("IMPORT")
	items, consumed, err := parseBraceList(lines, from, "COMPONENT IMPORT")
	if err != nil || consumed == 0 {
		return nil, consumed, err
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
		return nil, 0, &HeaderError{Line: 1, Col: restOffset + 1,
			Err: fmt.Errorf("invalid COMPONENT IMPORT entry %q: use Name or Name as Alias", item)}
	}
	imp.Aliases = aliases
	return imp, consumed, nil
}

func parseGoImport(lines []string) ([]string, int, error) {
	line := lines[0]
	kw := strings.Index(line, "GOIMPORT")
	if kw < 0 {
		return nil, 0, nil
	}
	from := kw + len("GOIMPORT")
	items, consumed, err := parseBraceList(lines, from, "GOIMPORT")
	if err != nil || consumed == 0 {
		return nil, consumed, err
	}
	if len(items) == 0 {
		return nil, 0, &HeaderError{Line: 1, Col: kw + 1,
			Err: fmt.Errorf("invalid GOIMPORT value: directive requires at least one path")}
	}
	paths := make([]string, 0, len(items))
	for _, item := range items {
		path := strings.Trim(strings.TrimSpace(item), `"`)
		if path == "" || strings.ContainsAny(path, " \t\r\n") {
			return nil, 0, &HeaderError{Line: 1, Col: kw + 1,
				Err: fmt.Errorf("invalid GOIMPORT path %q: path must be a non-empty import path without spaces", item)}
		}
		paths = append(paths, path)
	}
	return paths, consumed, nil
}

func parseBraceList(lines []string, from int, directive string) ([]string, int, error) {
	first := lines[0]
	if from > len(first) {
		return nil, 0, nil
	}
	open := strings.IndexByte(first[from:], '{')
	if open < 0 {
		return nil, 0, nil
	}
	open += from
	var parts []string
	var current strings.Builder
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		j := 0
		if i == 0 {
			j = open + 1
		}
		for ; j < len(line); j++ {
			switch line[j] {
			case '{':
				return nil, 0, &HeaderError{Line: i + 1, Col: j + 1,
					Err: fmt.Errorf("%s list: nested '{' is not allowed", directive)}
			case '}':
				if strings.TrimSpace(line[j+1:]) != "" {
					return nil, 0, &HeaderError{Line: i + 1, Col: j + 2,
						Err: fmt.Errorf("%s list: unexpected content after closing '}'", directive)}
				}
				parts = append(parts, current.String())
				return braceItems(parts), i + 1, nil
			case ',':
				parts = append(parts, current.String())
				current.Reset()
			default:
				current.WriteByte(line[j])
			}
		}
		if i > 0 {
			current.WriteByte('\n')
		}
	}
	return nil, 0, &HeaderError{Line: 1, Col: open + 1,
		Err: fmt.Errorf("%s list: unbalanced '{', missing closing '}'", directive)}
}

func braceItems(parts []string) []string {
	var items []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}
