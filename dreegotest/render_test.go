package dreegotest_test

import (
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/dreegotest"
)

func TestRenderMatchesRenderComponent(t *testing.T) {
	comp := dreego.ComponentFunc(func(ctx dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte("<section><h1>Welcome</h1><p>safe</p></section>")}, nil
	})
	want := dreegotest.RenderComponent(t, comp)
	got, err := dreego.Render(comp)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if string(got.HTML) != want {
		t.Errorf("Render output differs from RenderComponent:\n got: %q\nwant: %q", got.HTML, want)
	}
}
