package transpiler

import "testing"

func TestLexSingleBraceIsLiteralText(t *testing.T) {
	tokens, err := Lex(`<body>{value}</body>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(file.Body.Nodes) != 1 || file.Body.Nodes[0].Type != NodeText {
		t.Fatalf("single braces must remain literal text, got %+v", file.Body.Nodes)
	}
	if file.Body.Nodes[0].Content != "{value}" {
		t.Fatalf("expected literal {value}, got %q", file.Body.Nodes[0].Content)
	}
}
