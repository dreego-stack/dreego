package core

import (
	"strings"
	"testing"
)

func scanBraceAt(input string) (Token, int, error) {
	pos := 0
	tok, err := scanBrace(input, &pos)
	return tok, pos, err
}

func TestScanBraceAllControlTokens(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  TokenType
		value string
		pos   int
	}{
		{"else", "{#else}", TokenElse, "", 7},
		{"if close", "{/if}", TokenIfClose, "", 5},
		{"each close", "{/each}", TokenEachClose, "", 7},
		{"each else", "{#each else}", TokenEachElse, "", 12},
		{"slot", "{#slot}", TokenSlot, "", 7},
		{"slot open", "{#slot name}", TokenSlotOpen, "name", 12},
		{"slot close", "{/slot}", TokenSlotClose, "", 7},
		{"else if", "{#else if cond}", TokenElseIf, "cond", 15},
		{"if open", "{#if cond}", TokenIfOpen, "cond", 10},
		{"each open", "{#each items}", TokenEachOpen, "items", 13},
		{"verbatim", "{#verbatim}raw{/verbatim}", TokenVerbatim, "raw", 25},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tok, pos, err := scanBraceAt(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != tc.want {
				t.Fatalf("expected type %v, got %v", tc.want, tok.Type)
			}
			if tok.Value != tc.value {
				t.Fatalf("expected value %q, got %q", tc.value, tok.Value)
			}
			if pos != tc.pos {
				t.Fatalf("expected pos %d, got %d", tc.pos, pos)
			}
		})
	}
}

func TestScanBraceUnclosedIf(t *testing.T) {
	_, _, err := scanBraceAt("{#if ")
	if err == nil {
		t.Fatal("expected error for unclosed {#if, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed {#if at position 0") {
		t.Fatalf("expected 'unclosed {#if at position 0', got %q", err.Error())
	}
}

func TestScanBraceUnclosedExpression(t *testing.T) {
	_, _, err := scanBraceAt("{{ foo")
	if err == nil {
		t.Fatal("expected error for unclosed expression, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed expression at position 0") {
		t.Fatalf("expected 'unclosed expression at position 0', got %q", err.Error())
	}
}

func TestScanBraceDoubleExpression(t *testing.T) {
	tok, pos, err := scanBraceAt("{{ user.Name }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != TokenExpression {
		t.Fatalf("expected TokenExpression, got %v", tok.Type)
	}
	if tok.Value != "user.Name" {
		t.Fatalf("expected expression user.Name, got %q", tok.Value)
	}
	if pos != len("{{ user.Name }}") {
		t.Fatalf("expected position %d, got %d", len("{{ user.Name }}"), pos)
	}
}

func TestLexSingleBraceIsLiteralText(t *testing.T) {
	tokens, err := Lex(`<div>{value}</div>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(file.Template.Nodes) != 1 || file.Template.Nodes[0].Type != NodeText {
		t.Fatalf("single braces must remain literal text, got %+v", file.Template.Nodes)
	}
	if file.Template.Nodes[0].Content != "{value}" {
		t.Fatalf("expected literal {value}, got %q", file.Template.Nodes[0].Content)
	}
}

func TestScanBraceUnclosedEach(t *testing.T) {
	_, _, err := scanBraceAt("{#each ")
	if err == nil {
		t.Fatal("expected error for unclosed {#each, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed {#each at position 0") {
		t.Fatalf("expected 'unclosed {#each at position 0', got %q", err.Error())
	}
}

func TestScanBraceUnclosedSlot(t *testing.T) {
	_, _, err := scanBraceAt("{#slot ")
	if err == nil {
		t.Fatal("expected error for unclosed {#slot, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed {#slot at position 0") {
		t.Fatalf("expected 'unclosed {#slot at position 0', got %q", err.Error())
	}
}

func TestScanBraceUnclosedVerbatim(t *testing.T) {
	_, _, err := scanBraceAt("{#verbatim}")
	if err == nil {
		t.Fatal("expected error for unclosed {#verbatim}, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed {#verbatim} at position 0") {
		t.Fatalf("expected 'unclosed {#verbatim} at position 0', got %q", err.Error())
	}
}

func TestScanBraceUnclosedElseIf(t *testing.T) {
	_, _, err := scanBraceAt("{#else if ")
	if err == nil {
		t.Fatal("expected error for unclosed {#else if, got nil")
	}
	if !strings.Contains(err.Error(), "unclosed {#else if at position 0") {
		t.Fatalf("expected 'unclosed {#else if at position 0', got %q", err.Error())
	}
}

func TestScanBraceEachElse(t *testing.T) {
	tok, pos, err := scanBraceAt("{#each items as item else}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != TokenEachOpen {
		t.Fatalf("expected TokenEachOpen, got %v", tok.Type)
	}
	if tok.Value != "items as item else" {
		t.Fatalf("expected value 'items as item else', got %q", tok.Value)
	}
	if pos != 26 {
		t.Fatalf("expected pos 26, got %d", pos)
	}
}
