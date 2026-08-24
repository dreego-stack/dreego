package transpiler

import (
	"fmt"
	"strings"
)

type compGen struct {
	gen       *generator
	inSection bool
	builder   string
}

func genTemplateNodeComp(gen *generator, n TemplateNode) (string, error) {
	g := &compGen{gen: gen, builder: "b"}
	return g.node(n)
}

func (g *compGen) node(n TemplateNode) (string, error) {
	switch n.Type {
	case NodeText:
		code, next := compTextSection(n.Content, g.inSection)
		g.inSection = next
		return fmt.Sprintf("%s.WriteString(%s)", g.builder, code), nil
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
			return fmt.Sprintf("%s.WriteString(%s)", g.builder, code), nil
		}
		return fmt.Sprintf("%s.WriteString(dreego.SafeText(%s))", g.builder, code), nil
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
			return fmt.Sprintf("%s.WriteString(ctx.Get(\"slot_%s\"))", g.builder, n.Content), nil
		}
		return fmt.Sprintf("%s.WriteString(ctx.Get(\"slot\"))", g.builder), nil
	case NodeComponentCall:
		return g.genComponentCall(n)
	case NodeVerbatim:
		return fmt.Sprintf("%s.WriteString(%s)", g.builder, goLiteral(n.Content)), nil
	}
	return "", fmt.Errorf("unsupported component node type %d", n.Type)
}

func (g *compGen) genComponentCall(n TemplateNode) (string, error) {
	parts := strings.SplitN(n.Tag, ".", 2)
	funcName := parts[len(parts)-1]
	def := g.gen.lookupDef(funcName)
	if def == nil {
		return "", fmt.Errorf("%s: unknown component %s", sourceRef(n.Source, n.Pos), funcName)
	}
	args, err := buildComponentArgs(def, n.Attrs, n.Source, n.Pos)
	if err != nil {
		return "", err
	}
	callName := g.gen.qualify(funcName)
	if n.SelfClose {
		return fmt.Sprintf("{ html, err := %s(%s).Render(ctx); if err != nil { return \"\", err }; %s.WriteString(html) }", callName, args, g.builder), nil
	}

	id := n.Pos
	slotBuilder := fmt.Sprintf("slotBuilder%d", id)
	previousSlot := fmt.Sprintf("previousSlot%d", id)
	var buf strings.Builder
	buf.WriteString("{\n")
	buf.WriteString(fmt.Sprintf("\t%s := ctx.Data(\"slot\")\n", previousSlot))
	buf.WriteString(fmt.Sprintf("\tvar %s strings.Builder\n", slotBuilder))
	var slotKeys []string
	var previousNamedSlots []string
	defaultGen := &compGen{gen: g.gen, builder: slotBuilder}
	for _, child := range n.Children {
		if child.Type == NodeSlot && child.Content != "" {
			if nested := findNestedSlot(child.Children); nested != nil {
				return "", nestedSlotError(n, def, nested, g.gen.src)
			}
			if err := validateSlotName(def, child.Content, n.Source, g.gen.src, n.Pos); err != nil {
				return "", err
			}
			key := "slot_" + child.Content
			slotKeys = append(slotKeys, key)
			previous := fmt.Sprintf("previousNamedSlot%d_%d", id, len(slotKeys))
			previousNamedSlots = append(previousNamedSlots, previous)
			buf.WriteString(fmt.Sprintf("\t%s := ctx.Data(%q)\n", previous, key))
			namedBuilder := fmt.Sprintf("namedSlotBuilder%d_%d", id, len(slotKeys))
			buf.WriteString(fmt.Sprintf("\tvar %s strings.Builder\n", namedBuilder))
			childGen := &compGen{gen: g.gen, builder: namedBuilder}
			for _, slotChild := range child.Children {
				code, err := childGen.node(slotChild)
				if err != nil {
					return "", err
				}
				buf.WriteString("\t" + code + "\n")
			}
			buf.WriteString(fmt.Sprintf("\tctx.Set(%q, %s.String())\n", key, namedBuilder))
			continue
		}
		if nested := findNestedSlot([]TemplateNode{child}); nested != nil {
			return "", nestedSlotError(n, def, nested, g.gen.src)
		}
		code, err := defaultGen.node(child)
		if err != nil {
			return "", err
		}
		buf.WriteString("\t" + code + "\n")
	}
	buf.WriteString(fmt.Sprintf("\tctx.Set(\"slot\", %s.String())\n", slotBuilder))
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(%s).Render(ctx)\n", callName, args))
	buf.WriteString(restoreComponentContextValue("slot", previousSlot))
	for i, key := range slotKeys {
		buf.WriteString(restoreComponentContextValue(key, previousNamedSlots[i]))
	}
	buf.WriteString("\tif err != nil { return \"\", err }\n")
	buf.WriteString(fmt.Sprintf("\t%s.WriteString(html)\n", g.builder))
	buf.WriteString("}")
	return buf.String(), nil
}

func restoreComponentContextValue(key, previous string) string {
	return fmt.Sprintf("\tif %s == nil { ctx.Delete(%q) } else { ctx.Set(%q, %s) }\n", previous, key, key, previous)
}

func GenerateComponent(gen *generator, file *File, scopeHash string) (string, error) {
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

	for _, g := range file.Server {
		if g.Code != "" {
			for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
				buf.WriteString("\t\t" + strings.TrimSpace(line) + "\n")
			}
			buf.WriteString("\n")
		}
	}

	if file.Body != nil {
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		g := &compGen{gen: gen, builder: "b"}
		for _, n := range file.Body.Nodes {
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
