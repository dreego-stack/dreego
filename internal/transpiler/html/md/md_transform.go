package md

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

type mdSegment struct {
	isExpr    bool
	text      string
	node      ir.TemplateNode
	protected ir.TemplateNode
}

func TransformNodes(nodes []ir.TemplateNode) ([]ir.TemplateNode, error) {
	lines := buildLines(nodes)
	var out []ir.TemplateNode
	var inline []ir.TemplateNode
	var blockOpen, blockClose string
	var inFence bool
	var fenceLang string
	var fenceContent [][]ir.TemplateNode

	flushBlock := func() {
		if len(inline) == 0 {
			return
		}
		if blockOpen != "" {
			if inline[0].Type == ir.NodeText {
				inline[0].Content = blockOpen + inline[0].Content
			} else {
				inline = append([]ir.TemplateNode{textNode(blockOpen)}, inline...)
			}
		}
		if blockClose != "" {
			last := len(inline) - 1
			if inline[last].Type == ir.NodeText {
				inline[last].Content += blockClose
			} else {
				inline = append(inline, textNode(blockClose))
			}
		}
		out = append(out, mergeText(inline)...)
		inline = nil
		blockOpen, blockClose = "", ""
	}

	emitFence := func() {
		var content []ir.TemplateNode
		for i, fl := range fenceContent {
			if i > 0 {
				content = append(content, textNode("\n"))
			}
			content = append(content, fl...)
		}
		open := "<pre><code"
		if fenceLang != "" {
			open += ` class="language-` + fenceLang + `"`
		}
		open += ">"
		nodes := []ir.TemplateNode{textNode(open)}
		nodes = append(nodes, content...)
		nodes = append(nodes, textNode("</code></pre>"))
		out = append(out, mergeText(nodes)...)
	}

	for idx := 0; idx < len(lines); idx++ {
		line := lines[idx]
		if len(line) == 1 && line[0].protected.Type != 0 {
			flushBlock()
			out = append(out, line[0].protected)
			continue
		}
		raw := lineRaw(line)
		trimmed := strings.TrimSpace(raw)

		if inFence {
			if strings.HasPrefix(trimmed, "```") {
				emitFence()
				inFence = false
				fenceContent = nil
			} else {
				fenceContent = append(fenceContent, lineSegments(line, true))
			}
			continue
		}

		if strings.HasPrefix(trimmed, "```") {
			flushBlock()
			fenceLang = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			inFence = true
			fenceContent = nil
			continue
		}

		if !hasText(line) && blockOpen == "" {
			for _, s := range line {
				out = append(out, s.node)
			}
			continue
		}

		switch {
		case trimmed == "":
			flushBlock()
		case isATX(raw):
			flushBlock()
			level := headingLevel(raw)
			content := trimTrailingSpace(stripPrefix(line, markerLen(raw, atxHeading)))
			if before, expr, after, found := splitAtExpr(content); found {
				inline = append(inline, lineSegments(trimTrailingSpace(before), false)...)
				blockOpen, blockClose = fmt.Sprintf("<h%d>", level), fmt.Sprintf("</h%d>", level)
				flushBlock()
				out = append(out, lineSegments(expr, false)...)
				inline = append(inline, lineSegments(after, false)...)
				blockOpen, blockClose = "<p>", "</p>"
			} else {
				inline = append(inline, lineSegments(content, false)...)
				blockOpen, blockClose = fmt.Sprintf("<h%d>", level), fmt.Sprintf("</h%d>", level)
				flushBlock()
			}
		case isListItem(raw):
			content := trimTrailingSpace(stripPrefix(line, markerLen(raw, listItemRe(raw))))
			if before, expr, after, found := splitAtExpr(content); found {
				if blockOpen == "<ul><li>" || blockOpen == "<ol><li>" {
					inline = append(inline, textNode("</li><li>"))
				} else {
					flushBlock()
					if isOL(raw) {
						blockOpen, blockClose = "<ol><li>", "</li></ol>"
					} else {
						blockOpen, blockClose = "<ul><li>", "</li></ul>"
					}
				}
				inline = append(inline, lineSegments(trimTrailingSpace(before), false)...)
				flushBlock()
				out = append(out, lineSegments(expr, false)...)
				inline = append(inline, lineSegments(after, false)...)
				blockOpen, blockClose = "<p>", "</p>"
			} else {
				if blockOpen == "<ul><li>" || blockOpen == "<ol><li>" {
					inline = append(inline, textNode("</li><li>"))
				} else {
					flushBlock()
					if isOL(raw) {
						blockOpen, blockClose = "<ol><li>", "</li></ol>"
					} else {
						blockOpen, blockClose = "<ul><li>", "</li></ul>"
					}
				}
				inline = append(inline, lineSegments(content, false)...)
			}
		case isBlockquote(raw):
			flushBlock()
			content := trimTrailingSpace(trimLeadingSpace(stripPrefix(line, 1)))
			if before, expr, after, found := splitAtExpr(content); found {
				inline = append(inline, lineSegments(trimTrailingSpace(before), false)...)
				blockOpen, blockClose = "<blockquote>", "</blockquote>"
				flushBlock()
				out = append(out, lineSegments(expr, false)...)
				inline = append(inline, lineSegments(after, false)...)
				blockOpen, blockClose = "<p>", "</p>"
			} else {
				inline = append(inline, lineSegments(content, false)...)
				blockOpen, blockClose = "<blockquote>", "</blockquote>"
				flushBlock()
			}
		case isHR(trimmed):
			flushBlock()
			out = append(out, textNode("<hr>"))
		case htmlBlockStart.MatchString(trimmed):
			flushBlock()
			var raw []string
			for idx < len(lines) && strings.TrimSpace(lineRaw(lines[idx])) != "" {
				raw = append(raw, lineRaw(lines[idx]))
				idx++
			}
			idx--
			out = append(out, textNode(strings.Join(raw, "\n")))
		default:
			if blockOpen != "" {
				inline = append(inline, textNode(" "))
			} else {
				flushBlock()
				blockOpen, blockClose = "<p>", "</p>"
			}
			inline = append(inline, lineSegments(trimTrailingSpace(line), false)...)
		}
	}
	if inFence {
		emitFence()
	}
	flushBlock()
	return out, nil
}

