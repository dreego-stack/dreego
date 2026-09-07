package parser

import (
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func TestParseTypeScriptBodyScript(t *testing.T) {
	tokens, err := lexer.Lex(`<body><main>Ready</main><script lang="ts">const ready: boolean = true;</script></body>`)
	if err != nil {
		t.Fatal(err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Body.Nodes) != 4 {
		t.Fatalf("nodes = %#v", file.Body.Nodes)
	}
	node := file.Body.Nodes[3]
	if node.Type != ir.NodeClientScript || node.Language != "ts" || node.Content != "const ready: boolean = true;" {
		t.Fatalf("TypeScript node = %#v", node)
	}
}

func TestParsePreservesPlainBodyScript(t *testing.T) {
	tokens, err := lexer.Lex(`<body><script type="application/ld+json">{"ready":true}</script></body>`)
	if err != nil {
		t.Fatal(err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Body.Nodes) != 3 || file.Body.Nodes[0].Type != ir.NodeText {
		t.Fatalf("plain script changed: %#v", file.Body.Nodes)
	}
}
