package parser

import (
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func TestParseEachClauseMissingAs(t *testing.T) {
	parseExpectError(t,
		"<body>{#each items}</body>",
		"expected 'items as item'")
}

func TestParseEachClauseValid(t *testing.T) {
	tokens, err := lexer.Lex(`<body>{#each items as item}<p>{{ x }}</p>{/each}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	each := file.Body.Nodes[0]
	if each.Type != ir.NodeEach {
		t.Fatalf("node type = %v, want NodeEach", each.Type)
	}
	if each.Items != "items" || each.Item != "item" {
		t.Errorf("each = items:%q item:%q, want items:%q item:%q",
			each.Items, each.Item, "items", "item")
	}
}

func TestParseUnclosedIf(t *testing.T) {
	parseExpectError(t,
		"<body>{#if x}<p>a</p></body>",
		"unclosed {#if}")
}

func TestParseElseInsideElse(t *testing.T) {
	parseExpectError(t,
		"<body>{#if x}<p>a</p>{#else}<p>b</p>{#else}<p>c</p>{/if}</body>",
		"unexpected {#else} or {#else if} inside {#else}")
}

func TestParseExpressionMultipleFilters(t *testing.T) {
	tokens, err := lexer.Lex(`<body>{{ x|upper|raw }}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	expr := file.Body.Nodes[0]
	if expr.Type != ir.NodeExpression {
		t.Fatalf("node type = %v, want NodeExpression", expr.Type)
	}
	if expr.Content != "x" {
		t.Errorf("expr = %q, want %q", expr.Content, "x")
	}
	if len(expr.Filters) != 2 || expr.Filters[0] != "upper" || expr.Filters[1] != "raw" {
		t.Errorf("filters = %v, want [upper raw]", expr.Filters)
	}
}
