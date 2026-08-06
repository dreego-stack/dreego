package core

import (
	"testing"
)

func TestParseEachClauseMissingAs(t *testing.T) {
	parseExpectError(t,
		"<div>{#each items}</div>",
		"expected 'items as item'")
}

func TestParseEachClauseValid(t *testing.T) {
	tokens, err := Lex(`<div>{#each items as item}<p>{x}</p>{/each}</div>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	each := file.Template.Nodes[0]
	if each.Type != NodeEach {
		t.Fatalf("node type = %v, want NodeEach", each.Type)
	}
	if each.Items != "items" || each.Item != "item" {
		t.Errorf("each = items:%q item:%q, want items:%q item:%q",
			each.Items, each.Item, "items", "item")
	}
}

func TestParseUnclosedIf(t *testing.T) {
	parseExpectError(t,
		"<div>{#if x}<p>a</p></div>",
		"unclosed {#if}")
}

func TestParseElseInsideElse(t *testing.T) {
	parseExpectError(t,
		"<div>{#if x}<p>a</p>{#else}<p>b</p>{#else}<p>c</p>{/if}</div>",
		"unexpected {#else} or {#else if} inside {#else}")
}

func TestParseExpressionMultipleFilters(t *testing.T) {
	tokens, err := Lex(`<div>{x|upper|raw}</div>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	expr := file.Template.Nodes[0]
	if expr.Type != NodeExpression {
		t.Fatalf("node type = %v, want NodeExpression", expr.Type)
	}
	if expr.Content != "x" {
		t.Errorf("expr = %q, want %q", expr.Content, "x")
	}
	if len(expr.Filters) != 2 || expr.Filters[0] != "upper" || expr.Filters[1] != "raw" {
		t.Errorf("filters = %v, want [upper raw]", expr.Filters)
	}
}
