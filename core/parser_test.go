package core

import (
	"strings"
	"testing"
)

// parseExpectError lexes src, parses it, and asserts the returned error
// contains the given substring.
func parseExpectError(t *testing.T, src, want string) {
	t.Helper()
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error %q does not contain %q", err.Error(), want)
	}
}

func TestParseDuplicateDiv(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div><div><p>b</p></div>",
		"duplicate <div> section")
}

func TestParseDuplicateHead(t *testing.T) {
	parseExpectError(t,
		"<head><title>a</title></head><head><title>b</title></head>",
		"duplicate <head> section")
}

func TestParseUnknownSection(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div><x>foo</x>",
		"unknown section <x>")
}

func TestParseExpectedSectionTag(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div>\nplain text after",
		"expected section tag, got Text")
}

// Item 8 (v0.0.25-feedback2): text before a section. A file starting with
// non-section text (e.g. <!doctype html>) used to make the parser swallow ALL
// following tokens as plain template — <go> blocks landed as HTML text instead
// of Go code (silent corrupt). Chosen fix (option a): the parser stops the
// plain template at a section tag and the main loop parses the section, so
// <go> after leading text works as Go code.

// <go> after leading template text is parsed as a Go section (file.Go set),
// not swallowed as template text.
func TestParseTextBeforeGoSectionParsesAsGo(t *testing.T) {
	tokens, err := Lex("<!doctype html>\n<go>msg := \"hi\"</go>\n<div><p>x</p></div>")
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Go) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Go))
	}
	if file.Go[0].Code != "msg := \"hi\"" {
		t.Fatalf("expected Go code %q, got %q", "msg := \"hi\"", file.Go[0].Code)
	}
	if file.Template == nil {
		t.Fatal("expected template section")
	}
	for _, n := range file.Template.Nodes {
		if n.Type == NodeText && strings.Contains(n.Content, "msg :=") {
			t.Fatalf("go body must not land as template text, got %q", n.Content)
		}
	}
}

// A <div> after leading template text is template content (div is only a
// section at the file start), so the whole file parses as one template.
func TestParseTextBeforeDivSectionParsesAsTemplate(t *testing.T) {
	tokens, err := Lex("plain text\n<div><p>x</p></div>")
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Template == nil {
		t.Fatal("expected template section")
	}
	if len(file.Template.Nodes) != 6 {
		t.Fatalf("expected 6 template nodes, got %d", len(file.Template.Nodes))
	}
	if file.Template.Nodes[0].Type != NodeText || file.Template.Nodes[0].Content != "plain text\n" {
		t.Fatalf("expected leading text node, got %+v", file.Template.Nodes[0])
	}
	if file.Template.Nodes[1].Type != NodeText || file.Template.Nodes[1].Content != "<div>" {
		t.Fatalf("expected div open node, got %+v", file.Template.Nodes[1])
	}
}

// Control: a <go> section at the file start still parses as Go code.
func TestParseGoSectionFirstStillWorks(t *testing.T) {
	tokens, err := Lex("<go>msg := \"hi\"</go>\n<div><p>{msg}</p></div>")
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Go) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Go))
	}
	if file.Go[0].Code != "msg := \"hi\"" {
		t.Fatalf("expected Go code %q, got %q", "msg := \"hi\"", file.Go[0].Code)
	}
}

func TestParseGoAttrsVariants(t *testing.T) {
	cases := []struct {
		name  string
		attrs string
		want  string
	}{
		{"double quoted", `type="json"`, "json"},
		{"single quoted", `type='json'`, "json"},
		{"unquoted", `type=json`, "json"},
		{"no type", `method="POST"`, ""},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseGoAttrs(tc.attrs); got != tc.want {
				t.Errorf("parseGoAttrs(%q) = %q, want %q", tc.attrs, got, tc.want)
			}
		})
	}
}
