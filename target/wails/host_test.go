package wails

import (
	"errors"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestHostRendersRegisteredPage(t *testing.T) {
	app := dreego.New()
	if err := app.RegisterRender("/", dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte("<main>Timer</main>")}, nil
	})); err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	host, err := New(app)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	result, err := host.Render("/")
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if got := string(result.HTML); got != "<main>Timer</main>" {
		t.Fatalf("HTML = %q", got)
	}
}

func TestHostPropagatesUnknownRoute(t *testing.T) {
	host, err := New(dreego.New())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := host.Render("/missing"); !errors.Is(err, dreego.ErrRenderRouteNotFound) {
		t.Fatalf("Render error = %v, want ErrRenderRouteNotFound", err)
	}
}
