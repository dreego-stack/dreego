package core

import (
	"errors"
	"html"
	"testing"
)

func TestRenderNonHTTP(t *testing.T) {
	comp := ComponentFunc(func(ctx RenderContext) (Result, error) {
		ctx.Set("greeting", "hello")
		return Result{HTML: []byte("<p>" + ctx.Get("greeting") + " world</p>")}, nil
	})
	out, err := Render(comp)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if want := "<p>hello world</p>"; string(out.HTML) != want {
		t.Errorf("output = %q, want %q", out.HTML, want)
	}
}

func TestRenderEscapesProps(t *testing.T) {
	name := "<script>alert(1)</script>"
	comp := ComponentFunc(func(ctx RenderContext) (Result, error) {
		return Result{HTML: []byte("<p>" + html.EscapeString(name) + "</p>")}, nil
	})
	out, err := Render(comp)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if want := "<p>&lt;script&gt;alert(1)&lt;/script&gt;</p>"; string(out.HTML) != want {
		t.Errorf("output = %q, want %q", out.HTML, want)
	}
}

func TestRenderPropagatesComponentError(t *testing.T) {
	comp := ComponentFunc(func(ctx RenderContext) (Result, error) {
		return Result{}, errors.New("boom")
	})
	if _, err := Render(comp); err == nil {
		t.Fatal("expected component error to propagate")
	}
}
