package core

import (
	"bytes"
	"html"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/render"
)

func TestRenderParity(t *testing.T) {
	comp := ComponentFunc(func(ctx RenderContext) (Result, error) {
		ctx.Set("greeting", "hello")
		return Result{HTML: []byte("<p>" + ctx.Get("greeting") + "</p>")}, nil
	})
	publicOut, err := Render(comp)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	internalRes, err := render.Component(func(c RenderContext) (Result, error) {
		return comp.Render(c)
	})
	if err != nil {
		t.Fatalf("render.Component returned error: %v", err)
	}
	if !bytes.Equal(publicOut.HTML, internalRes.HTML) {
		t.Errorf("public output %q != internal HTML %q", publicOut, internalRes.HTML)
	}
}

func TestRenderParityEscapesProps(t *testing.T) {
	comp := ComponentFunc(func(ctx RenderContext) (Result, error) {
		ctx.Set("name", "<script>alert(1)</script>")
		return Result{HTML: []byte("<p>" + html.EscapeString(ctx.Get("name")) + "</p>")}, nil
	})
	publicOut, err := Render(comp)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	internalRes, err := render.Component(func(c RenderContext) (Result, error) {
		return comp.Render(c)
	})
	if err != nil {
		t.Fatalf("render.Component returned error: %v", err)
	}
	if !bytes.Equal(publicOut.HTML, internalRes.HTML) {
		t.Errorf("public output %q != internal HTML %q", publicOut, internalRes.HTML)
	}
}
