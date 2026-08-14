package core

import (
	"fmt"
	"strings"
)

type compGen struct {
	inSection bool
}

func genTemplateNodeComp(n TemplateNode) (string, error) {
	g := &compGen{}
	return g.node(n)
}

func (g *compGen) node(n TemplateNode) (string, error) {
	switch n.Type {
	case NodeText:
		code, next := compTextSection(n.Content, g.inSection)
		g.inSection = next
		return fmt.Sprintf("b.WriteString(%s)", code), nil
	case NodeExpression:
		code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", n.Content)
		raw := false
		for _, f := range n.Filters {
			switch f {
			case "raw":
				raw = true
			case "upper":
				code = fmt.Sprintf("strings.ToUpper(%s)", code)
			default:
				return "", fmt.Errorf("unknown filter '%s' at position %d", f, n.Pos)
			}
		}
		if raw {
			return fmt.Sprintf("b.WriteString(%s)", code), nil
		}
		return fmt.Sprintf("b.WriteString(html.EscapeString(%s))", code), nil
	case NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("if %s {\n", n.Cond))
		for _, child := range n.Children {
			code, err := g.node(child)
			if err != nil {
				return "", err
			}
			buf.WriteString("\t\t" + code + "\n")
		}
		chain := true
		for _, ec := range n.ElseChildren {
			if ec.Type != NodeIf {
				chain = false
				break
			}
		}
		if chain {
			for _, ec := range n.ElseChildren {
				buf.WriteString(fmt.Sprintf("\t} else if %s {\n", ec.Cond))
				for _, child := range ec.Children {
					code, err := g.node(child)
					if err != nil {
						return "", err
					}
					buf.WriteString("\t\t" + code + "\n")
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString("\t} else {\n")
					for _, child := range ec.ElseChildren {
						code, err := g.node(child)
						if err != nil {
							return "", err
						}
						buf.WriteString("\t\t" + code + "\n")
					}
				}
			}
		} else {
			buf.WriteString("\t} else {\n")
			for _, ec := range n.ElseChildren {
				code, err := g.node(ec)
				if err != nil {
					return "", err
				}
				buf.WriteString("\t\t" + code + "\n")
			}
		}
		buf.WriteString("\t}")
		return buf.String(), nil
	case NodeEach:
		var buf strings.Builder
		hasElse := len(n.ElseChildren) > 0
		if hasElse {
			buf.WriteString(fmt.Sprintf("if len(%s) > 0 {\n", n.Items))
		}
		buf.WriteString(fmt.Sprintf("\tfor i, %s := range %s {\n", n.Item, n.Items))
		buf.WriteString(fmt.Sprintf("\t\tloop := dreego.EachLoop{Index: i, First: i == 0, Last: i == len(%s)-1, Even: i%%2 == 0, Odd: i%%2 != 0}\n", n.Items))
		buf.WriteString("\t\t_ = loop\n")
		for _, child := range n.Children {
			code, err := g.node(child)
			if err != nil {
				return "", err
			}
			code = strings.ReplaceAll(code, "$loop.", "loop.")
			buf.WriteString("\t\t" + code + "\n")
		}
		buf.WriteString("\t}\n")
		if hasElse {
			buf.WriteString("} else {\n")
			for _, child := range n.ElseChildren {
				code, err := g.node(child)
				if err != nil {
					return "", err
				}
				buf.WriteString("\t\t" + code + "\n")
			}
			buf.WriteString("\t}")
		}
		return buf.String(), nil
	case NodeSlot:
		if n.Content != "" {
			return fmt.Sprintf("b.WriteString(ctx.Get(\"slot_%s\"))", n.Content), nil
		}
		return "b.WriteString(ctx.Get(\"slot\"))", nil
	case NodeComponentCall:
		return genComponentCall(n)
	case NodeVerbatim:
		return fmt.Sprintf("b.WriteString(%s)", goLiteral(n.Content)), nil
	}
	return "", fmt.Errorf("unsupported component node type %d", n.Type)
}

func genComponentCall(n TemplateNode) (string, error) {
	parts := strings.SplitN(n.Tag, ".", 2)
	funcName := parts[len(parts)-1]
	if n.SelfClose {
		return fmt.Sprintf("%s(%s).Render(ctx)", funcName, extractAttrValues(n.Attrs)), nil
	}
	return fmt.Sprintf("b.WriteString(\"<@%s>\")", n.Tag), nil
}

