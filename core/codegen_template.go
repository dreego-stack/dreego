package core

import (
	"fmt"
	"strings"
)

func genTemplateNode(gen *generator, n TemplateNode, depth int) (string, error) {
	return genTemplateNodeTo(gen, n, depth, "b")
}

func genTemplateNodeTo(gen *generator, n TemplateNode, depth int, builder string) (string, error) {
	indent := strings.Repeat("\t", depth)
	switch n.Type {
	case NodeText:
		if n.Content == "" {
			return "", nil
		}
		code, _ := compTextSection(n.Content, false)
		return fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, code), nil
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
			return fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, code), nil
		}
		return fmt.Sprintf(`%s%s.WriteString(dreego.SafeText(%s))`+"\n", indent, builder, code), nil
	case NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sif %s {\n", indent, n.Cond))
		for _, child := range n.Children {
			code, err := genTemplateNodeTo(gen, child, depth+1, builder)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
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
				buf.WriteString(fmt.Sprintf("%s} else if %s {\n", indent, ec.Cond))
				for _, child := range ec.Children {
					code, err := genTemplateNodeTo(gen, child, depth+1, builder)
					if err != nil {
						return "", err
					}
					buf.WriteString(code)
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
					for _, child := range ec.ElseChildren {
						code, err := genTemplateNodeTo(gen, child, depth+1, builder)
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
				code, err := genTemplateNodeTo(gen, ec, depth+1, builder)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
		buf.WriteString(indent + "}\n")
		return buf.String(), nil
	case NodeEach:
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
			code, err := genTemplateNodeTo(gen, child, forDepth+1, builder)
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
				code, err := genTemplateNodeTo(gen, child, depth+1, builder)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
			buf.WriteString(indent + "}\n")
		}
		return buf.String(), nil
	case NodeSlot:
		if n.Content != "" && len(n.Children) > 0 {
			slotBuilder := fmt.Sprintf("slotBuilder%d", depth)
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("%s{\n", indent))
			buf.WriteString(fmt.Sprintf("%s\tvar %s strings.Builder\n", indent, slotBuilder))
			for _, child := range n.Children {
				code, err := genTemplateNodeTo(gen, child, depth+2, slotBuilder)
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
	case NodeVerbatim:
		return fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, goLiteral(n.Content)), nil
	case NodeComponentCall:
		funcName := n.Tag
		if idx := strings.LastIndexByte(n.Tag, '.'); idx >= 0 {
			funcName = n.Tag[idx+1:]
		}
		def := gen.lookupDef(funcName)
		args, err := buildComponentArgs(def, n.Attrs, n.Source, n.Pos)
		if err != nil {
			return "", err
		}
		if n.SelfClose {
			return fmt.Sprintf("%s%s.WriteString(func() string { h, _ := %s(%s).Render(c); return h }())\n", indent, builder, funcName, args), nil
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
			if child.Type == NodeSlot && child.Content != "" {
				if nested := findNestedSlot(child.Children); nested != nil {
					return "", nestedSlotError(n, def, nested, gen.src)
				}
				if err := validateSlotName(def, child.Content, n.Source, gen.src, n.Pos); err != nil {
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
					code, err := genTemplateNodeTo(gen, sc, depth+3, namedSlotBuilder)
					if err != nil {
						return "", err
					}
					buf.WriteString(code)
				}
				buf.WriteString(fmt.Sprintf("%s\t\tc.Set(\"slot_%s\", %s.String())\n", indent, child.Content, namedSlotBuilder))
				buf.WriteString(fmt.Sprintf("%s\t}\n", indent))
			} else {
				if nested := findNestedSlot([]TemplateNode{child}); nested != nil {
					return "", nestedSlotError(n, def, nested, gen.src)
				}
				code, err := genTemplateNodeTo(gen, child, depth+2, slotBuilder)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
		buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot\", %s.String())\n", indent, slotBuilder))
		buf.WriteString(fmt.Sprintf("%s\thtml, err := %s(%s).Render(c)\n", indent, funcName, args))
		buf.WriteString(restoreContextValue(indent, "slot", previousSlot))
		for i, key := range slotKeys {
			buf.WriteString(restoreContextValue(indent, key, previousNamedSlots[i]))
		}
		buf.WriteString(fmt.Sprintf("%s\tif err != nil { return \"\", err }\n", indent))
		buf.WriteString(fmt.Sprintf("%s\t%s.WriteString(html)\n", indent, builder))
		buf.WriteString(fmt.Sprintf("%s}\n", indent))
		return buf.String(), nil
	}
	return "", fmt.Errorf("unsupported template node type %d", n.Type)
}

func restoreContextValue(indent, key, previous string) string {
	return fmt.Sprintf("%s\tif %s == nil { c.Delete(%q) } else { c.Set(%q, %s) }\n", indent, previous, key, key, previous)
}
