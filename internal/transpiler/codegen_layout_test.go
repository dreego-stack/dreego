package transpiler

import (
	"strings"
	"testing"
)

func TestGenTemplEmitsLayoutWrapping(t *testing.T) {
	file := &File{
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>page</p>"},
			},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Body: &BodySection{
				Nodes: []TemplateNode{
					{Type: NodeText, Content: "<html>{#slot}</html>"},
				},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		"pageContent := b.String()",
		"b.Reset()",
		`c.Set("slot", pageContent)`,
		"layouts.Default(c, pageContent, head)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("generated layout code missing %q, got:\n%s", want, out)
		}
	}
}

func TestGenTemplLayoutSlotNodeUsesSlot(t *testing.T) {
	file := &File{
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>page</p>"},
			},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Body: &BodySection{
				Nodes: []TemplateNode{
					{Type: NodeSlot},
				},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `c.Set("slot", pageContent)`) {
		t.Errorf("expected page content set into slot, got:\n%s", out)
	}
	if !strings.Contains(out, "layouts.Default(c, pageContent, head)") {
		t.Errorf("expected layout call, got:\n%s", out)
	}
}

func TestSplitLayoutTextMixed(t *testing.T) {
	out := splitLayoutText(`x{#head}y{#slot}z`)
	want := []string{"x", "{#head}", "y", "{#slot}", "z"}
	if len(out) != len(want) {
		t.Fatalf("splitLayoutText mixed = %#v, want %d parts", out, len(want))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("splitLayoutText part %d = %q, want %q", i, out[i], w)
		}
	}
}

func TestSplitLayoutTextHeadOnly(t *testing.T) {
	out := splitLayoutText(`{#head}`)
	if len(out) != 1 || out[0] != "{#head}" {
		t.Errorf("splitLayoutText head-only = %#v, want [\"{#head}\"]", out)
	}
}

func TestSplitLayoutTextNoPlaceholders(t *testing.T) {
	out := splitLayoutText(`plain text`)
	if len(out) != 1 || out[0] != "plain text" {
		t.Errorf("splitLayoutText no-placeholders = %#v, want [\"plain text\"]", out)
	}
}

// genLayoutNode with a named NodeSlot must read the named slot from context.
func TestGenLayoutNodeNamedSlot(t *testing.T) {
	n := TemplateNode{Type: NodeSlot, Content: "footer"}
	out, err := genLayoutNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `b.WriteString(c.Get("slot_footer"))`) {
		t.Errorf("genLayoutNode named slot must read slot_footer, got:\n%s", out)
	}
}