func compTextWithAttrs(s string) string {
	code, _ := compTextSection(s, false)
	return code
}

// compTextSection renders a NodeText content segment to Go code, resolving {{ … }}
// placeholders inside quoted attribute values but leaving <script>/<style> section
// bodies literal (the lexer treats those as raw text where {{ … }} is not an expression).
// It returns the generated code and the section state after this segment so the
// caller can carry section tracking across sibling text nodes.
func compTextSection(content string, inSection bool) (string, bool) {
	var parts []string
	cur := inSection
	var quote byte
	start := 0
	i := 0
	for i < len(content) {
		if cur {
			if closeLen := sectionCloseLen(content[i:]); closeLen > 0 {
				end := i + closeLen
				parts = append(parts, goLiteral(content[start:end]))
				start = end
				i = end
				cur = false
				continue
			}
			i++
			continue
		}
		if strings.HasPrefix(content[i:], "<script") || strings.HasPrefix(content[i:], "<style") {
			if start < i {
				parts = append(parts, goLiteral(content[start:i]))
			}
			start = i
			cur = true
			quote = 0
			continue
		}
		if (content[i] == '"' || content[i] == '\'') && (i == 0 || content[i-1] != '\\') {
			if quote == 0 {
				quote = content[i]
			} else if quote == content[i] {
				quote = 0
			}
			i++
			continue
		}
		if quote != 0 && strings.HasPrefix(content[i:], "{{") {
			closeIdx := strings.Index(content[i+2:], "}}")
			if closeIdx < 0 {
				i++
				continue
			}
			if start < i {
				parts = append(parts, goLiteral(content[start:i]))
			}
			expr := strings.TrimSpace(content[i+2 : i+2+closeIdx])
			code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
			parts = append(parts, fmt.Sprintf("html.EscapeString(%s)", code))
			i += 2 + closeIdx + 2
			start = i
			continue
		}
		i++
	}
	if start < len(content) {
		parts = append(parts, goLiteral(content[start:]))
	}
	return strings.Join(parts, " + "), cur
}

func sectionCloseLen(s string) int {
	for _, tag := range []string{"</script>", "</style>"} {
		if strings.HasPrefix(s, tag) {
			return len(tag)
		}
	}
	return 0
}

func GenerateComponent(file *File, scopeHash string) (string, error) {
	comp := file.Component
	if comp == nil {
		return "", fmt.Errorf("no component definition")
	}

	var buf strings.Builder

	declParams, implParams, callArgs, variadicName := componentParams(comp)

	if variadicName != "" {
		buf.WriteString(fmt.Sprintf("func %s(%s) dreego.Component {\n", comp.Name, declParams))
		buf.WriteString("\t" + variadicName + "0 := \"\"\n")
		buf.WriteString("\tif len(" + variadicName + ") > 0 {\n\t\t" + variadicName + "0 = " + variadicName + "[0]\n\t}\n")
		buf.WriteString("\treturn component" + comp.Name + "(" + callArgs + ")\n")
		buf.WriteString("}\n\n")
		buf.WriteString(fmt.Sprintf("func component%s(%s) dreego.Component {\n", comp.Name, implParams))
	} else {
		buf.WriteString(fmt.Sprintf("func %s(%s) dreego.Component {\n", comp.Name, declParams))
	}
	buf.WriteString("\treturn dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {\n")
	writePropDefaultFallbacks(&buf, comp)
	buf.WriteString("\t\tvar b strings.Builder\n\n")

	for _, g := range file.Go {
		if g.Code != "" {
			for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
				buf.WriteString("\t\t" + strings.TrimSpace(line) + "\n")
			}
			buf.WriteString("\n")
		}
	}

	if file.Template != nil {
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		g := &compGen{}
		for _, n := range file.Template.Nodes {
			code, err := g.node(n)
			if err != nil {
				return "", err
			}
			buf.WriteString("\t\t" + code + "\n")
		}
		buf.WriteString("\t\tb.WriteString(\"</div>\")\n")
	}

	if file.Style != nil {
		scoped := scopeCSS(file.Style.Code, scopeHash)
		buf.WriteString("\t\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(%s)\n", goLiteral(scoped)))
		buf.WriteString("\t\tb.WriteString(\"</style>\")\n")
	}

	buf.WriteString("\n\t\treturn b.String(), nil\n")
	buf.WriteString("\t})\n")
	buf.WriteString("}\n\n")

	return buf.String(), nil
}
