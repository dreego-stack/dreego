package core

import (
	"strings"
	"testing"
)

// genComponentCall must resolve {expr} attribute values when a component calls
// another component, mirroring how the route path uses extractAttrValues.
// Currently it passes n.Attrs through raw, emitting invalid Go like
// Link(href={url} label="x").Render(ctx).
func TestGenComponentCallResolvesAttrExpression(t *testing.T) {
	n := TemplateNode{
		Type:      NodeComponentCall,
		Tag:       "Nav.Link",
		Attrs:     `href={url} label="x"`,
		SelfClose: true,
	}
	out, err := genComponentCall(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "{url}") {
		t.Errorf("genComponentCall left {url} literal, must resolve to expression url. got:\n%s", out)
	}
	if !strings.Contains(out, "Link(url, \"x\")") {
		t.Errorf("genComponentCall must pass resolved args Link(url, \"x\"), got:\n%s", out)
	}
}

// compTextWithAttrs is applied to every NodeText in a component body, including
// text inside <script>/<style> sections. The lexer treats those sections as raw
// text blocks where {…} is NOT a Go expression, so compTextWithAttrs must leave
// them untouched. Currently it replaces {x} inside "<script>const s = "{x}";"
// with html.EscapeString(fmt.Sprintf("%v", x)).
func TestCompTextWithAttrsLeavesScriptStyleBodiesUntouched(t *testing.T) {
	cases := []string{
		`<script>const s = "{x}";</script>`,
		`<style>.a { color: "{x}"; }</style>`,
	}
	for _, in := range cases {
		out := compTextWithAttrs(in)
		if !strings.Contains(out, "{x}") {
			t.Errorf("compTextWithAttrs must leave script/style body {x} literal (not resolve as Go expr). input=%q got=%q", in, out)
		}
	}
}
