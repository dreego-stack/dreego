package md

import (
	"strings"

	shared "github.com/dreego-stack/dreego/internal/md"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func buildTableNodes(header, sep []mdSegment, body [][]mdSegment, r *Renderer) []ir.TemplateNode {
	headers := splitCellSegments(header)
	sepCells := splitCellSegments(sep)
	n := len(headers)
	var out []ir.TemplateNode
	out = append(out, textNode("<table><thead><tr>"))
	for i, h := range headers {
		out = append(out, cellTagNodes("th", h, alignAtSegments(sepCells, i), r)...)
	}
	out = append(out, textNode("</tr></thead><tbody>"))
	for _, row := range body {
		cells := splitCellSegments(row)
		out = append(out, textNode("<tr>"))
		for i := range n {
			var c []mdSegment
			if i < len(cells) {
				c = cells[i]
			}
			out = append(out, cellTagNodes("td", c, alignAtSegments(sepCells, i), r)...)
		}
		out = append(out, textNode("</tr>"))
	}
	out = append(out, textNode("</tbody></table>"))
	return mergeText(out)
}

func emitTableLines(lines [][]mdSegment, start int, consumed *int, r *Renderer) []ir.TemplateNode {
	header := lines[start]
	sep := lines[start+1]
	idx := start + 2
	var body [][]mdSegment
	for idx < len(lines) {
		bl := lineRaw(lines[idx])
		if strings.TrimSpace(bl) == "" || !strings.Contains(bl, "|") {
			break
		}
		body = append(body, lines[idx])
		idx++
	}
	*consumed = idx - start - 1
	return buildTableNodes(header, sep, body, r)
}

func splitCellSegments(line []mdSegment) [][]mdSegment {
	var cells [][]mdSegment
	var cur []mdSegment
	flush := func() {
		cells = append(cells, trimCellSegments(cur))
		cur = nil
	}
	for _, s := range line {
		if s.isExpr {
			cur = append(cur, s)
			continue
		}
		parts := strings.Split(s.text, "|")
		for j, p := range parts {
			if j > 0 {
				flush()
			}
			if p != "" {
				cur = append(cur, mdSegment{isExpr: false, text: p})
			}
		}
	}
	flush()
	raw := lineRaw(line)
	if strings.HasPrefix(strings.TrimSpace(raw), "|") && len(cells) > 0 && cellEmpty(cells[0]) {
		cells = cells[1:]
	}
	if strings.HasSuffix(strings.TrimSpace(raw), "|") && len(cells) > 0 && cellEmpty(cells[len(cells)-1]) {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func trimCellSegments(c []mdSegment) []mdSegment {
	for len(c) > 0 && !c[0].isExpr {
		t := strings.TrimLeft(c[0].text, " \t")
		if t == "" {
			c = c[1:]
			continue
		}
		c[0].text = t
		break
	}
	for len(c) > 0 && !c[len(c)-1].isExpr {
		last := len(c) - 1
		t := strings.TrimRight(c[last].text, " \t")
		if t == "" {
			c = c[:last]
			continue
		}
		c[last].text = t
		break
	}
	return c
}

func cellEmpty(c []mdSegment) bool {
	for _, s := range c {
		if s.isExpr {
			return false
		}
		if strings.TrimSpace(s.text) != "" {
			return false
		}
	}
	return true
}

func alignAtSegments(sepCells [][]mdSegment, i int) string {
	if i < len(sepCells) {
		return shared.CellAlign(lineRaw(sepCells[i]))
	}
	return ""
}

func cellTagNodes(tag string, cell []mdSegment, align string, r *Renderer) []ir.TemplateNode {
	var out []ir.TemplateNode
	open := "<" + tag
	if align != "" {
		open += ` style="text-align: ` + align + `"`
	}
	open += ">"
	out = append(out, textNode(open))
	out = append(out, lineSegments(cell, false, r)...)
	out = append(out, textNode("</"+tag+">"))
	return out
}
