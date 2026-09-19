package dreefile

import (
	"testing"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func TestParseBodyMarkdownWithDreegoConstructs(t *testing.T) {
	tokens, err := Lex(`<body lang="md"># Title

{#if cond}<@Card/>{/if}

trailing *para*</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Body.Language != "md" {
		t.Fatalf("language = %q, want md", file.Body.Language)
	}
	nodes := file.Body.Nodes
	if len(nodes) != 3 {
		t.Fatalf("got %d nodes, want 3: %+v", len(nodes), nodes)
	}
	if nodes[0].Type != ir.NodeText || nodes[0].Content != "<h1>Title</h1>" {
		t.Errorf("node 0 = %+v, want converted heading", nodes[0])
	}
	if nodes[1].Type != ir.NodeIf || nodes[1].Cond != "cond" {
		t.Errorf("node 1 = %+v, want NodeIf", nodes[1])
	}
	if len(nodes[1].Children) != 1 || nodes[1].Children[0].Type != ir.NodeComponentCall || nodes[1].Children[0].Tag != "Card" {
		t.Errorf("node 1 children = %+v, want component call", nodes[1].Children)
	}
	if nodes[2].Type != ir.NodeText || nodes[2].Content != "<p>trailing <em>para</em></p>" {
		t.Errorf("node 2 = %+v, want converted paragraph", nodes[2])
	}
}

func TestParseBodySharesNodeSliceBetweenBodyAndBodies(t *testing.T) {
	tokens, err := Lex(`<body lang="md"># Title

trailing *para*</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Bodies) == 0 {
		t.Fatal("expected at least one entry in file.Bodies")
	}
	if len(file.Bodies[0].Nodes) != len(file.Body.Nodes) {
		t.Fatalf("file.Bodies[0] has %d nodes, file.Body has %d", len(file.Bodies[0].Nodes), len(file.Body.Nodes))
	}

	gogen.SetNodeSource(file.Body.Nodes, "route.dreego", 100)

	for i := range file.Body.Nodes {
		if got := file.Bodies[0].Nodes[i].Source; got != "route.dreego" {
			t.Fatalf("file.Bodies[0].Nodes[%d].Source = %q, want route.dreego; the body and bodies views do not share the backing slice", i, got)
		}
		if got := file.Bodies[0].Nodes[i].Pos; got != file.Body.Nodes[i].Pos {
			t.Fatalf("file.Bodies[0].Nodes[%d].Pos = %d, want %d", i, got, file.Body.Nodes[i].Pos)
		}
	}

	if file.Body.Nodes[0].Pos < 100 {
		t.Fatalf("SetNodeSource did not offset positions: %d", file.Body.Nodes[0].Pos)
	}
}
