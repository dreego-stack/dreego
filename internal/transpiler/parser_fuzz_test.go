package transpiler

import (
	"reflect"
	"strings"
	"testing"
)

func FuzzParser(f *testing.F) {
	seeds := []string{
		"<go>msg := \"hi\"</go>\n<div><p>{{ msg }}</p></div>",
		"<go method=\"post\">x := 1</go><div>{#if x}<p>yes</p>{/if}</div>",
		"<div>{#each items as item}<li>{{ item.name }}</li>{/each}</div>",
		"<div><@Card title=\"Hi\" count={n}/><@Card>slot</@Card></div>",
		"<div>{#slot name}content{/slot}{#slot}default{/slot}</div>",
		"<head><title>T</title></head><div>body</div>",
		"<script>let a = 1 < 2;</script><style>.a > .b { color: red }</style><div>x</div>",
		"<div>{#verbatim}{{ raw }}{/verbatim}</div>",
		"<div><input value=\"x>y\" data-x='a{b'><a href=\"/p/{{ id }}\">link</a></div>",
		"<div><p>{{ user.name | upper }}</p></div>",
		"<div><div><p>nested</p></div></div>",
		"<div><a class=\"nav {#if cond}active{/if}\">x</a></div>",
		"<div>{#if a}{#else if b}{#else}c{/if}</div>",
		"<div>{#each xs as x}{#if x}y{/if}{/each}</div>",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 64<<10 {
			t.Skip("input too large")
		}

		tokens, err := Lex(input)
		if err != nil {
			return
		}

		file, err := NewParser(tokens).Parse()
		if err != nil {
			return
		}

		if n := countTemplateNodes(file.Template); n > 2*len(tokens)+4 {
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

func countTemplateNodes(section *TemplateSection) int {
	if section == nil {
		return 0
	}
	return countNodes(section.Nodes)
}

func countNodes(nodes []TemplateNode) int {
	n := len(nodes)
	for _, node := range nodes {
		n += countNodes(node.Children)
		n += countNodes(node.ElseChildren)
	}
	return n
}

func FuzzParserPreservesGoSection(f *testing.F) {
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
		if strings.Contains(code, "</go>") {
			t.Skip("code closes the section; not a valid <go> body")
		}

		input := "<go>" + code + "</go><div></div>"
		tokens, err := Lex(input)
		if err != nil {
			t.Fatalf("valid source rejected by lexer: %v", err)
		}
		file, err := NewParser(tokens).Parse()
		if err != nil {
			t.Fatalf("valid source rejected by parser: %v", err)
		}
		if len(file.Go) != 1 {
			t.Fatalf("expected 1 go section, got %d", len(file.Go))
		}
		if want := strings.TrimSpace(code); file.Go[0].Code != want {
			t.Fatalf("go section not preserved: got %q, want %q", file.Go[0].Code, want)
		}
	})
}
