package core

import "testing"

func TestNewContextExposesOnlyRenderCapabilities(t *testing.T) {
	ctx := NewContext()
	ctx.Set("name", "Dreego")
	if got := ctx.Get("name"); got != "Dreego" {
		t.Fatalf("Get(name) = %q, want Dreego", got)
	}
	if _, ok := ctx.(interface{ Query(string) string }); ok {
		t.Fatal("render context unexpectedly exposes HTTP query access")
	}
}

func TestRenderReturnsStructuredResultFromTypedComponent(t *testing.T) {
	component := ComponentFunc(func(ctx RenderContext) (Result, error) {
		return Result{HTML: []byte("<p>typed</p>")}, nil
	})
	result, err := Render(component)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if got := string(result.HTML); got != "<p>typed</p>" {
		t.Fatalf("HTML = %q, want %q", got, "<p>typed</p>")
	}
}
