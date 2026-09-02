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
