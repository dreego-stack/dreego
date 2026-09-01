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
			src:  "*a* & <b>",
			want: []string{"<p><em>a</em> &amp; &lt;b&gt;</p>"},
		},
		{
			name: "mixed blocks",
			src:  "# Title\n\n- one\n- two\n\npara",
			want: []string{"<h1>Title</h1>", "<ul><li>one</li><li>two</li></ul>", "<p>para</p>"},
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
