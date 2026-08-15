package core

import (
	"fmt"
	"strings"
)

func genTemplateNode(gen *generator, n TemplateNode, depth int) (string, error) {
	indent := strings.Repeat("\t", depth)
	switch n.Type {
	case NodeText:
		if n.Content == "" {
			return "", nil
		}
		code, _ := compTextSection(n.Content, false)
		return fmt.Sprintf("%sb.WriteString(%s)\n", indent, code), nil
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
			return fmt.Sprintf("%sb.WriteString(%s)\n", indent, code), nil
		}
		return fmt.Sprintf(`%sb.WriteString(html.EscapeString(%s))`+"\n", indent, code), nil
	case NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sif %s {\n", indent, n.Cond))
		for _, child := range n.Children {
			code, err := genTemplateNode(gen, child, depth+1)
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
					code, err := genTemplateNode(gen, child, depth+1)
					if err != nil {
						return "", err
					}
					buf.WriteString(code)
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
					for _, child := range ec.ElseChildren {
						code, err := genTemplateNode(gen, child, depth+1)
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
				code, err := genTemplateNode(gen, ec, depth+1)
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
			code, err := genTemplateNode(gen, child, forDepth+1)
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
				code, err := genTemplateNode(gen, child, depth+1)
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
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("%s{\n", indent))
			buf.WriteString(fmt.Sprintf("%s\tvar cb strings.Builder\n", indent))
			for _, child := range n.Children {
				code, err := genTemplateNode(gen, child, depth+2)
				if err != nil {
					return "", err
				}
				code = strings.ReplaceAll(code, "b.WriteString(", "cb.WriteString(")
				buf.WriteString(code)
			}
			buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot_%s\", cb.String())\n", indent, n.Content))
			buf.WriteString(fmt.Sprintf("%s}\n", indent))
			return buf.String(), nil
		}
		if n.Content != "" {
			return fmt.Sprintf("%sb.WriteString(c.Get(\"slot_%s\"))\n", indent, n.Content), nil
		}
		return fmt.Sprintf("%sb.WriteString(c.Get(\"slot\"))\n", indent), nil
	case NodeVerbatim:
		return fmt.Sprintf("%sb.WriteString(%s)\n", indent, goLiteral(n.Content)), nil
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
			return fmt.Sprintf("%sb.WriteString(func() string { h, _ := %s(%s).Render(c); return h }())\n", indent, funcName, args), nil
		}
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%s{\n", indent))
		buf.WriteString(fmt.Sprintf("%s\tvar cb strings.Builder\n", indent))
		for _, child := range n.Children {
			if child.Type == NodeSlot && child.Content != "" && len(child.Children) > 0 {
				buf.WriteString(fmt.Sprintf("%s\t{\n", indent))
				buf.WriteString(fmt.Sprintf("%s\t\tvar sb strings.Builder\n", indent))
				for _, sc := range child.Children {
					code, err := genTemplateNode(gen, sc, depth+3)
					if err != nil {
						return "", err
					}
					code = strings.ReplaceAll(code, "b.WriteString(", "sb.WriteString(")
					buf.WriteString(code)
				}
				buf.WriteString(fmt.Sprintf("%s\t\tc.Set(\"slot_%s\", sb.String())\n", indent, child.Content))
				buf.WriteString(fmt.Sprintf("%s\t}\n", indent))
			} else {
				code, err := genTemplateNode(gen, child, depth+2)
				if err != nil {
					return "", err
				}
				code = strings.ReplaceAll(code, "b.WriteString(", "cb.WriteString(")
				buf.WriteString(code)
			}
		}
		buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot\", cb.String())\n", indent))
		buf.WriteString(fmt.Sprintf("%s\thtml, err := %s(%s).Render(c)\n", indent, funcName, args))
		buf.WriteString(fmt.Sprintf("%s\tif err != nil { return \"\", err }\n", indent))
		buf.WriteString(fmt.Sprintf("%s\tb.WriteString(html)\n", indent))
		buf.WriteString(fmt.Sprintf("%s}\n", indent))
		return buf.String(), nil
	}
	return "", fmt.Errorf("unsupported template node type %d", n.Type)
}
