package transpiler

import "strings"

func ParseHeader(input string) (comp *ComponentDef, imports []Import, body string) {
	lines := strings.Split(input, "\n")
	i := 0

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])

		if strings.HasPrefix(trimmed, "Component ") {
			comp = parseComponentHeader(trimmed)
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "import ") {
			imp := parseImportLine(trimmed)
			if imp != nil {
				imports = append(imports, *imp)
			}
			i++
			continue
		}

		if trimmed == "" {
			i++
			continue
		}

		break
	}

	body = strings.Join(lines[i:], "\n")
	return
}

func parseComponentHeader(line string) *ComponentDef {
	line = strings.TrimPrefix(line, "Component ")
	openParen := strings.IndexByte(line, '(')
	if openParen < 0 {
		return &ComponentDef{Name: strings.TrimSpace(line)}
	}
	name := strings.TrimSpace(line[:openParen])
	rest := line[openParen:]

	closeParen := strings.IndexByte(rest, ')')
	if closeParen < 0 {
		return &ComponentDef{Name: name}
	}

	comp := &ComponentDef{Name: name}
	params := strings.TrimSpace(rest[1:closeParen])
	comp.Props = parseProps(params)

	slots := strings.TrimSpace(rest[closeParen+1:])
	if strings.HasPrefix(slots, "(") && strings.HasSuffix(slots, ")") {
		inner := strings.Trim(slots[1:len(slots)-1], " ")
		if inner != "" {
			for _, s := range strings.Split(inner, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					comp.Slots = append(comp.Slots, s)
				}
			}
		}
	}

	return comp
}

func parseProps(s string) []Prop {
	var props []Prop
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}
		p := Prop{Name: fields[0]}
		if len(fields) >= 2 {
			p.Type = fields[1]
		}
		if eq := strings.IndexByte(part, '='); eq >= 0 {
			p.Default = strings.TrimSpace(part[eq+1:])
		}
		if p.Type == "" {
			p.Type = "string"
		}
		props = append(props, p)
	}
	return props
}

func parseImportLine(line string) *Import {
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
		return &Import{Path: path}
	}
	imp := &Import{Path: strings.Trim(fields[len(fields)-1], "\"")}
	imp.Alias = fields[0]
	return imp
}
