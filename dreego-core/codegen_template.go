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
		return fmt.Sprintf(`%sb.WriteString(html.EscapeString(fmt.Sprintf("%%v", %s)))`+"\n", indent, n.Content)
	case NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sif %s {\n", indent, n.Cond))
		for _, child := range n.Children {
			buf.WriteString(genTemplateNode(child, depth+1))
		}
		buf.WriteString(indent + "}\n")
		return buf.String()
	case NodeEach:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sfor _, %s := range %s {\n", indent, n.Item, n.Items))
		for _, child := range n.Children {
			buf.WriteString(genTemplateNode(child, depth+1))
		}
		buf.WriteString(indent + "}\n")
		return buf.String()
	case NodeSlot:
		return fmt.Sprintf("%sb.WriteString(c.Get(\"slot\"))\n", indent)
	case NodeComponentCall:
		funcName := n.Tag
		if idx := strings.LastIndexByte(n.Tag, '.'); idx >= 0 {
			funcName = n.Tag[idx+1:]
		}
		args := extractAttrValues(n.Attrs)
		if n.SelfClose {
			return fmt.Sprintf("%sh, _ := %s(%s).Render(c); %sb.WriteString(h)\n", indent, funcName, args, indent)
		}
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%s{\n", indent))
		buf.WriteString(fmt.Sprintf("%s\tvar cb strings.Builder\n", indent))
		for _, child := range n.Children {
			childCode := genTemplateNode(child, depth+2)
			childCode = strings.ReplaceAll(childCode, "b.WriteString(", "cb.WriteString(")
			buf.WriteString(childCode)
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
	val := strings.Trim(part[eq+1:], "\"")
	return fmt.Sprintf("%q", val)
}
