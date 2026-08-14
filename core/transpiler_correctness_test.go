package core

import (
	"strings"
	"testing"
)

// Issue 1: `<` comparisons inside <go> blocks must not corrupt the Go code.
// The whole go body is raw text, so `if a < b` round-trips verbatim.
func TestLexGoSectionComparisonRoundTrips(t *testing.T) {
	src := "<go>if a < b { return a }</go>"
	tokens, err := Lex(src)
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
	want := "if a < b { return a }"
	if file.Go[0].Code != want {
		t.Fatalf("expected Go code %q, got %q", want, file.Go[0].Code)
	}
}

// Issue 1: a `<` inside a Go string literal must not be scanned as a tag.
func TestLexGoSectionLtInStringIsRaw(t *testing.T) {
	tokens, err := Lex(`<go>msg := "a < b"</go>`)
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
	want := "TagOpen(go), TagClose(go)"
	if strings.Join(kinds, ", ") != want {
		t.Fatalf("expected %q, got %q", want, strings.Join(kinds, ", "))
	}
}

// Issue 2: `>` inside a quoted attribute value must not terminate the tag.
func TestLexQuotedAttrGtDoesNotEndTag(t *testing.T) {
	src := `<div><input value="x>y"></div>`
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Template.Nodes) != 1 {
		t.Fatalf("expected 1 template node, got %d", len(file.Template.Nodes))
	}
	want := `<input value="x>y">`
	if file.Template.Nodes[0].Content != want {
		t.Fatalf("expected %q, got %q", want, file.Template.Nodes[0].Content)
	}
}

// Issue 2: single-quoted attribute values are quote-aware too.
func TestLexQuotedAttrGtSingleQuotes(t *testing.T) {
	src := `<div><input value='x>y'></div>`
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := `<input value='x>y'>`
	if file.Template.Nodes[0].Content != want {
		t.Fatalf("expected %q, got %q", want, file.Template.Nodes[0].Content)
	}
}

// Issue 2: a `>` inside a quoted component attribute must not end the tag.
func TestLexComponentAttrGtInQuotes(t *testing.T) {
	src := `<div><@Card title="a>b"/></div>`
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := file.Template.Nodes[0]
	if node.Type != NodeComponentCall {
		t.Fatalf("expected NodeComponentCall, got %v", node.Type)
	}
	if node.Attrs != `title="a>b"` {
		t.Fatalf("expected attrs %q, got %q", `title="a>b"`, node.Attrs)
	}
}

// Issue 3: `<` and `>` inside a <script> body are raw text, not tags.
func TestLexScriptBodyPreservesLtGt(t *testing.T) {
	src := "<script>if (a < b && c > d) { run(); }</script>"
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Script == nil {
		t.Fatal("expected script section")
	}
	want := "if (a < b && c > d) { run(); }"
	if file.Script.Code != want {
		t.Fatalf("expected script %q, got %q", want, file.Script.Code)
	}
}

// Issue 3: `<` and `>` inside a <style> body are raw text, not tags.
func TestLexStyleBodyPreservesLtGt(t *testing.T) {
	src := "<style>.a > .b { content: '<x>'; }</style>"
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Style == nil {
		t.Fatal("expected style section")
	}
	want := ".a > .b { content: '<x>'; }"
	if file.Style.Code != want {
		t.Fatalf("expected style %q, got %q", want, file.Style.Code)
	}
}

// Issue 4: `||` inside {{ }} is a Go operator, not a filter pipeline.
func TestParseExpressionOrOperatorNotFilter(t *testing.T) {
	tokens, err := Lex(`<div>{{ a || b }}</div>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := file.Template.Nodes[0]
	if node.Type != NodeExpression {
		t.Fatalf("expected NodeExpression, got %v", node.Type)
	}
	if node.Content != "a || b" {
		t.Fatalf("expected expression %q, got %q", "a || b", node.Content)
	}
	if len(node.Filters) != 0 {
		t.Fatalf("expected no filters, got %v", node.Filters)
	}
}

// Issue 4: a real filter pipeline still splits, and `||` before a filter
// stays in the expression.
func TestParseExpressionOrOperatorWithFilter(t *testing.T) {
	tokens, err := Lex(`<div>{{ a || b | upper }}</div>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := file.Template.Nodes[0]
	if node.Content != "a || b" {
		t.Fatalf("expected expression %q, got %q", "a || b", node.Content)
	}
	if len(node.Filters) != 1 || node.Filters[0] != "upper" {
		t.Fatalf("expected filter [upper], got %v", node.Filters)
	}
}

// Issue 5: an unknown filter in a template expression must fail codegen with
// the source position of the expression.
func TestGenTemplateNodeUnknownFilterErrors(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "x",
		Filters: []string{"nosuchfilter"},
		Pos:     42,
	}
	_, err := genTemplateNode(n, 0)
	if err == nil {
		t.Fatal("expected error for unknown filter, got nil")
	}
	if !strings.Contains(err.Error(), "unknown filter 'nosuchfilter' at position 42") {
		t.Fatalf("expected 'unknown filter ... at position 42', got %q", err.Error())
	}
}

// Issue 5: an unknown filter in a component expression must fail codegen with
// the source position of the expression.
func TestGenTemplateNodeCompUnknownFilterErrors(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "x",
		Filters: []string{"nosuchfilter"},
		Pos:     7,
	}
	_, err := genTemplateNodeComp(n)
	if err == nil {
		t.Fatal("expected error for unknown filter, got nil")
	}
	if !strings.Contains(err.Error(), "unknown filter 'nosuchfilter' at position 7") {
		t.Fatalf("expected 'unknown filter ... at position 7', got %q", err.Error())
	}
}

// Issue 5: an unknown filter in a head expression must fail codegen with the
// source position of the expression.
func TestGenHeadUnknownFilterErrors(t *testing.T) {
	_, err := genHead(`<title>{{ title|nosuchfilter }}</title>`, "b")
	if err == nil {
		t.Fatal("expected error for unknown filter, got nil")
	}
	if !strings.Contains(err.Error(), "unknown filter 'nosuchfilter' at position 7") {
		t.Fatalf("expected 'unknown filter ... at position 7', got %q", err.Error())
	}
}
