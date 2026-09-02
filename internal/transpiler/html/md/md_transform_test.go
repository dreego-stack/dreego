package md

import (
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestTransformNodesLineAware(t *testing.T) {
	tests := []struct {
		name string
		in   []ir.TemplateNode
		want []ir.TemplateNode
	}{
		{
			name: "heading with expression",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "## "},
				{Type: ir.NodeExpression, Content: "title"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<h2>"},
				{Type: ir.NodeExpression, Content: "title"},
				{Type: ir.NodeText, Content: "</h2>"},
			},
		},
		{
			name: "unordered list item with expression",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "- item "},
				{Type: ir.NodeExpression, Content: "x"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<ul><li>item "},
				{Type: ir.NodeExpression, Content: "x"},
				{Type: ir.NodeText, Content: "</li></ul>"},
			},
		},
		{
			name: "blockquote with expression",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "> "},
				{Type: ir.NodeExpression, Content: "q"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<blockquote>"},
				{Type: ir.NodeExpression, Content: "q"},
				{Type: ir.NodeText, Content: "</blockquote>"},
			},
		},
		{
			name: "paragraph flow across expression",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "Welcome to the "},
				{Type: ir.NodeExpression, Content: "n"},
				{Type: ir.NodeText, Content: " blog."},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<p>Welcome to the "},
				{Type: ir.NodeExpression, Content: "n"},
				{Type: ir.NodeText, Content: " blog.</p>"},
			},
		},
		{
			name: "fenced code block stays raw",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "```go\nfunc main() {}\n```"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<pre><code class="language-go">func main() {}</code></pre>`},
			},
		},
		{
			name: "fenced code block with expression passthrough",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "```\n<b>raw</b>\n```"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<pre><code><b>raw</b></code></pre>"},
			},
		},
		{
			name: "inline html in paragraph",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `Use <a href="/x" class="btn">the link</a> now`},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p>Use <a href="/x" class="btn">the link</a> now</p>`},
			},
		},
		{
			name: "html block raw passthrough",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<div class=\"note\">\nHTML content <b>bold</b>\n</div>"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<div class=\"note\">\nHTML content <b>bold</b>\n</div>"},
			},
		},
		{
			name: "html block with markdown around it",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "# Heading\n\n<div>raw</div>\n\npara after"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<h1>Heading</h1>"},
				{Type: ir.NodeText, Content: "<div>raw</div>"},
				{Type: ir.NodeText, Content: "<p>para after</p>"},
			},
		},
		{
			name: "bare less than stays escaped",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "5 < 6 and x > y"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<p>5 &lt; 6 and x &gt; y</p>"},
			},
		},
		{
			name: "table with expression in cell",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "| a | b |\n| --- | --- |\n| x | "},
				{Type: ir.NodeExpression, Content: "val"},
				{Type: ir.NodeText, Content: " |"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<table><thead><tr><th>a</th><th>b</th></tr></thead><tbody><tr><td>x</td><td>"},
				{Type: ir.NodeExpression, Content: "val"},
				{Type: ir.NodeText, Content: "</td></tr></tbody></table>"},
			},
		},
		{
			name: "nested ul in ul",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "- a\n  - b\n  - c\n- d"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<ul><li>a<ul><li>b</li><li>c</li></ul></li><li>d</li></ul>"},
			},
		},
		{
			name: "nested ol in ul",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "- a\n  1. b\n  2. c\n- d"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<ul><li>a<ol><li>b</li><li>c</li></ol></li><li>d</li></ul>"},
			},
		},
		{
			name: "nested three levels",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "- a\n  - b\n    - c\n- d"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "<ul><li>a<ul><li>b<ul><li>c</li></ul></li></ul></li><li>d</li></ul>"},
			},
		},
		{
			name: "image inline in paragraph",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "see ![alt text](https://example.com/img.png) here"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p>see <img src="https://example.com/img.png" alt="alt text"> here</p>`},
			},
		},
		{
			name: "image standalone paragraph",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "![alt](https://example.com/img.png)"},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p><img src="https://example.com/img.png" alt="alt"></p>`},
			},
		},
		{
			name: "footnote single ref",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "text[^1]\n\n[^1]: The footnote text."},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p>`},
				{Type: ir.NodeText, Content: `<section class="footnotes"><ol><li id="fn-1">The footnote text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`},
			},
		},
		{
			name: "footnote multiple refs",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "a[^1] b[^1]\n\n[^1]: The footnote text."},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p>a<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup> b<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">2</a></sup></p>`},
				{Type: ir.NodeText, Content: `<section class="footnotes"><ol><li id="fn-1">The footnote text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`},
			},
		},
		{
			name: "footnote definition after usage",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "text[^1]\n\n[^1]: The footnote text."},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p>`},
				{Type: ir.NodeText, Content: `<section class="footnotes"><ol><li id="fn-1">The footnote text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`},
			},
		},
		{
			name: "footnote formatting inside definition",
			in: []ir.TemplateNode{
				{Type: ir.NodeText, Content: "text[^1]\n\n[^1]: The *footnote* text."},
			},
			want: []ir.TemplateNode{
				{Type: ir.NodeText, Content: `<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p>`},
				{Type: ir.NodeText, Content: `<section class="footnotes"><ol><li id="fn-1">The <em>footnote</em> text <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformNodes(tt.in)
			if err != nil {
				t.Fatalf("TransformNodes() error = %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d nodes, want %d: %+v", len(got), len(tt.want), got)
			}
			for i, w := range tt.want {
				if got[i].Type != w.Type {
					t.Errorf("node %d type = %v, want %v", i, got[i].Type, w.Type)
				}
				if got[i].Content != w.Content {
					t.Errorf("node %d content = %q, want %q", i, got[i].Content, w.Content)
				}
			}
		})
	}
}

func TestTransformNodesBlankLineSeparatesParagraphs(t *testing.T) {
	in := []ir.TemplateNode{
		{Type: ir.NodeText, Content: "first line\n\nsecond line"},
	}
	got, err := TransformNodes(in)
	if err != nil {
		t.Fatalf("TransformNodes() error = %v", err)
	}
	want := []ir.TemplateNode{
		{Type: ir.NodeText, Content: "<p>first line</p>"},
		{Type: ir.NodeText, Content: "<p>second line</p>"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d nodes, want %d: %+v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i].Content != w.Content {
			t.Errorf("node %d content = %q, want %q", i, got[i].Content, w.Content)
		}
	}
}
