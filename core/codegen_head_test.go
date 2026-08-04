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

	out, err := genTempl(file, nil, "abc123", true)
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

	out, err := genTempl(file, nil, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out, "WriteString(`<title>") {
		t.Errorf("no head section: must not emit head code, got:\n%s", out)
	}
}
