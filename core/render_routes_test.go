package core

import (
	"errors"
	"testing"
)

func TestAppRendersRegisteredPage(t *testing.T) {
	app := New()
	component := ComponentFunc(func(RenderContext) (Result, error) {
		return Result{HTML: []byte("<main>Timer</main>")}, nil
	})
	if err := app.RegisterRender("/", component); err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	result, err := app.RenderPage("/")
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}
	if got := string(result.HTML); got != "<main>Timer</main>" {
		t.Fatalf("HTML = %q", got)
	}
}

func TestAppRenderPagePropagatesError(t *testing.T) {
	want := errors.New("render failed")
	app := New()
	if err := app.RegisterRender("/", ComponentFunc(func(RenderContext) (Result, error) {
		return Result{}, want
	})); err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	if _, err := app.RenderPage("/"); !errors.Is(err, want) {
		t.Fatalf("RenderPage error = %v, want %v", err, want)
	}
}

func TestAppRejectsInvalidRenderRegistration(t *testing.T) {
	app := New()
	component := ComponentFunc(func(RenderContext) (Result, error) { return Result{}, nil })
	for _, routePath := range []string{"", "timer", "/timer/../admin"} {
		if err := app.RegisterRender(routePath, component); err == nil {
			t.Fatalf("RegisterRender(%q) must fail", routePath)
		}
	}
	if err := app.RegisterRender("/", nil); err == nil {
		t.Fatal("RegisterRender with nil component must fail")
	}
}

func TestAppRejectsDuplicateRenderRoute(t *testing.T) {
	app := New()
	component := ComponentFunc(func(RenderContext) (Result, error) { return Result{}, nil })
	if err := app.RegisterRender("/", component); err != nil {
		t.Fatalf("first RegisterRender: %v", err)
	}
	if err := app.RegisterRender("/", component); !errors.Is(err, ErrRouteConflict) {
		t.Fatalf("duplicate error = %v, want ErrRouteConflict", err)
	}
}

func TestAppRejectsUnknownRenderRoute(t *testing.T) {
	if _, err := New().RenderPage("/http-only"); !errors.Is(err, ErrRenderRouteNotFound) {
		t.Fatalf("RenderPage error = %v, want ErrRenderRouteNotFound", err)
	}
}

func TestAppFreezesRenderRegistrationAfterBuild(t *testing.T) {
	app := New()
	if err := app.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}
	component := ComponentFunc(func(RenderContext) (Result, error) { return Result{}, nil })
	if err := app.RegisterRender("/", component); !errors.Is(err, ErrAppBuilt) {
		t.Fatalf("RegisterRender error = %v, want ErrAppBuilt", err)
	}
}
