package core

import "testing"

// Tests for error paths in parseTemplateNode / parseDivNodes /
// parseComponentNodes / parseSlotNodes (core/parser_section_div.go).
// Each test lexes the input and drives the full parse chain via
// Lex + NewParser(tokens).Parse(), asserting the returned error.

func TestParseTemplateUnexpectedSlotClose(t *testing.T) {
	parseExpectError(t,
		"<div>{/slot}</div>",
		"unexpected {/slot}")
}

func TestParseTemplateElseOutsideIf(t *testing.T) {
	parseExpectError(t,
		"<div>{#else}</div>",
		"unexpected {#else} outside {#if}")
}

func TestParseTemplateUnexpectedToken(t *testing.T) {
	parseExpectError(t,
		"<div>{/if}</div>",
		"unexpected token IfClose in template")
}

func TestParseComponentMismatchedClose(t *testing.T) {
	parseExpectError(t,
		"<div><@a></@b></div>",
		"unexpected </@b>, expected </@a>")
}

func TestParseUnclosedDiv(t *testing.T) {
	parseExpectError(t,
		"<div>",
		"unclosed <div>")
}

func TestParseIfInsideAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div><a class="nav {#if cond}active{/if}">x</a></div>`,
		"{#if} inside attribute value")
}

func TestParseEachInsideAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div><a class="nav {#each items as item}active{/each}">x</a></div>`,
		"{#each} inside attribute value")
}

func TestParseIfInsideComponentAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div><@Card class="nav {#if cond}active{/if}"/></div>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideDivOpenAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div class="nav {#if cond}active{/if}"><p>x</p></div>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideNestedDivAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div><div class="nav {#if cond}active{/if}">x</div></div>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideSingleQuotedAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div><a class='nav {#if cond}active{/if}'>x</a></div>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideSingleQuotedDivAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div class='nav {#if cond}active{/if}'><p>x</p></div>`,
		"{#if} inside attribute value")
}

func TestParseEachInsideSingleQuotedAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<div><a class='nav {#each items as item}active{/each}'>x</a></div>`,
		"{#each} inside attribute value")
}

func TestParseIfOutsideDivAttrStillWorks(t *testing.T) {
	tokens, err := Lex(`<div>{#if cond}<span class="nav active">x</span>{/if}</div>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Template.Nodes) != 1 || file.Template.Nodes[0].Type != NodeIf {
		t.Fatalf("expected NodeIf wrapping the div, got %+v", file.Template.Nodes)
	}
}

func TestParseIfOutsideAttrStillWorks(t *testing.T) {
	tokens, err := Lex(`<div>{#if cond}<a class="nav">x</a>{/if}</div>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Template.Nodes) != 1 || file.Template.Nodes[0].Type != NodeIf {
		t.Fatalf("expected NodeIf wrapping the tag, got %+v", file.Template.Nodes)
	}
}
