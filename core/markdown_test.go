package core

import (
	"strings"
	"testing"
)

func TestMarkdownToHTMLHeadings(t *testing.T) {
	out, err := MarkdownToHTML("# Title\n\n## Sub")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}
	if want := "<h1>Title</h1><h2>Sub</h2>"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestMarkdownToHTMLParagraph(t *testing.T) {
	out, err := MarkdownToHTML("Hello world")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}
	if want := "<p>Hello world</p>"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestMarkdownToHTMLBoldAndLink(t *testing.T) {
	out, err := MarkdownToHTML("**bold** and [link](https://example.com)")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}
	if want := `<p><strong>bold</strong> and <a href="https://example.com">link</a></p>`; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestMarkdownToHTMLTable(t *testing.T) {
	src := "| a | b |\n|---|---|\n| 1 | 2 |"
	out, err := MarkdownToHTML(src)
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}
	if !strings.Contains(out, "<table>") {
		t.Errorf("output %q does not contain a table", out)
	}
}

func TestMarkdownToHTMLControlFlowRenderedAsText(t *testing.T) {
	out, err := MarkdownToHTML("{#if user.IsAdmin}\nadmin\n{/if}")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}
	if want := "<p>{#if user.IsAdmin} admin {/if}</p>"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestMarkdownToHTMLSafeEscapesScript(t *testing.T) {
	out, err := MarkdownToHTML("<script>alert(1)</script>")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}
	if strings.Contains(out, "<script>") {
		t.Errorf("safe mode leaked raw script tag: %q", out)
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Errorf("safe mode did not escape script tag: %q", out)
	}
}

func TestMarkdownToHTMLTrustedPassesScript(t *testing.T) {
	out, err := MarkdownToHTMLTrusted("<script>alert(1)</script>")
	if err != nil {
		t.Fatalf("MarkdownToHTMLTrusted returned error: %v", err)
	}
	if !strings.Contains(out, "<script>alert(1)</script>") {
		t.Errorf("trusted mode did not pass raw script through: %q", out)
	}
}

func TestMarkdownToHTMLTrustedRendersMarkdown(t *testing.T) {
	out, err := MarkdownToHTMLTrusted("# Title\n\n**bold**")
	if err != nil {
		t.Fatalf("MarkdownToHTMLTrusted returned error: %v", err)
	}
	if want := "<h1>Title</h1><p><strong>bold</strong></p>"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}
