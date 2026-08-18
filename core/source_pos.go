package core

import "fmt"

func posToLineCol(src string, pos int) (line, col int) {
	line = 1
	col = 1
	for i := 0; i < pos && i < len(src); i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

func sourceLocation(src string, pos int) string {
	line, col := posToLineCol(src, pos)
	if line == 0 && col == 0 {
		return "?:?"
	}
	return fmt.Sprintf("%d:%d", line, col)
}

func setNodeSource(nodes []TemplateNode, src string, posOffset int) {
	for i := range nodes {
		nodes[i].Source = src
		nodes[i].Pos += posOffset
		setNodeSource(nodes[i].Children, src, posOffset)
		setNodeSource(nodes[i].ElseChildren, src, posOffset)
	}
}

func setSourceText(nodes []TemplateNode, src string) {
	for i := range nodes {
		nodes[i].SourceText = src
		setSourceText(nodes[i].Children, src)
		setSourceText(nodes[i].ElseChildren, src)
	}
}
