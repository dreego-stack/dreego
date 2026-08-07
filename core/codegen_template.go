package core

import (
	"fmt"
	"strings"
)

func genTemplateNode(n TemplateNode, depth int) (string, error) {
	indent := strings.Repeat("\t", depth)
	switch n.Type {
	case NodeText:
		if n.Content == "" {
			return "", nil
		}
		return fmt.Sprintf("%sb.WriteString(%s)\n", indent, goLiteral(n.Content)), nil
	case NodeExpression:
		code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", n.Content)
		raw := false
		for _, f := range n.Filters {
			switch f {
			case "raw":
				raw = true
			case "upper":
				code = fmt.Sprintf("strings.ToUpper(%s)", code)
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
			code, err := genTemplateNode(child, depth+1)
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
					code, err := genTemplateNode(child, depth+1)
					if err != nil {
						return "", err
					}
					buf.WriteString(code)
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
					for _, child := range ec.ElseChildren {
						code, err := genTemplateNode(child, depth+1)
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
				code, err := genTemplateNode(ec, depth+1)
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
			code, err := genTemplateNode(child, forDepth+1)
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
				code, err := genTemplateNode(child, depth+1)
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
				code, err := genTemplateNode(child, depth+2)
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
		args := extractAttrValues(n.Attrs)
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
					code, err := genTemplateNode(sc, depth+3)
					if err != nil {
						return "", err
					}
					code = strings.ReplaceAll(code, "b.WriteString(", "sb.WriteString(")
					buf.WriteString(code)
				}
				buf.WriteString(fmt.Sprintf("%s\t\tc.Set(\"slot_%s\", sb.String())\n", indent, child.Content))
				buf.WriteString(fmt.Sprintf("%s\t}\n", indent))
			} else {
				code, err := genTemplateNode(child, depth+2)
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

func genTempl(file *File, layout *File, scopeHash string, isGET bool) (string, error) {
	var buf strings.Builder

	if layout == nil && file.Head != nil && isGET {
		headCode, err := genHead(file.Head.Content, "b")
		if err != nil {
			return "", err
		}
		buf.WriteString(headCode)
	}
	if isGET {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
	}
	for _, n := range file.Template.Nodes {
		code, err := genTemplateNode(n, 1)
		if err != nil {
			return "", err
		}
		buf.WriteString(code)
	}
	if isGET {
		buf.WriteString("\tb.WriteString(\"</div>\")\n")
	}

	if file.Script != nil {
		buf.WriteString("\tb.WriteString(\"<script>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Script.Code)))
		buf.WriteString("\tb.WriteString(\"</script>\")\n")
	}
	if file.Style != nil {
		scoped := scopeCSS(file.Style.Code, scopeHash)
		buf.WriteString("\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(scoped)))
		buf.WriteString("\tb.WriteString(\"</style>\")\n")
	}

	if layout != nil && isGET {
		buf.WriteString("\tpageContent := b.String()\n")
		buf.WriteString("\tb.Reset()\n")

		if layout.Head != nil && layout.Head.Content != "" {
			headContent := layout.Head.Content
			if strings.Contains(headContent, "{#head}") {
				parts := strings.SplitN(headContent, "{#head}", 2)
				if parts[0] != "" {
					buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(parts[0])))
				}
				if file.Head != nil {
					headCode, err := genHead(file.Head.Content, "b")
					if err != nil {
						return "", err
					}
					buf.WriteString(headCode)
				} else {
					buf.WriteString("\tb.WriteString(\"\")\n")
				}
				if len(parts) > 1 && parts[1] != "" {
					buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(parts[1])))
				}
			} else {
				buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(headContent)))
			}
		}

		if file.Head != nil {
			buf.WriteString("\tvar headBuf strings.Builder\n")
			headCode, err := genHead(file.Head.Content, "headBuf")
			if err != nil {
				return "", err
			}
			buf.WriteString(headCode)
			buf.WriteString("\tc.Set(\"head\", headBuf.String())\n")
		}
		buf.WriteString("\tc.Set(\"slot\", pageContent)\n")
		if layout.Template != nil {
			for _, n := range layout.Template.Nodes {
				code, err := genLayoutNode(n, 1)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
	}

	return buf.String(), nil
}
