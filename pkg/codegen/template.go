package codegen

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func genTemplateNode(n ast.TemplateNode, depth int) string {
	indent := strings.Repeat("\t", depth)
	switch n.Type {
	case ast.NodeText:
		if n.Content == "" {
			return ""
		}
		return fmt.Sprintf("%sb.WriteString(%s)\n", indent, goLiteral(n.Content))
	case ast.NodeExpression:
		return fmt.Sprintf(`%sb.WriteString(html.EscapeString(fmt.Sprintf("%%v", %s)))`+"\n", indent, n.Content)
	case ast.NodeIf:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sif %s {\n", indent, n.Cond))
		for _, child := range n.Children {
			buf.WriteString(genTemplateNode(child, depth+1))
		}
		buf.WriteString(indent + "}\n")
		return buf.String()
	case ast.NodeEach:
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%sfor _, %s := range %s {\n", indent, n.Item, n.Items))
		for _, child := range n.Children {
			buf.WriteString(genTemplateNode(child, depth+1))
		}
		buf.WriteString(indent + "}\n")
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
