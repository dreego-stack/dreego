package md

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func emitListLines(lines [][]mdSegment, start int, consumed *int, r *Renderer) []ir.TemplateNode {
	indent := indentOfSegs(lines[start])
	items, next := emitListAt(lines, start, indent, r)
	*consumed = next - start
	var out []ir.TemplateNode
	if ulItem.MatchString(strings.TrimSpace(lineRaw(lines[start]))) {
		out = append(out, textNode("<ul>"))
	} else {
		out = append(out, textNode("<ol>"))
	}
	out = append(out, items...)
	if ulItem.MatchString(strings.TrimSpace(lineRaw(lines[start]))) {
		out = append(out, textNode("</ul>"))
	} else {
		out = append(out, textNode("</ol>"))
	}
	return mergeText(out)
}

func emitListAt(lines [][]mdSegment, start, indent int, r *Renderer) ([]ir.TemplateNode, int) {
	var items []ir.TemplateNode
	i := start
	ordered := isOL(strings.TrimSpace(lineRaw(lines[start])))
	for i < len(lines) {
		ind := indentOfSegs(lines[i])
		if ind < indent {
			break
		}
		line := stripIndentSegs(lines[i], indent)
		raw := lineRaw(line)
		trimmed := strings.TrimSpace(raw)
		if isOL(trimmed) != ordered {
			break
		}
		var text []mdSegment
		if m := ulItem.FindStringSubmatch(trimmed); m != nil {
			text = stripPrefix(line, markerLen(raw, ulItem))
		} else if m := olItem.FindStringSubmatch(trimmed); m != nil {
			text = stripPrefix(line, markerLen(raw, olItem))
		} else {
			break
		}
		var li []ir.TemplateNode
		li = append(li, textNode("<li>"))
		li = append(li, lineSegments(trimTrailingSpace(text), false, r)...)
		if i+1 < len(lines) && indentOfSegs(lines[i+1]) > ind {
			childIndent := indentOfSegs(lines[i+1])
			childItems, next := emitListAt(lines, i+1, childIndent, r)
			if ulItem.MatchString(strings.TrimSpace(lineRaw(lines[i+1]))) {
				li = append(li, textNode("<ul>"))
			} else {
				li = append(li, textNode("<ol>"))
			}
			li = append(li, childItems...)
			if ulItem.MatchString(strings.TrimSpace(lineRaw(lines[i+1]))) {
				li = append(li, textNode("</ul>"))
			} else {
				li = append(li, textNode("</ol>"))
			}
			i = next
		}
		li = append(li, textNode("</li>"))
		items = append(items, mergeText(li)...)
		i++
		for i < len(lines) && indentOfSegs(lines[i]) > ind {
			i++
		}
	}
	return items, i - 1
}

func indentOfSegs(line []mdSegment) int {
	n := 0
	for _, s := range line {
		if s.isExpr {
			break
		}
		for _, r := range s.text {
			if r == ' ' {
				n++
			} else if r == '\t' {
				n += 2
			} else {
				return n
			}
		}
	}
	return n
}

func stripIndentSegs(line []mdSegment, n int) []mdSegment {
	var out []mdSegment
	remaining := n
	for _, s := range line {
		if s.isExpr {
			out = append(out, s)
			continue
		}
		if remaining <= 0 {
			out = append(out, s)
			continue
		}
		trimmed := strings.TrimLeft(s.text, " \t")
		removed := len(s.text) - len(trimmed)
		if removed >= remaining {
			out = append(out, mdSegment{isExpr: false, text: s.text[remaining:]})
			remaining = 0
		} else {
			remaining -= removed
		}
	}
	return out
}
