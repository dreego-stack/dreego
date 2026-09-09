package lexer

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func TestLexMessageExpressionsInBodyAndHead(t *testing.T) {
	source := `<head><title>[[ page.title ]]</title></head><body><p>[[ cart.items count=len(items) ]]</p></body>`
	toks, err := Lex(source)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}

	var messages []tokens.Token
	for _, tok := range toks {
		if tok.Type == tokens.TokenMessage {
			messages = append(messages, tok)
		}
	}
	if len(messages) != 2 {
		t.Fatalf("message token count = %d, want 2", len(messages))
	}
	if messages[0].Value != "page.title" {
		t.Errorf("first message = %q, want page.title", messages[0].Value)
	}
	if messages[1].Value != "cart.items count=len(items)" {
		t.Errorf("second message = %q", messages[1].Value)
	}
	if messages[0].Pos != strings.Index(source, "[[ page.title ]]") {
		t.Errorf("first message position = %d", messages[0].Pos)
	}
}

func TestLexMessageExpressionIsRawInServer(t *testing.T) {
	toks, err := Lex(`<server>value := "[[ not.a.message ]]"</server><body></body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	for _, tok := range toks {
		if tok.Type == tokens.TokenMessage {
			t.Fatal("server source produced a message token")
		}
	}
}

func TestLexUnclosedMessageExpression(t *testing.T) {
	_, err := Lex(`<body>[[ page.title</body>`)
	if err == nil {
		t.Fatal("expected an unclosed message error")
	}
	if !strings.Contains(err.Error(), "unclosed message expression") {
		t.Fatalf("unexpected error: %v", err)
	}
}
