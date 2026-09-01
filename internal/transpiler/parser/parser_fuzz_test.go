package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func FuzzParser(f *testing.F) {
	seeds := []string{
		"<server>msg := \"hi\"</server>\n<body><p>{{ msg }}</p></body>",
		"<server method=\"post\">x := 1</server><body>{#if x}<p>yes</p>{/if}</body>",
		"<body>{#each items as item}<li>{{ item.name }}</li>{/each}</body>",
		"<body><@Card title=\"Hi\" count={n}/><@Card>slot</@Card></body>",
		"<body>{#slot name}content{/slot}{#slot}default{/slot}</body>",
		"<head><title>T</title></head><body>body</body>",
		"<client>let a = 1 < 2;</client><style>.a > .b { color: red }</style><body>x</body>",
		"<body>{#verbatim}{{ raw }}{/verbatim}</body>",
		"<body><input value=\"x>y\" data-x='a{b'><a href=\"/p/{{ id }}\">link</a></body>",
		"<body><p>{{ user.name | upper }}</p></body>",
		"<body><div><p>nested</p></div></body>",
		"<body><a class=\"nav {#if cond}active{/if}\">x</a></body>",
		"<body>{#if a}{#else if b}{#else}c{/if}</body>",
		"<body>{#each xs as x}{#if x}y{/if}{/each}</body>",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 64<<10 {
			t.Skip("input too large")
		}

		tokens, err := lexer.Lex(input)
		if err != nil {
			return
		}

		file, err := NewParser(tokens).Parse()
		if err != nil {
			return
		}

		if n := countTemplateNodes(file.Body); n > 2*len(tokens)+4 {
			t.Fatalf("node count %d exceeds bound for %d tokens", n, len(tokens))
		}

		again, err := NewParser(tokens).Parse()
		if err != nil {
			t.Fatalf("second parse failed: %v", err)
		}
		if !reflect.DeepEqual(file, again) {
			t.Fatalf("non-deterministic parse output for %q", input)
		}
	})
}

func countTemplateNodes(section *ir.BodySection) int {
	if section == nil {
		return 0
	}
	return countNodes(section.Nodes)
}

func countNodes(nodes []ir.TemplateNode) int {
	n := len(nodes)
	for _, node := range nodes {
		n += countNodes(node.Children)
		n += countNodes(node.ElseChildren)
	}
	return n
}

func FuzzParserPreservesServerSection(f *testing.F) {
	seeds := []string{
		"msg := \"hi\"",
		"if a < b { return a }",
		"x := 1\nfmt.Println(x)",
		"svg := `<svg><path d=\"M12 2\"/></svg>`",
		"for i, v := range xs {\n\tfmt.Println(i, v)\n}",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, code string) {
		if len(code) > 64<<10 {
			t.Skip("input too large")
		}
		if strings.Contains(code, "</server>") {
			t.Skip("code closes the section; not a valid <server> body")
		}

		input := "<server>" + code + "</server><body></body>"
		tokens, err := lexer.Lex(input)
		if err != nil {
			t.Fatalf("valid source rejected by lexer: %v", err)
		}
		file, err := NewParser(tokens).Parse()
		if err != nil {
			t.Fatalf("valid source rejected by parser: %v", err)
		}
		if len(file.Server) != 1 {
			t.Fatalf("expected 1 go section, got %d", len(file.Server))
		}
		if want := strings.TrimSpace(code); file.Server[0].Code != want {
			t.Fatalf("go section not preserved: got %q, want %q", file.Server[0].Code, want)
		}
	})
}
