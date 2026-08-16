package core

import (
	"strings"
	"testing"
)

func TestGenTemplEmitsHeadStandaloneWithoutLayout(t *testing.T) {
	file := &File{
		Head: &HeadSection{
			Content: "<title>Home</title>",
		},
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>body</p>"},
			},
		},
	}

	out, err := genTempl(NewGenerator(), file, nil, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "b.WriteString(`<title>Home</title>`)") {
		t.Errorf("expected standalone head emitted to b, got:\n%s", out)
	}
	if strings.Contains(out, `c.Set("slot"`) {
		t.Errorf("no layout: must not emit slot wrapping, got:\n%s", out)
	}
}

func TestGenTemplNoHeadWithoutLayout(t *testing.T) {
	file := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>body</p>"},
			},
		},
	}

	out, err := genTempl(NewGenerator(), file, nil, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out, "WriteString(`<title>") {
		t.Errorf("no head section: must not emit head code, got:\n%s", out)
	}
}

// genHead must resolve a {{ title }} expression to a Go expression and escape
// it by default.
func TestGenHeadExpression(t *testing.T) {
	out, err := genHead(`<title>{{ title }}</title>`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `fmt.Sprintf("%v", title)`) {
		t.Errorf("genHead must resolve {{ title }} to an expression, got:\n%s", out)
	}
	if !strings.Contains(out, "dreego.SafeText") {
		t.Errorf("genHead must escape by default, got:\n%s", out)
	}
}

// genHead with raw + upper filters: upper wraps, raw skips escaping.
func TestGenHeadFilterRawUpper(t *testing.T) {
	out, err := genHead(`<title>{{ title|raw|upper }}</title>`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "strings.ToUpper") {
		t.Errorf("genHead upper filter must wrap in strings.ToUpper, got:\n%s", out)
	}
	if strings.Contains(out, "dreego.SafeText") {
		t.Errorf("genHead raw filter must skip escaping, got:\n%s", out)
	}
}

// genHead must not panic or error on an unclosed brace, and must still emit the
// remaining literal text.
func TestGenHeadUnclosedBrace(t *testing.T) {
	out, err := genHead(`<title>{{ title</title>`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "{{ title</title>") {
		t.Errorf("genHead unclosed brace must keep the literal remainder, got:\n%s", out)
	}
}

// A route whose template starts with a doctype must emit the head section
// AFTER the doctype node (at the position of the <head> tag), never before it:
// the rendered body must start with <!doctype html>.
func TestGenTemplHeadAfterDoctypeWithoutLayout(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<meta charset="utf-8"><title>Home</title>`},
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<!doctype html><html><head>"},
				{Type: NodeText, Content: "</head><body><p>body</p></body></html>"},
			},
		},
	}

	out, err := genTempl(NewGenerator(), file, nil, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doctypeIdx := strings.Index(out, "<!doctype html>")
	headIdx := strings.Index(out, "<title>Home</title>")
	if doctypeIdx < 0 {
		t.Fatalf("doctype missing, got:\n%s", out)
	}
	if headIdx < 0 {
		t.Fatalf("head content missing, got:\n%s", out)
	}
	if headIdx < doctypeIdx {
		t.Errorf("head must be emitted after the doctype, got:\n%s", out)
	}
}
