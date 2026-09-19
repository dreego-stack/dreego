package lexer

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func ParseFileHeader(input string) (*ir.FileHeader, string) {
	header, body, _ := ParseFileHeaderStrict(input)
	return header, body
}

func ParseFileHeaderStrict(input string) (*ir.FileHeader, string, error) {
	header := &ir.FileHeader{}
	lines := strings.Split(input, "\n")
	i := 0
	seenDreefile := false
	seenLayout := false

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "DREEFILE" || strings.HasPrefix(trimmed, "DREEFILE ") {
			if seenDreefile {
				return header, "", &HeaderError{Line: i + 1, Col: strings.Index(line, "DREEFILE") + 1,
					Err: fmt.Errorf("duplicate DREEFILE directive: declare the file kind exactly once")}
			}
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "DREEFILE"))
			if err := applyDreefileLine(header, rest); err != nil {
				return header, "", &HeaderError{Line: i + 1, Col: strings.Index(line, "DREEFILE") + 1, Err: err}
			}
			seenDreefile = true
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "LAYOUT ") {
			if seenLayout {
				return header, "", &HeaderError{Line: i + 1, Col: strings.Index(line, "LAYOUT") + 1,
					Err: fmt.Errorf("duplicate LAYOUT directive: declare the layout exactly once")}
			}
			path, err := parseLayoutLine(trimmed)
			if err != nil {
				return header, "", &HeaderError{Line: i + 1, Col: strings.Index(line, "LAYOUT") + 1, Err: err}
			}
			header.Layout = path
			seenLayout = true
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "COMPONENT ") {
			imp, consumed, err := parseComponentImport(lines[i:])
			if err != nil {
				return header, "", rebaseHeaderError(err, i)
			}
			if imp != nil {
				header.Imports = append(header.Imports, *imp)
				i += consumed
				continue
			}
		}

		if strings.HasPrefix(trimmed, "GOIMPORT") {
			paths, consumed, err := parseGoImport(lines[i:])
			if err != nil {
				return header, "", rebaseHeaderError(err, i)
			}
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

func rebaseHeaderError(err error, base int) error {
	he, ok := err.(*HeaderError)
	if !ok {
		return err
	}
	he.Line += base
	return he
}

func ParseHeader(input string) (comp *ir.ComponentDef, imports []ir.Import, body string) {
	header, body := ParseFileHeader(input)
	if header.Kind == ir.FileKindComponent {
		comp = &ir.ComponentDef{Name: header.Name, Props: header.Props, Slots: header.Slots}
	}
	return comp, header.Imports, body
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
