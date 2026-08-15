package core

import (
	"fmt"
	"strings"
)

func genLayoutNode(gen *generator, n TemplateNode, depth int) (string, error) {
	indent := strings.Repeat("\t", depth)
	if n.Type == NodeSlot {
		if n.Content != "" {
			return indent + fmt.Sprintf("b.WriteString(c.Get(\"slot_%s\"))\n", n.Content), nil
		}
		return indent + "b.WriteString(c.Get(\"slot\"))\n", nil
	}
	if n.Type == NodeText && (strings.Contains(n.Content, "{#head}") || strings.Contains(n.Content, "{#slot}")) {
		parts := splitLayoutText(n.Content)
		var out string
		for _, p := range parts {
			switch p {
			case "{#head}":
				out += indent + "b.WriteString(c.Get(\"head\"))\n"
			case "{#slot}":
				out += indent + "b.WriteString(c.Get(\"slot\"))\n"
			default:
				out += indent + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(p))
			}
		}
		return out, nil
	}
	return genTemplateNode(gen, n, depth)
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
