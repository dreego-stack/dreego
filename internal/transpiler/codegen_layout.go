package transpiler

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func GenerateLayout(gen *Generator, file *File, funcName string) (string, error) {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func %s(c dreego.RenderContext, content, head string) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if file.Head != nil {
		buf.WriteString("\tb.WriteString(head)\n")
	}

	if file.Body != nil {
		inSection := false
		for _, n := range file.Body.Nodes {
			code, err := genLayoutNodeState(gen, n, 1, &inSection)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
		}
	}

	if file.Style != nil {
		buf.WriteString("\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", ir.GoLiteral(file.Style.Code)))
		buf.WriteString("\tb.WriteString(\"</style>\")\n")
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")
	return buf.String(), nil
}

func genLayoutNode(gen *Generator, n TemplateNode, depth int) (string, error) {
	inSection := false
	return genLayoutNodeState(gen, n, depth, &inSection)
}

func genLayoutNodeState(gen *Generator, n TemplateNode, depth int, inSection *bool) (string, error) {
	indent := strings.Repeat("\t", depth)
	if n.Type == ir.NodeSlot {
		if n.Content != "" {
			return indent + fmt.Sprintf("b.WriteString(c.Get(\"slot_%s\"))\n", n.Content), nil
		}
		return indent + "b.WriteString(content)\n", nil
	}
	if n.Type == ir.NodeText && (strings.Contains(n.Content, "{#head}") || strings.Contains(n.Content, "{#slot}")) {
		parts := splitLayoutText(n.Content)
		var out string
		for _, p := range parts {
			switch p {
			case "{#head}":
				out += indent + "b.WriteString(head)\n"
			case "{#slot}":
				out += indent + "b.WriteString(content)\n"
			default:
				out += indent + fmt.Sprintf("b.WriteString(%s)\n", ir.GoLiteral(p))
			}
		}
		return out, nil
	}
	return html.GenTemplateNodeToState(gen, n, depth, "b", inSection)
}

func splitLayoutText(s string) []string {
	var result []string
	for s != "" {
		headIdx := strings.Index(s, "{#head}")
		slotIdx := strings.Index(s, "{#slot}")

		next := -1
		nextLen := 0
		if headIdx >= 0 && (slotIdx < 0 || headIdx <= slotIdx) {
			next = headIdx
			nextLen = 7
		} else if slotIdx >= 0 {
			next = slotIdx
			nextLen = 7
		}

		if next < 0 {
			if s != "" {
				result = append(result, s)
			}
			break
		}
		if next > 0 {
			result = append(result, s[:next])
		}
		result = append(result, s[next:next+nextLen])
		s = s[next+nextLen:]
	}
	return result
}
