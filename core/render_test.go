package core

import (
	"errors"
	"testing"
)

func TestRenderNonHTTP(t *testing.T) {
	comp := ComponentFunc(func(ctx *SSRContext) (string, error) {
		ctx.Set("greeting", "hello")
		return "<p>" + ctx.Get("greeting") + " " + ctx.Data("name").(string) + "</p>", nil
	})
	out, err := Render(comp, "name", "world")
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if want := "<p>hello world</p>"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRenderEscapesProps(t *testing.T) {
	comp := ComponentFunc(func(ctx *SSRContext) (string, error) {
		return "<p>" + ctx.Get("name") + "</p>", nil
	})
	out, err := Render(comp, "name", "<script>alert(1)</script>")
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if want := "<p>&lt;script&gt;alert(1)&lt;/script&gt;</p>"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRenderPropKeyNotString(t *testing.T) {
	comp := ComponentFunc(func(ctx *SSRContext) (string, error) {
		return "", nil
	})
	if _, err := Render(comp, 42, "value"); err == nil {
		t.Fatal("expected error for non-string prop key")
	}
}

func TestRenderPropagatesComponentError(t *testing.T) {
	comp := ComponentFunc(func(ctx *SSRContext) (string, error) {
		return "", errors.New("boom")
	})
	if _, err := Render(comp); err == nil {
		t.Fatal("expected component error to propagate")
	}
}
