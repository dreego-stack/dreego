package transpiler

import (
	"strings"
	"testing"
)

// Issue 1: `<` comparisons inside <server> blocks must not corrupt the Go code.
// The whole go body is raw text, so `if a < b` round-trips verbatim.
func TestLexServerSectionComparisonRoundTrips(t *testing.T) {
	src := "<server>if a < b { return a }</server>"
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Server) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Server))
	}
	want := "if a < b { return a }"
	if file.Server[0].Code != want {
		t.Fatalf("expected Go code %q, got %q", want, file.Server[0].Code)
	}
}

// Issue 1: a `<` inside a Go string literal must not be scanned as a tag.
func TestLexServerSectionLtInStringIsRaw(t *testing.T) {
	tokens, err := Lex(`<server>msg := "a < b"</server>`)
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

// Issue 2: `>` inside a quoted attribute value must not terminate the tag.
func TestLexQuotedAttrGtDoesNotEndTag(t *testing.T) {
	src := `<body><input value="x>y"></body>`
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Body.Nodes) != 1 {
		t.Fatalf("expected 1 template node, got %d", len(file.Body.Nodes))
	}
	want := `<input value="x>y">`
	if file.Body.Nodes[0].Content != want {
		t.Fatalf("expected %q, got %q", want, file.Body.Nodes[0].Content)
	}
}

// Issue 2: single-quoted attribute values are quote-aware too.
func TestLexQuotedAttrGtSingleQuotes(t *testing.T) {
	src := `<body><input value='x>y'></body>`
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := `<input value='x>y'>`
	if file.Body.Nodes[0].Content != want {
		t.Fatalf("expected %q, got %q", want, file.Body.Nodes[0].Content)
	}
}

// Issue 2: a `>` inside a quoted component attribute must not end the tag.
func TestLexComponentAttrGtInQuotes(t *testing.T) {
	src := `<body><@Card title="a>b"/></body>`
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := file.Body.Nodes[0]
	if node.Type != NodeComponentCall {
		t.Fatalf("expected NodeComponentCall, got %v", node.Type)
	}
	if node.Attrs != `title="a>b"` {
		t.Fatalf("expected attrs %q, got %q", `title="a>b"`, node.Attrs)
	}
}

// Issue 3: `<` and `>` inside a <script> body are raw text, not tags.
func TestLexScriptBodyPreservesLtGt(t *testing.T) {
	src := "<client>if (a < b && c > d) { run(); }</client>"
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Client == nil {
		t.Fatal("expected script section")
	}
	want := "if (a < b && c > d) { run(); }"
	if file.Client.Code != want {
		t.Fatalf("expected script %q, got %q", want, file.Client.Code)
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
	tokens, err := Lex(`<body>{{ a || b }}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := file.Body.Nodes[0]
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
	tokens, err := Lex(`<body>{{ a || b | upper }}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := file.Body.Nodes[0]
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
	_, err := genTemplateNode(NewGenerator(), n, 0)
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
	_, err := genTemplateNodeComp(NewGenerator(), n)
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
