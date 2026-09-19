package md

import (
	shared "github.com/dreego-stack/dreego/internal/md"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

type Mode = shared.Mode

const (
	ModeTrusted = shared.ModeTrusted
	ModeSafe    = shared.ModeSafe
)

func ToNodes(src string, mode Mode) ([]ir.TemplateNode, error) {
	blocks, err := shared.ParseBlocks(src, mode)
	if err != nil {
		return nil, err
	}
	var nodes []ir.TemplateNode
	for _, block := range blocks {
		nodes = append(nodes, textNode(block))
	}
	return nodes, nil
}

func textNode(content string) ir.TemplateNode {
	return ir.TemplateNode{Type: ir.NodeText, Content: content}
}
