package parser

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func TestParseRejectsDuplicateServerSameMethod(t *testing.T) {
	t.Parallel()
	tokens, err := lexer.Lex(`<server>x</server><server>y</server>`)
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatal("Parse succeeded, want duplicate <server> error")
	}
	if !strings.Contains(err.Error(), "duplicate <server> section for method GET") {
		t.Fatalf("error = %q, want duplicate <server> section for method GET", err)
	}
}

func TestParseRejectsDuplicateServerExplicitPost(t *testing.T) {
	t.Parallel()
	tokens, err := lexer.Lex(`<server method="post">a</server><server method="post">b</server>`)
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatal("Parse succeeded, want duplicate <server> error")
	}
	if !strings.Contains(err.Error(), "duplicate <server> section for method POST") {
		t.Fatalf("error = %q, want duplicate <server> section for method POST", err)
	}
}

func TestParseAllowsServerDefaultAndPost(t *testing.T) {
	t.Parallel()
	tokens, err := lexer.Lex(`<server>x</server><server method="post">y</server>`)
	if err != nil {
		t.Fatalf("Lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(file.Server) != 2 {
		t.Fatalf("len(file.Server) = %d, want 2", len(file.Server))
	}
	if file.Server[0].Method != "GET" {
		t.Fatalf("Server[0].Method = %q, want GET", file.Server[0].Method)
	}
	if file.Server[1].Method != "POST" {
		t.Fatalf("Server[1].Method = %q, want POST", file.Server[1].Method)
	}
}
