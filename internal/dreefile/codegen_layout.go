package dreefile

import (
	"fmt"
	"os"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func GenerateLayout(gen *Generator, file *File, funcName string) (string, error) {
	if warning, ok := layoutHeadDedupeWarning(file); ok {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func %s(c dreego.RenderContext, content, head string) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if file.Head != nil {
		buf.WriteString("\tb.WriteString(head)\n")
	}

	var styleCode string
	if file.Style != nil {
		styleCode = "\tb.WriteString(\"<style>\")\n"
		styleCode += fmt.Sprintf("\tb.WriteString(%s)\n", gogen.GoLiteral(file.Style.Code))
		styleCode += "\tb.WriteString(\"</style>\")\n"
	}

	if file.Body != nil {
		inSection := false
		headIdx, headAfter := -1, ""
		headEnd := -1
		headPrefix, headTail := "", ""
		if file.Head == nil {
			if idx, prefix, after, ok := file.StaticBodyHead(); ok {
				tail, end, tailOK := file.StaticBodyHeadTail(idx, after)
				captureTail := tailOK && ir.HasHeadDedupeTag(tail)
				if captureTail || ir.HasHeadDedupeTag(prefix) {
					headIdx, headAfter = idx, after
					headPrefix = prefix
					if captureTail {
						headTail, headEnd = tail, end
					}
				}
			}
		}
		if headPrefix != "" || headTail != "" {
			if headPrefix != "" {
				buf.WriteString(fmt.Sprintf("\tlayoutHead := %s\n", gogen.GoLiteral(headPrefix)))
				writeHeadDedupe(&buf, "layoutHead")
				buf.WriteString("\tb.WriteString(layoutHead)\n")
			}
			buf.WriteString("\tb.WriteString(head)\n")
			if headTail != "" {
				buf.WriteString(fmt.Sprintf("\tlayoutTail := %s\n", gogen.GoLiteral(headTail)))
				writeHeadDedupe(&buf, "layoutTail")
				buf.WriteString("\tb.WriteString(layoutTail)\n")
			}
		}
		for i, n := range file.Body.Nodes {
			if headEnd >= 0 {
				if i < headEnd {
					continue
				}
			} else if headIdx >= 0 {
				if i < headIdx {
					continue
				}
				if i == headIdx {
					n = TemplateNode{Type: ir.NodeText, Content: headAfter}
				}
			}
			if styleCode != "" && n.Type == ir.NodeText && strings.Contains(n.Content, "</html>") {
				buf.WriteString(styleCode)
				styleCode = ""
			}
			code, err := genLayoutNodeState(gen, n, 1, &inSection)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
		}
	}

	if styleCode != "" {
		buf.WriteString(styleCode)
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")
	return buf.String(), nil
}

func writeHeadDedupe(buf *strings.Builder, varName string) {
	buf.WriteString(fmt.Sprintf("\tif strings.Contains(strings.ToLower(head), \"<title\") {\n"))
	buf.WriteString(fmt.Sprintf("\t\t%s = stripTitleTag(%s)\n", varName, varName))
	buf.WriteString("\t}\n")
	buf.WriteString("\tif strings.Contains(strings.ToLower(head), `name=\"description\"`) || strings.Contains(strings.ToLower(head), `name='description'`) {\n")
	buf.WriteString(fmt.Sprintf("\t\t%s = stripMetaDescriptionTag(%s)\n", varName, varName))
	buf.WriteString("\t}\n")
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
	if n.Type == ir.NodeText && (strings.Contains(n.Content, ir.HeadPlaceholder) || strings.Contains(n.Content, "{#slot}")) {
		parts := splitLayoutText(n.Content)
		var out strings.Builder
		for _, p := range parts {
			switch p {
			case ir.HeadPlaceholder:
				out.WriteString(indent + "b.WriteString(head)\n")
			case "{#slot}":
				out.WriteString(indent + "b.WriteString(content)\n")
			default:
				out.WriteString(indent + fmt.Sprintf("b.WriteString(%s)\n", gogen.GoLiteral(p)))
			}
		}
		return out.String(), nil
	}
	return genTemplateNodeToState(gen, n, depth, "b", inSection)
}

func splitLayoutText(s string) []string {
	var result []string
	for s != "" {
		headIdx := strings.Index(s, ir.HeadPlaceholder)
		slotIdx := strings.Index(s, "{#slot}")

		next := -1
		nextLen := 0
		if headIdx >= 0 && (slotIdx < 0 || headIdx <= slotIdx) {
			next = headIdx
			nextLen = len(ir.HeadPlaceholder)
		} else if slotIdx >= 0 {
			next = slotIdx
			nextLen = len("{#slot}")
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
