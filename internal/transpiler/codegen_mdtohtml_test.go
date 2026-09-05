package transpiler

import "testing"

func TestTranslateMdtohtmlSafeDefault(t *testing.T) {
	got := translateMdtohtml(`html, err := dreego.mdtohtml(post.Content)`)
	want := `html, err := dreego.MarkdownToHTML(post.Content)`
	if got != want {
		t.Errorf("translateMdtohtml = %q, want %q", got, want)
	}
}

func TestTranslateMdtohtmlTrustedTrue(t *testing.T) {
	got := translateMdtohtml(`html, err := dreego.mdtohtml(post.Content, trusted: true)`)
	want := `html, err := dreego.MarkdownToHTMLTrusted(post.Content)`
	if got != want {
		t.Errorf("translateMdtohtml = %q, want %q", got, want)
	}
}

func TestTranslateMdtohtmlTrustedFalse(t *testing.T) {
	got := translateMdtohtml(`html, err := dreego.mdtohtml(post.Content, trusted: false)`)
	want := `html, err := dreego.MarkdownToHTML(post.Content)`
	if got != want {
		t.Errorf("translateMdtohtml = %q, want %q", got, want)
	}
}

func TestTranslateMdtohtmlNestedParens(t *testing.T) {
	got := translateMdtohtml(`html, err := dreego.mdtohtml(render(post.Content), trusted: true)`)
	want := `html, err := dreego.MarkdownToHTMLTrusted(render(post.Content))`
	if got != want {
		t.Errorf("translateMdtohtml = %q, want %q", got, want)
	}
}

func TestTranslateMdtohtmlMultipleCalls(t *testing.T) {
	got := translateMdtohtml(`a := dreego.mdtohtml(x)
b := dreego.mdtohtml(y, trusted: true)`)
	want := `a := dreego.MarkdownToHTML(x)
b := dreego.MarkdownToHTMLTrusted(y)`
	if got != want {
		t.Errorf("translateMdtohtml = %q, want %q", got, want)
	}
}

func TestTranslateMdtohtmlUnbalancedUnchanged(t *testing.T) {
	in := `html, err := dreego.mdtohtml(post.Content`
	got := translateMdtohtml(in)
	if got != in {
		t.Errorf("translateMdtohtml unbalanced = %q, want unchanged %q", got, in)
	}
}

func TestTranslateMdtohtmlNoMarkerUnchanged(t *testing.T) {
	in := `html, err := dreego.MarkdownToHTML(post.Content)`
	got := translateMdtohtml(in)
	if got != in {
		t.Errorf("translateMdtohtml no-marker = %q, want unchanged %q", got, in)
	}
}

func TestTranslateMdtohtmlNonTrustedArgUnchanged(t *testing.T) {
	in := `html, err := dreego.mdtohtml(post.Content, trusted: maybe)`
	got := translateMdtohtml(in)
	if got != in {
		t.Errorf("translateMdtohtml non-trusted arg = %q, want unchanged %q", got, in)
	}
}
