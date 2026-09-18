package lexer

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func ParseFileHeader(input string) (*ir.FileHeader, string) {
	header, body, _ := ParseFileHeaderStrict(input)
	return header, body
}

func ParseFileHeaderStrict(input string) (*ir.FileHeader, string, error) {
	header := &ir.FileHeader{}
	lines := strings.Split(input, "\n")
	i := 0

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "DREEFILE" || strings.HasPrefix(trimmed, "DREEFILE ") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "DREEFILE"))
			if err := applyDreefileLine(header, rest); err != nil {
				return header, "", &HeaderError{Line: i + 1, Col: strings.Index(line, "DREEFILE") + 1, Err: err}
			}
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "LAYOUT ") {
			header.Layout = parseLayoutLine(trimmed)
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "COMPONENT ") {
			imp, consumed := parseComponentImport(lines[i:])
			if imp != nil {
				header.Imports = append(header.Imports, *imp)
				i += consumed
				continue
			}
		}

		if strings.HasPrefix(trimmed, "GOIMPORT") {
			paths, consumed := parseGoImport(lines[i:])
			if consumed > 0 {
				header.GoImports = append(header.GoImports, paths...)
				i += consumed
				continue
			}
		}

		if trimmed == "Component" || strings.HasPrefix(trimmed, "Component ") {
			return header, "", legacyHeaderError(line, i, "Component",
				"use DREEFILE component (props); the component name comes from the filename")
		}

		if trimmed == "import" || strings.HasPrefix(trimmed, "import ") {
			return header, "", legacyHeaderError(line, i, "import",
				`use GOIMPORT { path } for Go imports and COMPONENT "path" IMPORT { A } for components`)
		}

		if strings.HasPrefix(trimmed, "from ") {
			return header, "", legacyHeaderError(line, i, "from",
				`use COMPONENT "path" IMPORT { A }`)
		}

		if trimmed == "" {
			i++
			continue
		}

		break
	}

	return header, strings.Join(lines[i:], "\n"), nil
}

func ParseHeader(input string) (comp *ir.ComponentDef, imports []ir.Import, body string) {
	header, body := ParseFileHeader(input)
	if header.Kind == ir.FileKindComponent {
		comp = &ir.ComponentDef{Name: header.Name, Props: header.Props, Slots: header.Slots}
	}
	return comp, header.Imports, body
}

func parseFromImport(lines []string) (*ir.Import, int) {
	if len(lines) < 2 {
		return nil, 0
	}
	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "from ") {
		return nil, 0
	}
	rest := strings.TrimSpace(strings.TrimPrefix(line, "from "))
	marker := " import {"
	if !strings.HasSuffix(rest, marker) {
		return nil, 0
	}
	path := strings.TrimSpace(strings.TrimSuffix(rest, marker))
	if len(path) < 2 || path[0] != '"' || path[len(path)-1] != '"' {
		return nil, 0
	}
	imp := &ir.Import{Path: path[1 : len(path)-1]}
	for offset := 1; offset < len(lines); offset++ {
		name := strings.TrimSpace(lines[offset])
		if name == "}" {
			return imp, offset + 1
		}
		name = strings.TrimSuffix(name, ",")
		if name == "" || strings.ContainsAny(name, "{} \t") {
			return nil, 0
		}
		imp.Names = append(imp.Names, name)
	}
	return nil, 0
}

func parseComponentHeader(line string) *ir.ComponentDef {
	line = strings.TrimPrefix(line, "Component ")
	openParen := strings.IndexByte(line, '(')
	if openParen < 0 {
		return &ir.ComponentDef{Name: strings.TrimSpace(line)}
	}
	name := strings.TrimSpace(line[:openParen])
	rest := line[openParen:]

	closeParen := strings.IndexByte(rest, ')')
	if closeParen < 0 {
		return &ir.ComponentDef{Name: name}
	}

	comp := &ir.ComponentDef{Name: name}
	params := strings.TrimSpace(rest[1:closeParen])
	comp.Props = parseProps(params)

	slots := strings.TrimSpace(rest[closeParen+1:])
	if strings.HasPrefix(slots, "(") && strings.HasSuffix(slots, ")") {
		inner := strings.Trim(slots[1:len(slots)-1], " ")
		if inner != "" {
			for s := range strings.SplitSeq(inner, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					comp.Slots = append(comp.Slots, s)
				}
			}
		}
	}

	return comp
}

func parseProps(s string) []ir.Prop {
	var props []ir.Prop
	for part := range strings.SplitSeq(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}
		p := ir.Prop{Name: fields[0]}
		if len(fields) >= 2 {
			p.Type = fields[1]
		}
		if _, after, ok := strings.Cut(part, "="); ok {
			p.Default = strings.TrimSpace(after)
		}
		if p.Type == "" {
			p.Type = "string"
		}
		props = append(props, p)
	}
	return props
}

func parseImportLine(line string) *ir.Import {
	line = strings.TrimPrefix(line, "import ")
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	if len(fields) == 1 {
		path := strings.Trim(fields[0], "\"")
		if path == fields[0] {
			return nil
		}
		return &ir.Import{Path: path}
	}
	imp := &ir.Import{Path: strings.Trim(fields[len(fields)-1], "\"")}
	imp.Alias = fields[0]
	return imp
}
