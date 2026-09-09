package md

import (
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestTransformPreservesMessageInMarkdownHeading(t *testing.T) {
	nodes, err := TransformNodes([]ir.TemplateNode{
		{Type: ir.NodeText, Content: "# "},
		{Type: ir.NodeMessage, MessageKey: "home.title"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 3 || nodes[0].Content != "<h1>" || nodes[1].Type != ir.NodeMessage || nodes[2].Content != "</h1>" {
		t.Fatalf("nodes = %+v", nodes)
	}
}
