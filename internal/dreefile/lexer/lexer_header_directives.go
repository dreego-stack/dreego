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

func parseGoImport(lines []string) ([]ir.GoImport, int, error) {
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
	imports := make([]ir.GoImport, 0, len(items))
	for _, item := range items {
		imp, err := parseGoImportItem(item, kw)
		if err != nil {
			return nil, 0, err
		}
		imports = append(imports, imp)
	}
	return imports, consumed, nil
}

func parseGoImportItem(item string, kw int) (ir.GoImport, error) {
	badPath := func() error {
		return &HeaderError{Line: 1, Col: kw + 1,
			Err: fmt.Errorf("invalid GOIMPORT path %q: use path or alias \"path\" with a non-empty import path without spaces", item)}
	}
	fields := strings.Fields(item)
	switch len(fields) {
	case 1:
		path := unquoteGoImportPath(fields[0])
		if path == "" {
			return ir.GoImport{}, badPath()
		}
		return ir.GoImport{Path: path}, nil
	case 2:
		alias := fields[0]
		if !validGoImportAlias(alias) {
			return ir.GoImport{}, badPath()
		}
		if !isQuotedGoImportPath(fields[1]) {
			return ir.GoImport{}, badPath()
		}
		path := fields[1][1 : len(fields[1])-1]
		if path == "" {
			return ir.GoImport{}, badPath()
		}
		return ir.GoImport{Alias: alias, Path: path}, nil
	default:
		return ir.GoImport{}, badPath()
	}
}

func unquoteGoImportPath(field string) string {
	if isQuotedGoImportPath(field) {
		field = field[1 : len(field)-1]
	}
	if field == "" || strings.ContainsAny(field, " \t\r\n\"") {
		return ""
	}
	return field
}

func isQuotedGoImportPath(field string) bool {
	return len(field) >= 2 && field[0] == '"' && field[len(field)-1] == '"'
}

func validGoImportAlias(alias string) bool {
	if alias == "" {
		return false
	}
	for i, r := range alias {
		if r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			continue
		}
		if i > 0 && r >= '0' && r <= '9' {
			continue
		}
		return false
	}
	return true
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
