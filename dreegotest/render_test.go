package dreegotest_test

import (
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/dreegotest"
)

func TestRenderMatchesRenderComponent(t *testing.T) {
	comp := dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
		return "<section><h1>" + ctx.Get("title") + "</h1><p>" + ctx.Get("name") + "</p></section>", nil
	})
	props := []any{"title", "Welcome <b>bold</b>", "name", "<script>alert(1)</script>"}
	want := dreegotest.RenderComponent(t, comp, props...)
	got, err := dreego.Render(comp, props...)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if got != want {
		t.Errorf("Render output differs from RenderComponent:\n got: %q\nwant: %q", got, want)
	}
}
