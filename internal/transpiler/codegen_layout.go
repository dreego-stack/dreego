package transpiler

import (
	"fmt"
	"strings"
)

func GenerateLayout(gen *Generator, file *File, funcName string) (string, error) {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func %s(c *dreego.SSRContext, content, head string) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	layoutHead := ""
	if file.Head != nil {
		layoutHead = file.Head.Content
	}
	headPrefix, headSuffix := splitHeadPlaceholder(layoutHead)

	if headPrefix != "" {
		buf.WriteString(fmt.Sprintf("\tlayoutHead := %s\n", goLiteral(headPrefix)))
		buf.WriteString("\tif strings.Contains(head, \"<title\") {\n")
		buf.WriteString("\t\tlayoutHead = stripTitleTag(layoutHead)\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tif strings.Contains(head, `name=\"description\"`) || strings.Contains(head, `name='description'`) {\n")
		buf.WriteString("\t\tlayoutHead = stripMetaDescriptionTag(layoutHead)\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tb.WriteString(layoutHead)\n")
		buf.WriteString("\tb.WriteString(head)\n")
	} else if headSuffix != "" {
		buf.WriteString("\tb.WriteString(head)\n")
	}

	if headSuffix != "" {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(headSuffix)))
	}

	if file.Template != nil {
		inSection := false
		for _, n := range file.Template.Nodes {
			code, err := genLayoutNodeState(gen, n, 1, &inSection)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
		}
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")
	return buf.String(), nil
}

func splitHeadPlaceholder(head string) (prefix, suffix string) {
	if !strings.Contains(head, "{#head}") {
		return head, ""
	}
	parts := strings.SplitN(head, "{#head}", 2)
	return parts[0], parts[1]
}

func genLayoutNode(gen *Generator, n TemplateNode, depth int) (string, error) {
	inSection := false
	return genLayoutNodeState(gen, n, depth, &inSection)
}

func genLayoutNodeState(gen *Generator, n TemplateNode, depth int, inSection *bool) (string, error) {
	indent := strings.Repeat("\t", depth)
	if n.Type == NodeSlot {
		if n.Content != "" {
			return indent + fmt.Sprintf("b.WriteString(c.Get(\"slot_%s\"))\n", n.Content), nil
		}
		return indent + "b.WriteString(content)\n", nil
	}
	if n.Type == NodeText && (strings.Contains(n.Content, "{#head}") || strings.Contains(n.Content, "{#slot}")) {
		parts := splitLayoutText(n.Content)
		var out string
		for _, p := range parts {
			switch p {
			case "{#head}":
				out += indent + "b.WriteString(head)\n"
			case "{#slot}":
				out += indent + "b.WriteString(content)\n"
			default:
				out += indent + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(p))
			}
		}
		return out, nil
	}
	return genTemplateNodeToState(gen, n, depth, "b", inSection)
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

func layoutHelpers() string {
	return `func stripTitleTag(s string) string {
	for {
		open := strings.Index(s, "<title")
		if open < 0 {
			return s
		}
		closeIdx := strings.Index(s[open:], "</title>")
		if closeIdx < 0 {
			return s
		}
		end := open + closeIdx + len("</title>")
		s = s[:open] + s[end:]
	}
}

func stripMetaDescriptionTag(s string) string {
	offset := 0
	for {
		open := strings.Index(s[offset:], "<meta")
		if open < 0 {
			return s
		}
		open += offset
		end := strings.IndexByte(s[open:], '>')
		if end < 0 {
			return s
		}
		tag := s[open : open+end+1]
		if strings.Contains(tag, ` + "`name=\"description\"`" + `) || strings.Contains(tag, ` + "`name='description'`" + `) {
			s = s[:open] + s[open+end+1:]
			offset = open
			continue
		}
		offset = open + end + 1
	}
}
`
}
