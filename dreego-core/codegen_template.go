package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

)

func genTemplateNode(n TemplateNode, depth int) string {
	indent := strings.Repeat("\t", depth)
	switch n.Type {
	case NodeText:
		if n.Content == "" {
			return ""
		}
		return fmt.Sprintf("%sb.WriteString(%s)\n", indent, goLiteral(n.Content))
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
			return fmt.Sprintf("%sb.WriteString(%s)\n", indent, code)
		}
		return fmt.Sprintf(`%sb.WriteString(html.EscapeString(%s))`+"\n", indent, code)
	case NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sif %s {\n", indent, n.Cond))
		for _, child := range n.Children {
			buf.WriteString(genTemplateNode(child, depth+1))
		}
		if len(n.ElseChildren) > 0 {
			buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
			for _, child := range n.ElseChildren {
				buf.WriteString(genTemplateNode(child, depth+1))
			}
		}
		buf.WriteString(indent + "}\n")
		return buf.String()
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
		buf.WriteString(fmt.Sprintf("%s\tloop := core.EachLoop{Index: i, First: i == 0, Last: i == len(%s)-1, Even: i%%2 == 0, Odd: i%%2 != 0}\n", forIndent, n.Items))
		buf.WriteString(fmt.Sprintf("%s\t_ = loop\n", forIndent))
		for _, child := range n.Children {
			code := genTemplateNode(child, forDepth+1)
			code = strings.ReplaceAll(code, "$loop.", "loop.")
			buf.WriteString(code)
		}
		buf.WriteString(forIndent + "}\n")
		if hasElse {
			buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
			for _, child := range n.ElseChildren {
				buf.WriteString(genTemplateNode(child, depth+1))
			}
			buf.WriteString(indent + "}\n")
		}
		return buf.String()
	case NodeSlot:
		if n.Content != "" && len(n.Children) > 0 {
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("%s{\n", indent))
			buf.WriteString(fmt.Sprintf("%s\tvar cb strings.Builder\n", indent))
			for _, child := range n.Children {
				code := genTemplateNode(child, depth+2)
				code = strings.ReplaceAll(code, "b.WriteString(", "cb.WriteString(")
				buf.WriteString(code)
			}
			buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot_%s\", cb.String())\n", indent, n.Content))
			buf.WriteString(fmt.Sprintf("%s}\n", indent))
			return buf.String()
		}
		if n.Content != "" {
			return fmt.Sprintf("%sb.WriteString(c.Get(\"slot_%s\"))\n", indent, n.Content)
		}
		return fmt.Sprintf("%sb.WriteString(c.Get(\"slot\"))\n", indent)
	case NodeVerbatim:
		return fmt.Sprintf("%sb.WriteString(%s)\n", indent, goLiteral(n.Content))
	case NodeComponentCall:
		funcName := n.Tag
		if idx := strings.LastIndexByte(n.Tag, '.'); idx >= 0 {
			funcName = n.Tag[idx+1:]
		}
		args := extractAttrValues(n.Attrs)
		if n.SelfClose {
			return fmt.Sprintf("%sb.WriteString(func() string { h, _ := %s(%s).Render(c); return h }())\n", indent, funcName, args)
		}
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%s{\n", indent))
		buf.WriteString(fmt.Sprintf("%s\tvar cb strings.Builder\n", indent))
		for _, child := range n.Children {
			if child.Type == NodeSlot && child.Content != "" && len(child.Children) > 0 {
				buf.WriteString(fmt.Sprintf("%s\t{\n", indent))
				buf.WriteString(fmt.Sprintf("%s\t\tvar sb strings.Builder\n", indent))
				for _, sc := range child.Children {
					code := genTemplateNode(sc, depth+3)
					code = strings.ReplaceAll(code, "b.WriteString(", "sb.WriteString(")
					buf.WriteString(code)
				}
				buf.WriteString(fmt.Sprintf("%s\t\tc.Set(\"slot_%s\", sb.String())\n", indent, child.Content))
				buf.WriteString(fmt.Sprintf("%s\t}\n", indent))
			} else {
				code := genTemplateNode(child, depth+2)
				code = strings.ReplaceAll(code, "b.WriteString(", "cb.WriteString(")
				buf.WriteString(code)
			}
		}
		buf.WriteString(fmt.Sprintf("%s\tc.Set(\"slot\", cb.String())\n", indent))
		buf.WriteString(fmt.Sprintf("%s\thtml, err := %s(%s).Render(c)\n", indent, funcName, args))
		buf.WriteString(fmt.Sprintf("%s\tif err != nil { return \"\", err }\n", indent))
		buf.WriteString(fmt.Sprintf("%s\tb.WriteString(html)\n", indent))
		buf.WriteString(fmt.Sprintf("%s}\n", indent))
		return buf.String()
	}
	return ""
}

func goLiteral(s string) string {
	if strings.Contains(s, "`") {
		return strconv.Quote(s)
	}
	return "`" + s + "`"
}

func scopeCSS(css string, hash string) string {
	prefix := fmt.Sprintf("[data-scope=%s] ", hash)
	var result strings.Builder
	for _, rule := range strings.Split(css, "}") {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		result.WriteString(prefix)
		result.WriteString(rule)
		result.WriteString("}\n")
	}
	return strings.TrimSpace(result.String())
}

func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	var result strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		result.WriteByte(byte(unicode.ToUpper(rune(p[0]))))
		if len(p) > 1 {
			result.WriteString(strings.ToLower(p[1:]))
		}
	}
	return result.String()
}

func extractAttrValues(attrs string) string {
	if attrs == "" {
		return ""
	}
	var vals []string
	inQuote := false
	start := 0
	for i := 0; i < len(attrs); i++ {
		if attrs[i] == '"' {
			inQuote = !inQuote
		}
		if attrs[i] == ' ' && !inQuote {
			if start < i {
				vals = append(vals, attrVal(attrs[start:i]))
			}
			start = i + 1
		}
	}
	if start < len(attrs) {
		vals = append(vals, attrVal(attrs[start:]))
	}
	return strings.Join(vals, ", ")
}

func attrVal(part string) string {
	eq := strings.IndexByte(part, '=')
	if eq < 0 {
		return fmt.Sprintf("%q", part)
	}
	val := strings.TrimSpace(part[eq+1:])
	if val == "" {
		return fmt.Sprintf("%q", "")
	}
	if val[0] == '{' && val[len(val)-1] == '}' {
		return val[1 : len(val)-1]
	}
	val = strings.Trim(val, "\"")
	return fmt.Sprintf("%q", val)
}
