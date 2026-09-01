package lexer

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func scanBraceAt(input string) (tokens.Token, int, error) {
	pos := 0
	tok, err := scanBrace(input, &pos)
	return tok, pos, err
}

func TestScanBraceAllControlTokens(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  tokens.TokenType
		value string
		pos   int
	}{
		{"else", "{#else}", tokens.TokenElse, "", 7},
		{"if close", "{/if}", tokens.TokenIfClose, "", 5},
		{"each close", "{/each}", tokens.TokenEachClose, "", 7},
		{"each else", "{#each else}", tokens.TokenEachElse, "", 12},
		{"slot", "{#slot}", tokens.TokenSlot, "", 7},
		{"slot open", "{#slot name}", tokens.TokenSlotOpen, "name", 12},
		{"slot close", "{/slot}", tokens.TokenSlotClose, "", 7},
		{"else if", "{#else if cond}", tokens.TokenElseIf, "cond", 15},
		{"if open", "{#if cond}", tokens.TokenIfOpen, "cond", 10},
		{"each open", "{#each items}", tokens.TokenEachOpen, "items", 13},
		{"verbatim", "{#verbatim}raw{/verbatim}", tokens.TokenVerbatim, "raw", 25},
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
	if tok.Type != tokens.TokenExpression {
		t.Fatalf("expected TokenExpression, got %v", tok.Type)
	}
	if tok.Value != "user.Name" {
		t.Fatalf("expected expression user.Name, got %q", tok.Value)
	}
	if pos != len("{{ user.Name }}") {
		t.Fatalf("expected position %d, got %d", len("{{ user.Name }}"), pos)
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
	if tok.Type != tokens.TokenEachOpen {
		t.Fatalf("expected TokenEachOpen, got %v", tok.Type)
	}
	if tok.Value != "items as item else" {
		t.Fatalf("expected value 'items as item else', got %q", tok.Value)
	}
	if pos != 26 {
		t.Fatalf("expected pos 26, got %d", pos)
	}
}
