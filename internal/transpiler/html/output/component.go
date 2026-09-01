package output

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

type CompGen struct {
	Gen       *ir.Generator
	InSection bool
	Builder   string
}

func GenTemplateNodeComp(gen *ir.Generator, n ir.TemplateNode) (string, error) {
	g := &CompGen{Gen: gen, Builder: "b"}
	return g.Node(n)
}

func GenComponentCall(gen *ir.Generator, builder string, n ir.TemplateNode) (string, error) {
	g := &CompGen{Gen: gen, Builder: builder}
	return g.genComponentCall(n)
}

func (g *CompGen) Node(n ir.TemplateNode) (string, error) {
	switch n.Type {
	case ir.NodeText:
		code, next := CompTextSection(n.Content, g.InSection)
		g.InSection = next
		return fmt.Sprintf("%s.WriteString(%s)", g.Builder, code), nil
	case ir.NodeExpression:
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
			return fmt.Sprintf("%s.WriteString(%s)", g.Builder, code), nil
		}
		return fmt.Sprintf("%s.WriteString(dreego.SafeText(%s))", g.Builder, code), nil
	case ir.NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("if %s {\n", n.Cond))
		for _, child := range n.Children {
			code, err := g.Node(child)
			if err != nil {
				return "", err
			}
			buf.WriteString("\t\t" + code + "\n")
		}
		chain := true
		for _, ec := range n.ElseChildren {
			if ec.Type != ir.NodeIf {
				chain = false
				break
			}
		}
		if chain {
			for _, ec := range n.ElseChildren {
				buf.WriteString(fmt.Sprintf("\t} else if %s {\n", ec.Cond))
				for _, child := range ec.Children {
					code, err := g.Node(child)
					if err != nil {
						return "", err
					}
					buf.WriteString("\t\t" + code + "\n")
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString("\t} else {\n")
					for _, child := range ec.ElseChildren {
						code, err := g.Node(child)
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
				code, err := g.Node(ec)
				if err != nil {
					return "", err
				}
				buf.WriteString("\t\t" + code + "\n")
			}
		}
		buf.WriteString("\t}")
		return buf.String(), nil
	case ir.NodeEach:
		var buf strings.Builder
		hasElse := len(n.ElseChildren) > 0
		if hasElse {
			buf.WriteString(fmt.Sprintf("if len(%s) > 0 {\n", n.Items))
		}
		buf.WriteString(fmt.Sprintf("\tfor i, %s := range %s {\n", n.Item, n.Items))
		buf.WriteString(fmt.Sprintf("\t\tloop := dreego.EachLoop{Index: i, First: i == 0, Last: i == len(%s)-1, Even: i%%2 == 0, Odd: i%%2 != 0}\n", n.Items))
		buf.WriteString("\t\t_ = loop\n")
		for _, child := range n.Children {
			code, err := g.Node(child)
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
				code, err := g.Node(child)
				if err != nil {
					return "", err
				}
				buf.WriteString("\t\t" + code + "\n")
			}
			buf.WriteString("\t}")
		}
		return buf.String(), nil
	case ir.NodeSlot:
		if n.Content != "" {
			return fmt.Sprintf("%s.WriteString(ctx.Get(\"slot_%s\"))", g.Builder, n.Content), nil
		}
		return fmt.Sprintf("%s.WriteString(ctx.Get(\"slot\"))", g.Builder), nil
	case ir.NodeComponentCall:
		return g.genComponentCall(n)
	case ir.NodeVerbatim:
		return fmt.Sprintf("%s.WriteString(%s)", g.Builder, ir.GoLiteral(n.Content)), nil
	}
	return "", fmt.Errorf("unsupported component node type %d", n.Type)
}

func (g *CompGen) genComponentCall(n ir.TemplateNode) (string, error) {
	parts := strings.SplitN(n.Tag, ".", 2)
	funcName := parts[len(parts)-1]
	def := g.Gen.LookupDef(funcName)
	if def == nil {
		return "", fmt.Errorf("%s: unknown component %s", ir.SourceRef(n.Source, n.Pos), funcName)
	}
	args, err := BuildComponentArgs(def, n.Attrs, n.Source, n.Pos)
	if err != nil {
		return "", err
	}
	callName := g.Gen.Qualify(funcName)
	if n.SelfClose {
		return fmt.Sprintf("{ result, err := %s(%s).Render(ctx); if err != nil { return dreego.Result{}, err }; %s.Write(result.HTML) }", callName, args, g.Builder), nil
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
	defaultGen := &CompGen{Gen: g.Gen, Builder: slotBuilder}
	for _, child := range n.Children {
		if child.Type == ir.NodeSlot && child.Content != "" {
			if nested := FindNestedSlot(child.Children); nested != nil {
				return "", NestedSlotError(n, def, nested, g.Gen.Src)
			}
			if err := ValidateSlotName(def, child.Content, n.Source, g.Gen.Src, n.Pos); err != nil {
				return "", err
			}
			key := "slot_" + child.Content
			slotKeys = append(slotKeys, key)
			previous := fmt.Sprintf("previousNamedSlot%d_%d", id, len(slotKeys))
			previousNamedSlots = append(previousNamedSlots, previous)
			buf.WriteString(fmt.Sprintf("\t%s := ctx.Data(%q)\n", previous, key))
			namedBuilder := fmt.Sprintf("namedSlotBuilder%d_%d", id, len(slotKeys))
			buf.WriteString(fmt.Sprintf("\tvar %s strings.Builder\n", namedBuilder))
			childGen := &CompGen{Gen: g.Gen, Builder: namedBuilder}
			for _, slotChild := range child.Children {
				code, err := childGen.Node(slotChild)
				if err != nil {
					return "", err
				}
				buf.WriteString("\t" + code + "\n")
			}
			buf.WriteString(fmt.Sprintf("\tctx.Set(%q, %s.String())\n", key, namedBuilder))
			continue
		}
		if nested := FindNestedSlot([]ir.TemplateNode{child}); nested != nil {
			return "", NestedSlotError(n, def, nested, g.Gen.Src)
		}
		code, err := defaultGen.Node(child)
		if err != nil {
			return "", err
		}
		buf.WriteString("\t" + code + "\n")
	}
	buf.WriteString(fmt.Sprintf("\tctx.Set(\"slot\", %s.String())\n", slotBuilder))
	buf.WriteString(fmt.Sprintf("\tresult, err := %s(%s).Render(ctx)\n", callName, args))
	buf.WriteString(RestoreComponentContextValue("slot", previousSlot))
	for i, key := range slotKeys {
		buf.WriteString(RestoreComponentContextValue(key, previousNamedSlots[i]))
	}
	buf.WriteString("\tif err != nil { return dreego.Result{}, err }\n")
	buf.WriteString(fmt.Sprintf("\t%s.Write(result.HTML)\n", g.Builder))
	buf.WriteString("}")
	return buf.String(), nil
}

func RestoreComponentContextValue(key, previous string) string {
	return fmt.Sprintf("\tif %s == nil { ctx.Delete(%q) } else { ctx.Set(%q, %s) }\n", previous, key, key, previous)
}
