package core

import (
	"fmt"
	"strings"
)

func genTemplateNodeComp(n TemplateNode) (string, error) {
	switch n.Type {
	case NodeText:
		return fmt.Sprintf("b.WriteString(%s)", goLiteral(n.Content)), nil
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
			return fmt.Sprintf("b.WriteString(%s)", code), nil
		}
		return fmt.Sprintf("b.WriteString(html.EscapeString(%s))", code), nil
	case NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("if %s {\n", n.Cond))
		for _, child := range n.Children {
			code, err := genTemplateNodeComp(child)
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
					code, err := genTemplateNodeComp(child)
					if err != nil {
						return "", err
					}
					buf.WriteString("\t\t" + code + "\n")
				}
				if len(ec.ElseChildren) > 0 {
					buf.WriteString("\t} else {\n")
					for _, child := range ec.ElseChildren {
						code, err := genTemplateNodeComp(child)
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
				code, err := genTemplateNodeComp(ec)
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
		buf.WriteString(fmt.Sprintf("\t\tloop := core.EachLoop{Index: i, First: i == 0, Last: i == len(%s)-1, Even: i%%2 == 0, Odd: i%%2 != 0}\n", n.Items))
		buf.WriteString("\t\t_ = loop\n")
		for _, child := range n.Children {
			code, err := genTemplateNodeComp(child)
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
				code, err := genTemplateNodeComp(child)
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
		return fmt.Sprintf("%s(%s).Render(ctx)", funcName, n.Attrs), nil
	}
	return fmt.Sprintf("b.WriteString(\"<@%s>\")", n.Tag), nil
}
