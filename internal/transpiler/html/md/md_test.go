package md

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestToNodes(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{
			name: "headings",
			src:  "# One\n## Two\n### Three",
			want: []string{"<h1>One</h1>", "<h2>Two</h2>", "<h3>Three</h3>"},
		},
		{
			name: "paragraph",
			src:  "hello world",
			want: []string{"<p>hello world</p>"},
		},
		{
			name: "paragraph multi line",
			src:  "first line\nsecond line",
			want: []string{"<p>first line second line</p>"},
		},
		{
			name: "unordered list",
			src:  "- a\n- b",
			want: []string{"<ul><li>a</li><li>b</li></ul>"},
		},
		{
			name: "ordered list",
			src:  "1. a\n2. b",
			want: []string{"<ol><li>a</li><li>b</li></ol>"},
		},
		{
			name: "emphasis",
			src:  "*em* and **strong**",
			want: []string{"<p><em>em</em> and <strong>strong</strong></p>"},
		},
		{
			name: "inline code",
			src:  "use `code` here",
			want: []string{"<p>use <code>code</code> here</p>"},
		},
		{
			name: "fenced code block",
			src:  "```go\nfunc main() {}\n```",
			want: []string{`<pre><code class="language-go">func main() {}</code></pre>`},
		},
		{
			name: "fenced code block raw",
			src:  "```\n<b>raw</b>\n```",
			want: []string{"<pre><code><b>raw</b></code></pre>"},
		},
		{
			name: "blockquote",
			src:  "> quoted text",
			want: []string{"<blockquote>quoted text</blockquote>"},
		},
		{
			name: "link",
			src:  "[text](https://example.com)",
			want: []string{`<p><a href="https://example.com">text</a></p>`},
		},
		{
			name: "horizontal rule",
			src:  "---",
			want: []string{"<hr>"},
		},
		{
			name: "escaping",
			src:  "*a* & < 5",
			want: []string{"<p><em>a</em> &amp; &lt; 5</p>"},
		},
		{
			name: "mixed blocks",
			src:  "# Title\n\n- one\n- two\n\npara",
			want: []string{"<h1>Title</h1>", "<ul><li>one</li><li>two</li></ul>", "<p>para</p>"},
		},
		{
			name: "inline html",
			src:  `Use <a href="/x" class="btn">the link</a> now`,
			want: []string{`<p>Use <a href="/x" class="btn">the link</a> now</p>`},
		},
		{
			name: "html block raw",
			src:  `<div class="note">\nHTML content <b>bold</b>\n</div>`,
			want: []string{`<div class="note">\nHTML content <b>bold</b>\n</div>`},
		},
		{
			name: "void element",
			src:  "line one<br>line two",
			want: []string{"<p>line one<br>line two</p>"},
		},
		{
			name: "mixed markdown around html",
			src:  "# Heading\n\n<div>raw</div>\n\npara after",
			want: []string{"<h1>Heading</h1>", "<div>raw</div>", "<p>para after</p>"},
		},
		{
			name: "bare less than stays escaped",
			src:  "5 < 6 and x > y",
			want: []string{"<p>5 &lt; 6 and x &gt; y</p>"},
		},
		{
			name: "less than followed by space stays escaped",
			src:  "a < b",
			want: []string{"<p>a &lt; b</p>"},
		},
		{
			name: "basic table",
			src:  "| a | b |\n| --- | --- |\n| x | y |",
			want: []string{"<table><thead><tr><th>a</th><th>b</th></tr></thead><tbody><tr><td>x</td><td>y</td></tr></tbody></table>"},
		},
		{
			name: "table alignment",
			src:  "| a | b | c |\n| :--- | :---: | ---: |\n| x | y | z |",
			want: []string{`<table><thead><tr><th style="text-align: left">a</th><th style="text-align: center">b</th><th style="text-align: right">c</th></tr></thead><tbody><tr><td style="text-align: left">x</td><td style="text-align: center">y</td><td style="text-align: right">z</td></tr></tbody></table>`},
		},
		{
			name: "table header only",
			src:  "| a | b |\n| --- | --- |",
			want: []string{"<table><thead><tr><th>a</th><th>b</th></tr></thead><tbody></tbody></table>"},
		},
		{
			name: "table inline formatting",
			src:  "| a | b |\n| --- | --- |\n| *em* | `code` |",
			want: []string{"<table><thead><tr><th>a</th><th>b</th></tr></thead><tbody><tr><td><em>em</em></td><td><code>code</code></td></tr></tbody></table>"},
		},
		{
			name: "table mismatched columns",
			src:  "| a | b | c |\n| --- | --- | --- |\n| x |\n| p | q | r | s |",
			want: []string{"<table><thead><tr><th>a</th><th>b</th><th>c</th></tr></thead><tbody><tr><td>x</td><td></td><td></td></tr><tr><td>p</td><td>q</td><td>r</td></tr></tbody></table>"},
		},
		{
			name: "nested ul in ul",
			src:  "- a\n  - b\n  - c\n- d",
			want: []string{"<ul><li>a<ul><li>b</li><li>c</li></ul></li><li>d</li></ul>"},
		},
		{
			name: "nested ol in ul",
			src:  "- a\n  1. b\n  2. c\n- d",
			want: []string{"<ul><li>a<ol><li>b</li><li>c</li></ol></li><li>d</li></ul>"},
		},
		{
			name: "nested three levels",
			src:  "- a\n  - b\n    - c\n- d",
			want: []string{"<ul><li>a<ul><li>b<ul><li>c</li></ul></li></ul></li><li>d</li></ul>"},
		},
		{
			name: "image inline in paragraph",
			src:  "see ![alt text](https://example.com/img.png) here",
			want: []string{`<p>see <img src="https://example.com/img.png" alt="alt text"> here</p>`},
		},
		{
			name: "image standalone paragraph",
			src:  "![alt](https://example.com/img.png)",
			want: []string{`<p><img src="https://example.com/img.png" alt="alt"></p>`},
		},
		{
			name: "footnote single ref",
			src:  "text[^1]\n\n[^1]: The footnote text.",
			want: []string{
				`<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p>`,
				`<section class="footnotes"><ol><li id="fn-1">The footnote text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`,
			},
		},
		{
			name: "footnote multiple refs",
			src:  "a[^1] b[^1]\n\n[^1]: The footnote text.",
			want: []string{
				`<p>a<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup> b<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">2</a></sup></p>`,
				`<section class="footnotes"><ol><li id="fn-1">The footnote text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`,
			},
		},
		{
			name: "footnote definition after usage",
			src:  "text[^1]\n\n[^1]: The footnote text.",
			want: []string{
				`<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p>`,
				`<section class="footnotes"><ol><li id="fn-1">The footnote text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`,
			},
		},
		{
			name: "footnote formatting inside definition",
			src:  "text[^1]\n\n[^1]: The *footnote* text.",
			want: []string{
				`<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p>`,
				`<section class="footnotes"><ol><li id="fn-1">The <em>footnote</em> text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ToNodes(tt.src)
			if err != nil {
				t.Fatalf("ToNodes() error = %v", err)
			}
			if len(nodes) != len(tt.want) {
				t.Fatalf("got %d nodes, want %d: %v", len(nodes), len(tt.want), nodes)
			}
			for i, w := range tt.want {
				if nodes[i].Type != ir.NodeText {
					t.Errorf("node %d type = %v, want NodeText", i, nodes[i].Type)
				}
				if nodes[i].Content != w {
					t.Errorf("node %d content = %q, want %q", i, nodes[i].Content, w)
				}
			}
		})
	}
}

func TestToNodesError(t *testing.T) {
	_, err := ToNodes("-")
	if err == nil {
		t.Fatal("expected error for bare list marker")
	}
	if !strings.Contains(err.Error(), "line 1") {
		t.Errorf("error should mention line number, got: %v", err)
	}
}

func TestTransformNodesPreservesDreegoConstructs(t *testing.T) {
	input := []ir.TemplateNode{
		{Type: ir.NodeText, Content: "# Title"},
		{Type: ir.NodeExpression, Content: "name"},
		{Type: ir.NodeText, Content: "some *em* text"},
		{Type: ir.NodeIf, Cond: "cond", Children: []ir.TemplateNode{
			{Type: ir.NodeComponentCall, Tag: "Card", SelfClose: true},
		}},
		{Type: ir.NodeText, Content: "trailing para"},
	}

	got, err := TransformNodes(input)
	if err != nil {
		t.Fatalf("TransformNodes() error = %v", err)
	}

	if len(got) != 5 {
		t.Fatalf("got %d nodes, want 5: %+v", len(got), got)
	}

	if got[0].Type != ir.NodeText || got[0].Content != "<h1>Title</h1>" {
		t.Errorf("node 0 = %+v, want converted heading", got[0])
	}
	if got[1].Type != ir.NodeExpression || got[1].Content != "name" {
		t.Errorf("node 1 = %+v, want unchanged expression", got[1])
	}
	if got[2].Type != ir.NodeText || got[2].Content != "<p>some <em>em</em> text</p>" {
		t.Errorf("node 2 = %+v, want converted paragraph", got[2])
	}
	if got[3].Type != ir.NodeIf || got[3].Cond != "cond" {
		t.Errorf("node 3 = %+v, want unchanged NodeIf", got[3])
	}
	if len(got[3].Children) != 1 || got[3].Children[0].Type != ir.NodeComponentCall || got[3].Children[0].Tag != "Card" {
		t.Errorf("node 3 children = %+v, want unchanged component call", got[3].Children)
	}
	if got[4].Type != ir.NodeText || got[4].Content != "<p>trailing para</p>" {
		t.Errorf("node 4 = %+v, want converted paragraph", got[4])
	}
}
