package md

import (
	"strings"
	"sync"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestSafeModeEscapesRawHTML(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{
			name: "script block escaped",
			src:  "<script>alert(1)</script>",
			want: []string{"&lt;script&gt;alert(1)&lt;/script&gt;"},
		},
		{
			name: "div block escaped",
			src:  "<div>raw</div>",
			want: []string{"&lt;div&gt;raw&lt;/div&gt;"},
		},
		{
			name: "inline img onerror escaped",
			src:  "before <img onerror=alert(1) src=x> after",
			want: []string{"<p>before &lt;img onerror=alert(1) src=x&gt; after</p>"},
		},
		{
			name: "markdown constructs still render",
			src:  "# Title\n\n- a\n- b\n\n*em* **strong**",
			want: []string{
				"<h1>Title</h1>",
				"<ul><li>a</li><li>b</li></ul>",
				"<p><em>em</em> <strong>strong</strong></p>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ToNodes(tt.src, ModeSafe)
			if err != nil {
				t.Fatalf("ToNodes() error = %v", err)
			}
			if len(nodes) != len(tt.want) {
				t.Fatalf("got %d nodes, want %d: %+v", len(nodes), len(tt.want), nodes)
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

func TestURLValidation(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "javascript link blocked", src: "[x](javascript:alert(1))", want: "<p>[x](javascript:alert(1))</p>"},
		{name: "mixed case javascript blocked", src: "[x](JaVaScRiPt:alert(1))", want: "<p>[x](JaVaScRiPt:alert(1))</p>"},
		{name: "control char in scheme blocked", src: "[x](jav\tascript:alert(1))", want: "<p>[x](jav\tascript:alert(1))</p>"},
		{name: "https allowed", src: "[x](https://ok.com)", want: `<p><a href="https://ok.com">x</a></p>`},
		{name: "mailto allowed", src: "[x](mailto:hi@example.com)", want: `<p><a href="mailto:hi@example.com">x</a></p>`},
		{name: "relative path allowed", src: "[x](/blog/x)", want: `<p><a href="/blog/x">x</a></p>`},
		{name: "data image allowed", src: "![x](data:image/png;base64,AAA)", want: `<p><img src="data:image/png;base64,AAA" alt="x"></p>`},
		{name: "data html blocked", src: "![x](data:text/html,<script>)", want: "<p>![x](data:text/html,&lt;script&gt;)</p>"},
		{name: "vbscript blocked", src: "[x](vbscript:msgbox(1))", want: "<p>[x](vbscript:msgbox(1))</p>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ToNodes(tt.src, ModeSafe)
			if err != nil {
				t.Fatalf("ToNodes() error = %v", err)
			}
			if len(nodes) != 1 {
				t.Fatalf("got %d nodes, want 1: %+v", len(nodes), nodes)
			}
			if nodes[0].Content != tt.want {
				t.Errorf("content = %q, want %q", nodes[0].Content, tt.want)
			}
		})
	}
}

func TestSafeModeEscapesFenceLang(t *testing.T) {
	nodes, err := ToNodes("```\"><script>\n```", ModeSafe)
	if err != nil {
		t.Fatalf("ToNodes() error = %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("got %d nodes, want 1: %+v", len(nodes), nodes)
	}
	want := `<pre><code class="language-&#34;&gt;&lt;script&gt;"></code></pre>`
	if nodes[0].Content != want {
		t.Errorf("content = %q, want %q", nodes[0].Content, want)
	}
}

func TestSafeModeEscapesUnclosedTag(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "unclosed iframe", src: "<ifrAme", want: "<p>&lt;ifrAme</p>"},
		{name: "unclosed iframe with attrs", src: "hello <iframe src=javascript:alert(1)//", want: "<p>hello &lt;iframe src=javascript:alert(1)//</p>"},
		{name: "unclosed img onerror", src: "hello <img src=x onerror=alert(1)", want: "<p>hello &lt;img src=x onerror=alert(1)</p>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := ToNodes(tt.src, ModeSafe)
			if err != nil {
				t.Fatalf("ToNodes() error = %v", err)
			}
			if len(nodes) != 1 {
				t.Fatalf("got %d nodes, want 1: %+v", len(nodes), nodes)
			}
			if nodes[0].Content != tt.want {
				t.Errorf("content = %q, want %q", nodes[0].Content, tt.want)
			}
			if strings.Contains(nodes[0].Content, "<iframe") || strings.Contains(nodes[0].Content, "<img") {
				t.Errorf("output contains unescaped executable tag: %q", nodes[0].Content)
			}
		})
	}
}

func TestTrustedModeStillPassesRawHTML(t *testing.T) {
	nodes, err := ToNodes("<script>alert(1)</script>", ModeTrusted)
	if err != nil {
		t.Fatalf("ToNodes() error = %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("got %d nodes, want 1: %+v", len(nodes), nodes)
	}
	if nodes[0].Content != "<script>alert(1)</script>" {
		t.Errorf("content = %q, want raw passthrough", nodes[0].Content)
	}
}

func TestConcurrentToNodesNoRace(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mode := ModeSafe
			if i%2 == 0 {
				mode = ModeTrusted
			}
			nodes, err := ToNodes("# h\n\npara *em* [x](https://ok.com)\n\n[^1]: foot", mode)
			if err != nil {
				t.Errorf("ToNodes() error = %v", err)
			}
			if len(nodes) == 0 {
				t.Error("expected nodes")
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentMarkdownToHTML(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				nodes, err := ToNodes("<div>raw</div>\n\n[a](https://x.com)", ModeSafe)
				if err != nil {
					t.Errorf("ToNodes() error = %v", err)
				}
				if len(nodes) == 0 {
					t.Error("expected nodes")
				}
			}
		}()
	}
	wg.Wait()
}

func TestSafeURLDirect(t *testing.T) {
	tests := []struct {
		raw     string
		isImage bool
		want    string
	}{
		{"https://ok.com", false, "https://ok.com"},
		{"javascript:alert(1)", false, ""},
		{" data:image/png;base64,AAA ", true, "data:image/png;base64,AAA"},
		{"data:text/html,x", true, ""},
		{"/rel/path", false, "/rel/path"},
		{"data:image/jpeg;base64,BB", true, "data:image/jpeg;base64,BB"},
		{"data:image/gif;base64,CC", true, "data:image/gif;base64,CC"},
		{"data:image/webp;base64,DD", true, "data:image/webp;base64,DD"},
		{"data:image/svg+xml;base64,EE", true, ""},
	}
	for _, tt := range tests {
		got := safeURL(tt.raw, tt.isImage)
		if got != tt.want {
			t.Errorf("safeURL(%q, %v) = %q, want %q", tt.raw, tt.isImage, got, tt.want)
		}
	}
}

func TestSafeURLUnescapesHTML(t *testing.T) {
	got := safeURL("https://ok.com/?a=1&amp;b=2", false)
	want := "https://ok.com/?a=1&b=2"
	if got != want {
		t.Errorf("safeURL() = %q, want %q", got, want)
	}
	if strings.Contains(got, "&amp;") {
		t.Errorf("safeURL() kept HTML entity: %q", got)
	}
}
