package parser

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
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
	_, err := lexer.Lex("<body>")
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
	tokens, err := lexer.Lex(`<body>{#if cond}<span class="nav active">x</span>{/if}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Body.Nodes) != 1 || file.Body.Nodes[0].Type != ir.NodeIf {
		t.Fatalf("expected NodeIf wrapping the div, got %+v", file.Body.Nodes)
	}
}

func TestParseIfOutsideAttrStillWorks(t *testing.T) {
	tokens, err := lexer.Lex(`<body>{#if cond}<a class="nav">x</a>{/if}</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Body.Nodes) != 1 || file.Body.Nodes[0].Type != ir.NodeIf {
		t.Fatalf("expected NodeIf wrapping the tag, got %+v", file.Body.Nodes)
	}
}

func TestParseBodyMarkdownWithDreegoConstructs(t *testing.T) {
	tokens, err := lexer.Lex(`<body lang="md"># Title

{#if cond}<@Card/>{/if}

trailing *para*</body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Body.Language != "md" {
		t.Fatalf("language = %q, want md", file.Body.Language)
	}
	nodes := file.Body.Nodes
	if len(nodes) != 3 {
		t.Fatalf("got %d nodes, want 3: %+v", len(nodes), nodes)
	}
	if nodes[0].Type != ir.NodeText || nodes[0].Content != "<h1>Title</h1>" {
		t.Errorf("node 0 = %+v, want converted heading", nodes[0])
	}
	if nodes[1].Type != ir.NodeIf || nodes[1].Cond != "cond" {
		t.Errorf("node 1 = %+v, want NodeIf", nodes[1])
	}
	if len(nodes[1].Children) != 1 || nodes[1].Children[0].Type != ir.NodeComponentCall || nodes[1].Children[0].Tag != "Card" {
		t.Errorf("node 1 children = %+v, want component call", nodes[1].Children)
	}
	if nodes[2].Type != ir.NodeText || nodes[2].Content != "<p>trailing <em>para</em></p>" {
		t.Errorf("node 2 = %+v, want converted paragraph", nodes[2])
	}
}
