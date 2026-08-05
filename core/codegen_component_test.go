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
