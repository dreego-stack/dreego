package parser

import (
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func TestParseLuaBodyScript(t *testing.T) {
	tokens, err := lexer.Lex(`<body><script lang="lua">local ready = true</script></body>`)
	if err != nil {
		t.Fatal(err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Body.Nodes) != 1 {
		t.Fatalf("nodes = %#v", file.Body.Nodes)
	}
	node := file.Body.Nodes[0]
	if node.Type != ir.NodeClientScript || node.Language != "lua" || node.Content != "local ready = true" {
		t.Fatalf("Lua node = %#v", node)
	}
}
