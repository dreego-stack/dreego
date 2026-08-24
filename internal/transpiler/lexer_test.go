package transpiler

import (
	"strings"
	"testing"
)

func TestLexUnclosedSectionTag(t *testing.T) {
	_, err := Lex("<server>hello")
	if err == nil {
		t.Fatal("expected error for unclosed section tag, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed tag <server>") {
		t.Fatalf("expected 'unclosed tag <server>', got %q", err.Error())
	}
}

func TestLexUnexpectedClosingTag(t *testing.T) {
	_, err := Lex("</server>")
	if err == nil {
		t.Fatal("expected error for unexpected closing tag, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected closing tag </server>") {
		t.Fatalf("expected 'unexpected closing tag </server>', got %q", err.Error())
	}
}

func TestLexMismatchedClosingTag(t *testing.T) {
	_, err := Lex("<head>...</style>")
	if err == nil {
		t.Fatal("expected error for mismatched closing tag, got nil")
	}
	if !strings.Contains(err.Error(), "mismatched closing tag </style>, expected </head>") {
		t.Fatalf("expected 'mismatched closing tag </style>, expected </head>', got %q", err.Error())
	}
}

func TestLexEmptyInputReturnsEOF(t *testing.T) {
	tokens, err := Lex("")
	if err != nil {
		t.Fatalf("expected no error for empty input, got %v", err)
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Type != TokenEOF {
		t.Fatalf("expected TokenEOF, got %v", tokens[0].Type)
	}
}

func TestLexPlainTextNoSections(t *testing.T) {
	tokens, err := Lex("hello")
	if err != nil {
		t.Fatalf("expected no error for plain text, got %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens (text + EOF), got %d", len(tokens))
	}
	if tokens[0].Type != TokenText {
		t.Fatalf("expected TokenText, got %v", tokens[0].Type)
	}
	if tokens[0].Value != "hello" {
		t.Fatalf("expected value 'hello', got %q", tokens[0].Value)
	}
	if tokens[1].Type != TokenEOF {
		t.Fatalf("expected TokenEOF, got %v", tokens[1].Type)
	}
}

func TestScanTagUnclosed(t *testing.T) {
	pos := 0
	tok := scanTag("<div", &pos)
	if tok.Type != TokenText {
		t.Fatalf("expected TokenText for unclosed tag, got %v", tok.Type)
	}
	if tok.Value != "<div" {
		t.Fatalf("expected value '<div', got %q", tok.Value)
	}
	if pos != 4 {
		t.Fatalf("expected pos 4, got %d", pos)
	}
}

func TestScanTagSelfClosing(t *testing.T) {
	pos := 0
	tok := scanTag("<img src=x/>", &pos)
	if tok.Type != TokenTagOpen {
		t.Fatalf("expected TokenTagOpen, got %v", tok.Type)
	}
	if tok.Tag != "img" {
		t.Fatalf("expected tag 'img', got %q", tok.Tag)
	}
	if tok.Attr != "src=x" {
		t.Fatalf("expected attr 'src=x', got %q", tok.Attr)
	}
	if !tok.SelfClose {
		t.Fatalf("expected SelfClose true, got false")
	}
	if pos != 12 {
		t.Fatalf("expected pos 12, got %d", pos)
	}
}

func TestScanTagNotSelfClosing(t *testing.T) {
	pos := 0
	tok := scanTag("<img src=x>", &pos)
	if tok.Type != TokenTagOpen {
		t.Fatalf("expected TokenTagOpen, got %v", tok.Type)
	}
	if tok.SelfClose {
		t.Fatalf("expected SelfClose false, got true")
	}
}

func TestScanComponentSelfClose(t *testing.T) {
	pos := 0
	tok := scanTag("<@Comp/>", &pos)
	if tok.Type != TokenComponentSelfClose {
		t.Fatalf("expected TokenComponentSelfClose, got %v", tok.Type)
	}
	if tok.Tag != "Comp" {
		t.Fatalf("expected tag 'Comp', got %q", tok.Tag)
	}
	if pos != 8 {
		t.Fatalf("expected pos 8, got %d", pos)
	}
}

func TestScanComponentUnclosed(t *testing.T) {
	pos := 0
	tok := scanTag("<@Comp", &pos)
	if tok.Type != TokenText {
		t.Fatalf("expected TokenText for unclosed component tag, got %v", tok.Type)
	}
	if tok.Value != "<@Comp" {
		t.Fatalf("expected value '<@Comp', got %q", tok.Value)
	}
	if pos != 6 {
		t.Fatalf("expected pos 6, got %d", pos)
	}
}

// A <server> section body is raw text: <...> inside it must NOT lex as tags, so
// Go comparisons and strings survive verbatim.
func TestLexServerSectionLtIsRawText(t *testing.T) {
	tokens, err := Lex(`<server>msg := "TO: <HASH>"</server>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	var kinds []string
	for _, tok := range tokens {
		if tok.Type == TokenEOF {
			break
		}
		if tok.Type == TokenTagOpen || tok.Type == TokenTagClose {
			kinds = append(kinds, tok.Type.String()+"("+tok.Tag+")")
		}
	}
	want := "TagOpen(server), TagClose(server)"
	if strings.Join(kinds, ", ") != want {
		t.Fatalf("expected %q, got %q", want, strings.Join(kinds, ", "))
	}
}
