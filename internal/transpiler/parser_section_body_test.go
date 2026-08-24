package transpiler

import (
	"strings"
	"testing"
)

// Tests for error paths in parseTemplateNode / parseBodyNodes /
// parseComponentNodes / parseSlotNodes (parser_section_body.go).
// Each test lexes the input and drives the full parse chain via
// Lex + NewParser(tokens).Parse(), asserting the returned error.

func TestParseTemplateUnexpectedSlotClose(t *testing.T) {
	parseExpectError(t,
		"<body>{/slot}</body>",
		"unexpected {/slot}")
}

func TestParseTemplateElseOutsideIf(t *testing.T) {
	parseExpectError(t,
		"<body>{#else}</body>",
		"unexpected {#else} outside {#if}")
}

func TestParseTemplateUnexpectedToken(t *testing.T) {
	parseExpectError(t,
		"<body>{/if}</body>",
		"unexpected token IfClose in template")
}

func TestParseComponentMismatchedClose(t *testing.T) {
	parseExpectError(t,
		"<body><@a></@b></body>",
		"unexpected </@b>, expected </@a>")
}

func TestParseUnclosedBody(t *testing.T) {
	_, err := Lex("<body>")
	if err == nil || !strings.Contains(err.Error(), "unclosed tag <body>") {
		t.Fatalf("expected unclosed body error, got %v", err)
	}
}

func TestParseIfInsideAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body><a class="nav {#if cond}active{/if}">x</a></body>`,
		"{#if} inside attribute value")
}

func TestParseEachInsideAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body><a class="nav {#each items as item}active{/each}">x</a></body>`,
		"{#each} inside attribute value")
}

func TestParseIfInsideComponentAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body><@Card class="nav {#if cond}active{/if}"/></body>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideDivOpenAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body class="nav {#if cond}active{/if}"><p>x</p></body>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideNestedDivAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body><div class="nav {#if cond}active{/if}">x</div></body>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideSingleQuotedAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body><a class='nav {#if cond}active{/if}'>x</a></body>`,
		"{#if} inside attribute value")
}

func TestParseIfInsideSingleQuotedDivAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body class='nav {#if cond}active{/if}'><p>x</p></body>`,
		"{#if} inside attribute value")
}

func TestParseEachInsideSingleQuotedAttrRejected(t *testing.T) {
	parseExpectError(t,
		`<body><a class='nav {#each items as item}active{/each}'>x</a></body>`,
		"{#each} inside attribute value")
}

func TestParseIfOutsideDivAttrStillWorks(t *testing.T) {
	tokens, err := Lex(`<body>{#if cond}<span class="nav active">x</span>{/if}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Body.Nodes) != 1 || file.Body.Nodes[0].Type != NodeIf {
		t.Fatalf("expected NodeIf wrapping the div, got %+v", file.Body.Nodes)
	}
}

func TestParseIfOutsideAttrStillWorks(t *testing.T) {
	tokens, err := Lex(`<body>{#if cond}<a class="nav">x</a>{/if}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Body.Nodes) != 1 || file.Body.Nodes[0].Type != NodeIf {
		t.Fatalf("expected NodeIf wrapping the tag, got %+v", file.Body.Nodes)
	}
}
