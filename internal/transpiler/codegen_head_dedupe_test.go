package transpiler

import (
	"strings"
	"testing"
)

// Head-merge semantics (v0.0.25-feedback2 item 6):
// When the route head defines a <title>, the layout <title> is dropped from
// the merged output — the route wins because it is more specific. The same
// applies to <meta name="description">. The route head is still inserted at
// the {#head} position and still set into the head slot.

func TestGenTemplHeadMergeDedupesTitle(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<title>Page</title><meta name="description" content="route desc">`},
		Template: &TemplateSection{
			Nodes: []TemplateNode{{Type: NodeText, Content: "<p>page</p>"}},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Head: &HeadSection{Content: "<title>Site</title>\n    {#head}"},
			Template: &TemplateSection{
				Nodes: []TemplateNode{{Type: NodeText, Content: "<div><main>{#slot}</main></div>"}},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `if strings.Contains(pageHead, "<title")`) {
		t.Errorf("runtime title dedupe must be emitted, got:\n%s", out)
	}
	if !strings.Contains(out, "layoutHead = stripTitleTag(layoutHead)") {
		t.Errorf("stripTitleTag call must be emitted, got:\n%s", out)
	}
	if !strings.Contains(out, "b.WriteString(`<title>Page</title>") {
		t.Errorf("route <title> must be emitted at the {#head} position, got:\n%s", out)
	}
	if !strings.Contains(out, "pageHead := b.String()") {
		t.Errorf("route head must still be captured for the head slot, got:\n%s", out)
	}
}

func TestGenTemplHeadMergeDedupesMetaDescription(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<meta name="description" content="route desc">`},
		Template: &TemplateSection{
			Nodes: []TemplateNode{{Type: NodeText, Content: "<p>page</p>"}},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Head: &HeadSection{Content: "<meta name=\"description\" content=\"site desc\">\n{#head}"},
			Template: &TemplateSection{
				Nodes: []TemplateNode{{Type: NodeText, Content: "<div><main>{#slot}</main></div>"}},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `strings.Contains(pageHead, `+"`name=\"description\"`"+`)`) {
		t.Errorf("runtime meta description dedupe must be emitted, got:\n%s", out)
	}
	if !strings.Contains(out, "layoutHead = stripMetaDescriptionTag(layoutHead)") {
		t.Errorf("stripMetaDescriptionTag call must be emitted, got:\n%s", out)
	}
	if !strings.Contains(out, "route desc") {
		t.Errorf("route meta description must be emitted at the {#head} position, got:\n%s", out)
	}
}

// Control: without a route <title>, the layout <title> must be kept — dedupe
// must not remove layout head content the route does not override.
func TestGenTemplHeadMergeKeepsLayoutTitleWithoutRouteTitle(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<meta name="description" content="route desc">`},
		Template: &TemplateSection{
			Nodes: []TemplateNode{{Type: NodeText, Content: "<p>page</p>"}},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Head: &HeadSection{Content: "<title>Site</title>\n    {#head}"},
			Template: &TemplateSection{
				Nodes: []TemplateNode{{Type: NodeText, Content: "<div><main>{#slot}</main></div>"}},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "layoutHead := `") {
		t.Errorf("layout head prefix must be emitted, got:\n%s", out)
	}
	if !strings.Contains(out, "route desc") {
		t.Errorf("route head must still be merged, got:\n%s", out)
	}
}

// Regression (v0.0.25-feedback2 item 6): stripping the layout meta description
// must keep the layout head prefix — non-description tags like <meta charset>
// and <link> before/after the description must survive the dedupe.
func TestGenTemplHeadMergeKeepsCharsetAndLinkWhenStrippingDescription(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<meta name="description" content="route desc">`},
		Template: &TemplateSection{
			Nodes: []TemplateNode{{Type: NodeText, Content: "<p>page</p>"}},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Head: &HeadSection{Content: `<meta charset="utf-8"><meta name="description" content="site desc"><link rel="stylesheet" href="/x.css">\n{#head}`},
			Template: &TemplateSection{
				Nodes: []TemplateNode{{Type: NodeText, Content: "<div><main>{#slot}</main></div>"}},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `charset="utf-8"`) {
		t.Errorf("layout <meta charset> must be kept when stripping the description, got:\n%s", out)
	}
	if !strings.Contains(out, `rel="stylesheet" href="/x.css"`) {
		t.Errorf("layout <link> must be kept when stripping the description, got:\n%s", out)
	}
	if !strings.Contains(out, "route desc") {
		t.Errorf("route meta description must be emitted at the {#head} position, got:\n%s", out)
	}
}
