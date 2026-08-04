package core

import (
	"strings"
	"testing"
)

func TestGenTemplEmitsLayoutWrapping(t *testing.T) {
	file := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>page</p>"},
			},
		},
	}
	layout := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<html>{#slot}</html>"},
			},
		},
	}

	out, err := genTempl(file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		"pageContent := b.String()",
		"b.Reset()",
		`c.Set("slot", pageContent)`,
		`b.WriteString(c.Get("slot"))`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("generated layout code missing %q, got:\n%s", want, out)
		}
	}
}

func TestGenTemplLayoutSlotNodeUsesSlot(t *testing.T) {
	file := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>page</p>"},
			},
		},
	}
	layout := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeSlot},
			},
		},
	}

	out, err := genTempl(file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `c.Set("slot", pageContent)`) {
		t.Errorf("expected page content set into slot, got:\n%s", out)
	}
	if !strings.Contains(out, `b.WriteString(c.Get("slot"))`) {
		t.Errorf("expected layout to render slot, got:\n%s", out)
	}
}
