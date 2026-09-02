package parser

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func parseBodyNodesOf(t *testing.T, src string) []ir.TemplateNode {
	t.Helper()
	tokens, err := lexer.Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file.Body.Nodes
}

func TestParseMdTagWrapsInDiv(t *testing.T) {
	nodes := parseBodyNodesOf(t, `<body><md class="prose">
# Hi

Some *text*
</md></body>`)
	var contents []string
	for _, n := range nodes {
		if n.Type == ir.NodeText {
			contents = append(contents, n.Content)
		}
	}
	joined := strings.Join(contents, "")
	if !strings.Contains(joined, `<div class="prose">`) {
		t.Fatalf("missing div wrapper, got: %q", joined)
	}
	if !strings.Contains(joined, "<h1>Hi</h1>") {
		t.Fatalf("missing converted heading, got: %q", joined)
	}
	if !strings.Contains(joined, "<p>Some <em>text</em></p>") {
		t.Fatalf("missing converted paragraph, got: %q", joined)
	}
	if !strings.Contains(joined, "</div>") {
		t.Fatalf("missing div close, got: %q", joined)
	}
}

func TestParseMdTagWithExpression(t *testing.T) {
	nodes := parseBodyNodesOf(t, `<body><md class="x">Hi {{ name }}</md></body>`)
	var types []ir.TemplateNodeType
	for _, n := range nodes {
		types = append(types, n.Type)
	}
	want := []ir.TemplateNodeType{ir.NodeText, ir.NodeText, ir.NodeExpression, ir.NodeText, ir.NodeText}
	if len(types) != len(want) {
		t.Fatalf("got %d nodes %v, want %d", len(types), types, len(want))
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("node %d type = %v, want %v (all: %v)", i, types[i], want[i], types)
		}
	}
	if nodes[0].Content != `<div class="x">` {
		t.Fatalf("open = %q, want div wrapper", nodes[0].Content)
	}
	if nodes[1].Content != "<p>Hi " {
		t.Fatalf("text before expr = %q, want paragraph open", nodes[1].Content)
	}
	if nodes[3].Content != "</p>" {
		t.Fatalf("text after expr = %q, want paragraph close", nodes[3].Content)
	}
	if nodes[4].Content != "</div>" {
		t.Fatalf("close = %q, want div close", nodes[4].Content)
	}
}

func TestParseMdTagUnclosed(t *testing.T) {
	parseExpectError(t,
		"<body><md class=\"prose\"># Hi</body>",
		"unclosed <md>")
}

func TestParseMdTagInMdBodyRejected(t *testing.T) {
	parseExpectError(t,
		`<body lang="md"><md class="prose"># Hi</md></body>`,
		"md is already the body language")
}

func TestParseMdTagNestedRejected(t *testing.T) {
	parseExpectError(t,
		"<body><md class=\"a\">outer <md class=\"b\">inner</md></md></body>",
		"nested <md>")
}

func TestParseMdTagControlFlowRejected(t *testing.T) {
	parseExpectError(t,
		"<body><md class=\"a\">{#if cond}x{/if}</md></body>",
		"md blocks cannot span dreego control flow")
}

func TestParseMdTagPassesOtherAttrs(t *testing.T) {
	nodes := parseBodyNodesOf(t, `<body><md class="prose" data-x="1"># Hi</md></body>`)
	if nodes[0].Content != `<div class="prose" data-x="1">` {
		t.Fatalf("open = %q, want div with passthrough attrs", nodes[0].Content)
	}
}
