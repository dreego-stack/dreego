package transpiler

import (
	"fmt"
	"strings"
)

type callAttr struct {
	Name  string
	Value string
}

func parseCallAttrs(attrs string) ([]callAttr, error) {
	if attrs == "" {
		return nil, nil
	}
	var out []callAttr
	inQuote := false
	var quote byte
	braceDepth := 0
	start := 0
	for i := 0; i < len(attrs); i++ {
		ch := attrs[i]
		switch ch {
		case '"', '\'':
			if braceDepth == 0 {
				if !inQuote {
					inQuote = true
					quote = ch
				} else if quote == ch {
					inQuote = false
				}
			}
		case '{':
			if !inQuote {
				braceDepth++
			}
		case '}':
			if !inQuote && braceDepth > 0 {
				braceDepth--
			}
		default:
			if isAttrSpace(ch) && !inQuote && braceDepth == 0 {
				if start < i {
					attr, err := parseCallAttr(attrs[start:i])
					if err != nil {
						return nil, err
					}
					out = append(out, attr)
				}
				start = i + 1
			}
		}
	}
	if start < len(attrs) {
		attr, err := parseCallAttr(attrs[start:])
		if err != nil {
			return nil, err
		}
		out = append(out, attr)
	}
	return out, nil
}

func parseCallAttr(part string) (callAttr, error) {
	part = strings.TrimSpace(part)
	eq := strings.IndexByte(part, '=')
	if eq < 0 {
		return callAttr{}, fmt.Errorf("attribute %q must use name=value syntax", part)
	}
	name := strings.TrimSpace(part[:eq])
	value := strings.TrimSpace(part[eq+1:])
	return callAttr{Name: name, Value: value}, nil
}

func attrExpressionValue(val string) string {
	if val == "" {
		return `""`
	}
	if val[0] == '{' && val[len(val)-1] == '}' && strings.Count(val, "{") == 1 && strings.Count(val, "}") == 1 {
		return val[1 : len(val)-1]
	}
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		inner := val[1 : len(val)-1]
		if strings.Contains(inner, "{") {
			return concatPlaceholders(inner)
		}
		return fmt.Sprintf("%q", inner)
	}
	return fmt.Sprintf("%q", val)
}

func buildComponentArgs(comp *ComponentDef, attrs string, src string, pos int) (string, error) {
	if comp == nil {
		return "", fmt.Errorf("%s: unknown component", sourceRef(src, pos))
	}

	provided, err := parseCallAttrs(attrs)
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", sourceRef(src, pos), comp.Name, err)
	}

	seen := map[string]bool{}
	providedMap := map[string]string{}
	var dups []string
	var unknowns []string
	var wrongType *propTypeError
	for _, a := range provided {
		if seen[a.Name] {
			dups = append(dups, a.Name)
			continue
		}
		seen[a.Name] = true
		if !propExists(comp, a.Name) {
			unknowns = append(unknowns, a.Name)
			continue
		}
		prop := propByName(comp, a.Name)
		expr := attrExpressionValue(a.Value)
		if mismatch := checkPropLiteralType(prop, expr, a.Value); mismatch != nil {
			wrongType = &propTypeError{prop: a.Name, got: mismatch.got, want: mismatch.want}
			break
		}
		providedMap[a.Name] = expr
	}

	var missing []string
	for _, p := range comp.Props {
		if p.Default == "" && !seen[p.Name] {
			missing = append(missing, p.Name)
		}
	}

	loc := sourceRef(src, pos)
	if len(dups) > 0 {
		return "", fmt.Errorf("%s: %s %s: duplicate prop \"%s\"", loc, comp.Name, dups[0], dups[0])
	}
	if len(unknowns) > 0 {
		return "", fmt.Errorf("%s: %s %s: unknown prop \"%s\"", loc, comp.Name, unknowns[0], unknowns[0])
	}
	if wrongType != nil {
		return "", fmt.Errorf("%s: %s %s: expected %s, got %s", loc, comp.Name, wrongType.prop, wrongType.want, wrongType.got)
	}
	if len(missing) > 0 {
		return "", fmt.Errorf("%s: %s %s: missing required prop \"%s\"", loc, comp.Name, missing[0], missing[0])
	}

	var args []string
	for _, p := range comp.Props {
		if v, ok := providedMap[p.Name]; ok {
			args = append(args, v)
		} else if p.Default != "" {
			if p.Type != "string" {
				return "", fmt.Errorf("%s: %s %s: default value for %s prop is not supported", loc, comp.Name, p.Name, p.Type)
			}
			args = append(args, p.Default)
		} else {
			args = append(args, zeroValue(p.Type))
		}
	}
	return strings.Join(args, ", "), nil
}

func propByName(comp *ComponentDef, name string) Prop {
	for _, p := range comp.Props {
		if p.Name == name {
			return p
		}
	}
	return Prop{}
}

type propTypeError struct {
	prop string
	want string
	got  string
}

func checkPropLiteralType(prop Prop, expr, rawValue string) *propTypeError {
	if !isExpressionValue(rawValue) {
		return nil
	}
	if prop.Type == "" {
		return nil
	}
	kind := classifyExpression(expr)
	if kind == exprKindOther {
		return nil
	}
	want := propTypeKind(prop.Type)
	if kind != want {
		return &propTypeError{prop: prop.Name, want: prop.Type, got: kindName(kind)}
	}
	return nil
}

func isExpressionValue(raw string) bool {
	raw = strings.TrimSpace(raw)
	return len(raw) >= 2 && raw[0] == '{' && raw[len(raw)-1] == '}' && strings.Count(raw, "{") == 1 && strings.Count(raw, "}") == 1
}

func propTypeKind(typ string) exprKind {
	switch typ {
	case "string":
		return exprKindStringLiteral
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		return exprKindIntLiteral
	default:
		return exprKindOther
	}
}

func propExists(comp *ComponentDef, name string) bool {
	for _, p := range comp.Props {
		if p.Name == name {
			return true
		}
	}
	return false
}

func zeroValue(typ string) string {
	switch typ {
	case "int", "int64", "int32":
		return "0"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return ""
	}
}

func sourceRef(src string, pos int) string {
	loc := sourceLocation(src, pos)
	if src == "" {
		return "?:" + loc
	}
	return fmt.Sprintf("%s:%s", src, loc)
}
