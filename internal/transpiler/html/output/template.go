package output

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func GenTemplateNode(gen *ir.Generator, n ir.TemplateNode, depth int) (string, error) {
	inSection := false
	return GenTemplateNodeToState(gen, n, depth, "b", &inSection)
}

func GenTemplateNodeToState(gen *ir.Generator, n ir.TemplateNode, depth int, builder string, inSection *bool) (string, error) {
	indent := strings.Repeat("\t", depth)
	switch n.Type {
	case ir.NodeText:
		if n.Content == "" {
			return "", nil
		}
		code, next := CompTextSection(n.Content, *inSection)
		*inSection = next
		return fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, code), nil
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
			return fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, code), nil
		}
		return fmt.Sprintf(`%s%s.WriteString(dreego.SafeText(%s))`+"\n", indent, builder, code), nil
	case ir.NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sif %s {\n", indent, n.Cond))
		for _, child := range n.Children {
			code, err := GenTemplateNodeToState(gen, child, depth+1, builder, inSection)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
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
				buf.WriteString(fmt.Sprintf("%s} else if %s {\n", indent, ec.Cond))
				for _, child := range ec.Children {
					code, err := GenTemplateNodeToState(gen, child, depth+1, builder, inSection)
					if err != nil {
						return "", err
					}
					buf.WriteString(code)
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
					for _, child := range ec.ElseChildren {
						code, err := GenTemplateNodeToState(gen, child, depth+1, builder, inSection)
						if err != nil {
							return "", err
						}
						buf.WriteString(code)
					}
				}
			}
		} else {
			buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
			for _, ec := range n.ElseChildren {
				code, err := GenTemplateNodeToState(gen, ec, depth+1, builder, inSection)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
		buf.WriteString(indent + "}\n")
		return buf.String(), nil
	case ir.NodeEach:
		var buf strings.Builder
		hasElse := len(n.ElseChildren) > 0
		forDepth := depth + 1
		forIndent := strings.Repeat("\t", forDepth)
		if hasElse {
			buf.WriteString(fmt.Sprintf("%sif len(%s) > 0 {\n", indent, n.Items))
			forDepth = depth + 2
			forIndent = strings.Repeat("\t", forDepth)
		}
		buf.WriteString(fmt.Sprintf("%sfor i, %s := range %s {\n", forIndent, n.Item, n.Items))
		buf.WriteString(fmt.Sprintf("%s\tloop := dreego.EachLoop{Index: i, First: i == 0, Last: i == len(%s)-1, Even: i%%2 == 0, Odd: i%%2 != 0}\n", forIndent, n.Items))
		buf.WriteString(fmt.Sprintf("%s\t_ = loop\n", forIndent))
		for _, child := range n.Children {
			code, err := GenTemplateNodeToState(gen, child, forDepth+1, builder, inSection)
			if err != nil {
				return "", err
			}
			code = strings.ReplaceAll(code, "$loop.", "loop.")
			buf.WriteString(code)
		}
		buf.WriteString(forIndent + "}\n")
		if hasElse {
			buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
			for _, child := range n.ElseChildren {
				code, err := GenTemplateNodeToState(gen, child, depth+1, builder, inSection)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
			buf.WriteString(indent + "}\n")
		}
		return buf.String(), nil
	case ir.NodeSlot:
		if n.Content != "" && len(n.Children) > 0 {
			slotBuilder := fmt.Sprintf("slotBuilder%d", depth)
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("%s{\n", indent))
			buf.WriteString(fmt.Sprintf("%s\tvar %s strings.Builder\n", indent, slotBuilder))
			for _, child := range n.Children {
				code, err := GenTemplateNodeToState(gen, child, depth+2, slotBuilder, inSection)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
			buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot_%s\", %s.String())\n", indent, n.Content, slotBuilder))
			buf.WriteString(fmt.Sprintf("%s}\n", indent))
			return buf.String(), nil
		}
		if n.Content != "" {
			return fmt.Sprintf("%s%s.WriteString(c.Get(\"slot_%s\"))\n", indent, builder, n.Content), nil
		}
		return fmt.Sprintf("%s%s.WriteString(c.Get(\"slot\"))\n", indent, builder), nil
	case ir.NodeVerbatim:
		return fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, ir.GoLiteral(n.Content)), nil
	case ir.NodeComponentCall:
		funcName := n.Tag
		if idx := strings.LastIndexByte(n.Tag, '.'); idx >= 0 {
			funcName = n.Tag[idx+1:]
		}
		def := gen.LookupDef(funcName)
		if def == nil {
			return "", fmt.Errorf("%s: unknown component %s", ir.SourceRef(n.Source, n.Pos), funcName)
		}
		args, err := BuildComponentArgs(def, n.Attrs, n.Source, n.Pos)
		if err != nil {
			return "", err
		}
		callName := gen.Qualify(funcName)
		if n.SelfClose {
			return fmt.Sprintf("%s{\n%s\tresult, err := %s(%s).Render(c)\n%s\tif err != nil { return \"\", err }\n%s\t%s.Write(result.HTML)\n%s}\n", indent, indent, callName, args, indent, indent, builder, indent), nil
		}
		slotBuilder := fmt.Sprintf("slotBuilder%d", depth)
		previousSlot := fmt.Sprintf("previousSlot%d", depth)
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%s{\n", indent))
		buf.WriteString(fmt.Sprintf("%s\t%s := c.Data(\"slot\")\n", indent, previousSlot))
		buf.WriteString(fmt.Sprintf("%s\tvar %s strings.Builder\n", indent, slotBuilder))
		var slotKeys []string
		var previousNamedSlots []string
		for _, child := range n.Children {
			if child.Type == ir.NodeSlot && child.Content != "" {
				if nested := FindNestedSlot(child.Children); nested != nil {
					return "", NestedSlotError(n, def, nested, gen.Src)
				}
				if err := ValidateSlotName(def, child.Content, n.Source, gen.Src, n.Pos); err != nil {
					return "", err
				}
				slotKeys = append(slotKeys, "slot_"+child.Content)
				previousNamedSlot := fmt.Sprintf("previousNamedSlot%d_%d", depth, len(slotKeys))
				previousNamedSlots = append(previousNamedSlots, previousNamedSlot)
				buf.WriteString(fmt.Sprintf("%s\t%s := c.Data(\"slot_%s\")\n", indent, previousNamedSlot, child.Content))
				if len(child.Children) == 0 {
					continue
				}
				namedSlotBuilder := fmt.Sprintf("namedSlotBuilder%d", depth)
				buf.WriteString(fmt.Sprintf("%s\t{\n", indent))
				buf.WriteString(fmt.Sprintf("%s\t\tvar %s strings.Builder\n", indent, namedSlotBuilder))
				for _, sc := range child.Children {
					code, err := GenTemplateNodeToState(gen, sc, depth+3, namedSlotBuilder, inSection)
					if err != nil {
						return "", err
					}
					buf.WriteString(code)
				}
				buf.WriteString(fmt.Sprintf("%s\t\tc.Set(\"slot_%s\", %s.String())\n", indent, child.Content, namedSlotBuilder))
				buf.WriteString(fmt.Sprintf("%s\t}\n", indent))
			} else {
				if nested := FindNestedSlot([]ir.TemplateNode{child}); nested != nil {
					return "", NestedSlotError(n, def, nested, gen.Src)
				}
				code, err := GenTemplateNodeToState(gen, child, depth+2, slotBuilder, inSection)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
		buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot\", %s.String())\n", indent, slotBuilder))
		buf.WriteString(fmt.Sprintf("%s\tresult, err := %s(%s).Render(c)\n", indent, callName, args))
		buf.WriteString(RestoreContextValue(indent, "slot", previousSlot))
		for i, key := range slotKeys {
			buf.WriteString(RestoreContextValue(indent, key, previousNamedSlots[i]))
		}
		buf.WriteString(fmt.Sprintf("%s\tif err != nil { return \"\", err }\n", indent))
		buf.WriteString(fmt.Sprintf("%s\t%s.Write(result.HTML)\n", indent, builder))
		buf.WriteString(fmt.Sprintf("%s}\n", indent))
		return buf.String(), nil
	}
	return "", fmt.Errorf("unsupported template node type %d", n.Type)
}

func RestoreContextValue(indent, key, previous string) string {
	return fmt.Sprintf("%s\tif %s == nil { c.Delete(%q) } else { c.Set(%q, %s) }\n", indent, previous, key, key, previous)
}