func buildLines(nodes []ir.TemplateNode) [][]mdSegment {
	var lines [][]mdSegment
	var cur []mdSegment
	flush := func() {
		if len(cur) > 0 {
			lines = append(lines, cur)
			cur = nil
		}
	}
	for _, n := range nodes {
		switch n.Type {
		case ir.NodeExpression:
			cur = append(cur, mdSegment{isExpr: true, node: n})
		case ir.NodeText:
			parts := strings.Split(n.Content, "\n")
			for j, part := range parts {
				if j > 0 {
					flush()
				}
				if part == "" {
					lines = append(lines, []mdSegment{{isExpr: false, text: ""}})
				} else {
					cur = append(cur, mdSegment{isExpr: false, text: part})
				}
			}
		default:
			flush()
			lines = append(lines, []mdSegment{{protected: n}})
		}
	}
	flush()
	return lines
}

func lineRaw(line []mdSegment) string {
	var b strings.Builder
	for _, s := range line {
		if s.isExpr {
			b.WriteString("{{expr}}")
		} else {
			b.WriteString(s.text)
		}
	}
	return b.String()
}

func hasText(line []mdSegment) bool {
	for _, s := range line {
		if !s.isExpr {
			return true
		}
	}
	return false
}

func lineSegments(line []mdSegment, raw bool) []ir.TemplateNode {
	var out []ir.TemplateNode
	for _, s := range line {
		if s.isExpr {
			out = append(out, s.node)
		} else if raw {
			out = append(out, textNode(s.text))
		} else {
			out = append(out, textNode(renderInline(s.text)))
		}
	}
	return out
}

func mergeText(nodes []ir.TemplateNode) []ir.TemplateNode {
	var out []ir.TemplateNode
	for _, n := range nodes {
		if len(out) > 0 && out[len(out)-1].Type == ir.NodeText && n.Type == ir.NodeText {
			out[len(out)-1].Content += n.Content
		} else {
			out = append(out, n)
		}
	}
	return out
}

func splitAtExpr(line []mdSegment) (before, expr, after []mdSegment, found bool) {
	for i, s := range line {
		if !s.isExpr {
			continue
		}
		if hasSignificantText(line[i+1:]) {
			return line[:i], line[i : i+1], line[i+1:], true
		}
	}
	return nil, nil, nil, false
}

func hasSignificantText(segs []mdSegment) bool {
	for _, s := range segs {
		if s.isExpr {
			return true
		}
		if strings.TrimSpace(s.text) != "" {
			return true
		}
	}
	return false
}

func stripPrefix(line []mdSegment, n int) []mdSegment {
	var out []mdSegment
	remaining := n
	for _, s := range line {
		if remaining <= 0 {
			out = append(out, s)
			continue
		}
		if s.isExpr {
			out = append(out, s)
			continue
		}
		if len(s.text) <= remaining {
			remaining -= len(s.text)
		} else {
			out = append(out, mdSegment{isExpr: false, text: s.text[remaining:]})
			remaining = 0
		}
	}
	return out
}

func trimLeadingSpace(line []mdSegment) []mdSegment {
	if len(line) == 0 || line[0].isExpr {
		return line
	}
	t := strings.TrimLeft(line[0].text, " ")
	if t == "" {
		return line[1:]
	}
	line[0].text = t
	return line
}

func trimTrailingSpace(line []mdSegment) []mdSegment {
	if len(line) == 0 || line[len(line)-1].isExpr {
		return line
	}
	last := len(line) - 1
	t := strings.TrimRight(line[last].text, " \t")
	if t == "" {
		return line[:last]
	}
	line[last].text = t
	return line
}

func markerLen(raw string, re *regexp.Regexp) int {
	m := re.FindStringSubmatch(raw)
	if m == nil {
		return 0
	}
	return len(raw) - len(m[len(m)-1])
}

func isATX(raw string) bool { return atxHeading.MatchString(raw) }

func headingLevel(raw string) int {
	m := atxHeading.FindStringSubmatch(raw)
	return len(m[1])
}

func isBlockquote(raw string) bool {
	return strings.HasPrefix(strings.TrimSpace(raw), ">")
}

func isListItem(raw string) bool {
	return ulItem.MatchString(raw) || olItem.MatchString(raw)
}

func isOL(raw string) bool { return olItem.MatchString(raw) }

func listItemRe(raw string) *regexp.Regexp {
	if olItem.MatchString(raw) {
		return olItem
	}
	return ulItem
}
