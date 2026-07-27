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
		buf.WriteString(fmt.Sprintf("%shtml, err := %s(%s).Render(c)\n", indent, funcName, args))
		buf.WriteString(indent + "if err != nil { return err }\n")
		buf.WriteString(fmt.Sprintf("%sb.WriteString(html)\n", indent))
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
	for _, part := range strings.Split(attrs, " ") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		eq := strings.IndexByte(part, '=')
		if eq < 0 {
			vals = append(vals, part)
			continue
		}
		val := strings.Trim(part[eq+1:], "\"")
		vals = append(vals, fmt.Sprintf("%q", val))
	}
	return strings.Join(vals, ", ")
}
